package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/oliveagle/jsonpath"
)

type TestExecutor struct {
	RequestBody  string
	Response     *http.Response
	ResponseBody []byte
}

const (
	BASE_URL = "http://localhost:3000"
)

func (te *TestExecutor) iUseThisBody(body *gherkin.DocString) error {
	te.RequestBody = strings.TrimSpace(body.Content)
	return nil
}

func (te *TestExecutor) hasStringValue(path, expectedValue string) error {
	pattern, err := jsonpath.Compile(path)
	if err != nil {
		return err
	}

	var actualData interface{}
	json.Unmarshal(te.ResponseBody, &actualData)
	foundValue, _ := pattern.Lookup(actualData)

	if foundValue != expectedValue {
		return fmt.Errorf("actual %s is not equal to expected %s", foundValue, expectedValue)
	}

	return nil
}

func (te *TestExecutor) iSendARequestOn(method, path string) error {
	var body *strings.Reader
	if method == http.MethodPost || method == http.MethodPut || method != http.MethodPatch {
		body = strings.NewReader(te.RequestBody)
	} else if method != http.MethodGet {
		return fmt.Errorf("unknown http method %s", method)
	}

	req, err := http.NewRequest(method, BASE_URL+path, body)

	if err != nil {
		return err
	}

	client := http.Client{Timeout: time.Second * 5}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	te.Response = res
	te.ResponseBody, err = ioutil.ReadAll(te.Response.Body)
	if err != nil {
		return err
	}

	return nil
}

func (te *TestExecutor) iCreateTheseUsers(userDataTable *gherkin.DataTable) error {
	for i, row := range userDataTable.Rows {
		if i != 0 {
			email := strings.TrimSpace(row.Cells[0].Value)
			password := strings.TrimSpace(row.Cells[1].Value)
			te.RequestBody = fmt.Sprintf(`
				{
					"email": "%s",
					"password": "%s"
				}
			`,
				email,
				password)
			err := te.iSendARequestOn(http.MethodPost, "/users")
			if err != nil {
				return err
			}

			err = te.theStatusCodeIs(http.StatusCreated)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (te *TestExecutor) theStatusCodeIs(code int) error {
	if te.Response.StatusCode != code {
		return fmt.Errorf("status codes are not the same: actual %d expected %d", te.Response.StatusCode, code)
	}
	return nil
}

func (te *TestExecutor) isANonEmptyString(path string) error {
	pattern, err := jsonpath.Compile(path)
	if err != nil {
		return err
	}

	var actualData interface{}
	json.Unmarshal(te.ResponseBody, &actualData)
	foundValue, _ := pattern.Lookup(actualData)

	if foundValue == "" {
		return fmt.Errorf("actual %s should not be empty", foundValue)
	}

	return nil
}

func FeatureContext(s *godog.Suite) {
	te := &TestExecutor{}
	s.Step(`^I use this body$`, te.iUseThisBody)
	s.Step(`^"([^"]*)" has string value "([^"]*)"$`, te.hasStringValue)
	s.Step(`^I send a "([^"]*)" request on "([^"]*)"$`, te.iSendARequestOn)
	s.Step(`^I create these users?$`, te.iCreateTheseUsers)
	s.Step(`^the status code is (\d+)$`, te.theStatusCodeIs)
	s.Step(`^"([^"]*)" is a non empty string$`, te.isANonEmptyString)
}
