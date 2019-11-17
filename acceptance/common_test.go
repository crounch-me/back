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
	"github.com/Sehsyha/crounch-back/util"
	"github.com/oliveagle/jsonpath"
)

type TestExecutor struct {
	RequestBody  string
	Response     *http.Response
	ResponseBody []byte
	UserEmail    string
	UserPassword string
	UserToken    string
}

const (
	BaseURL = "http://localhost:3000"
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

	req, err := http.NewRequest(method, BaseURL+path, body)

	if err != nil {
		return err
	}

	if te.UserToken != "" {
		req.Header.Add("Authorization", te.UserToken)
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

func (te *TestExecutor) imAuthenticatedWithThisRandomUSer() error {
	te.RequestBody = fmt.Sprintf(`
    {
      "email": "%s",
      "password": "%s"
    }
  `,
		te.UserEmail,
		te.UserPassword)
	err := te.iSendARequestOn(http.MethodPost, "/users/login")
	if err != nil {
		return err
	}

	err = te.theStatusCodeIs(http.StatusCreated)

	if err != nil {
		return err
	}

	pattern, err := jsonpath.Compile("$.accessToken")
	if err != nil {
		return err
	}

	var actualData interface{}
	json.Unmarshal(te.ResponseBody, &actualData)
	foundValue, _ := pattern.Lookup(actualData)

	te.UserToken = foundValue.(string)

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

func (te *TestExecutor) iCreateARandomUser() error {
	email := randomEmail()
	password := randomPassword()
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

	te.UserEmail = email
	te.UserPassword = password

	return nil
}

func (te *TestExecutor) theStatusCodeIs(code int) error {
	if te.Response.StatusCode != code {
		return fmt.Errorf("status codes are not the same: actual %d expected %d", te.Response.StatusCode, code)
	}
	return nil
}

func (te *TestExecutor) isAStringEqualTo(path string, expected string) error {
	pattern, err := jsonpath.Compile(path)
	if err != nil {
		return err
	}

	var actualData interface{}
	json.Unmarshal(te.ResponseBody, &actualData)
	foundValue, _ := pattern.Lookup(actualData)

	if foundValue != expected {
		return fmt.Errorf("actual %s should be equal to expected %s", foundValue, expected)
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

func randomEmail() string {
	characters := util.RandStringRunes(3)
	characters = strings.ToLower(characters)
	return fmt.Sprintf("%s@%s.%s", characters, characters, characters)
}

func randomPassword() string {
	return util.RandStringRunes(10)
}

func FeatureContext(s *godog.Suite) {
	te := &TestExecutor{}
	s.Step(`^I use this body$`, te.iUseThisBody)
	s.Step(`^"([^"]*)" has string value "([^"]*)"$`, te.hasStringValue)
	s.Step(`^I send a "([^"]*)" request on "([^"]*)"$`, te.iSendARequestOn)
	s.Step(`^I create these users?$`, te.iCreateTheseUsers)
	s.Step(`^the status code is (\d+)$`, te.theStatusCodeIs)
	s.Step(`^"([^"]*)" is a non empty string$`, te.isANonEmptyString)
	s.Step(`^"([^"]*)" is a string equal to "([^"]*)"$`, te.isAStringEqualTo)
	s.Step(`^I\'m authenticated with this random user$`, te.imAuthenticatedWithThisRandomUSer)
	s.Step(`^I create a random user$`, te.iCreateARandomUser)
}
