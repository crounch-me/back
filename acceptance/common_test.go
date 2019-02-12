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

func (tc *TestExecutor) iUseThisBody(body *gherkin.DocString) error {
	tc.RequestBody = strings.TrimSpace(body.Content)
	return nil
}

func (tc *TestExecutor) hasStringValue(path, expectedValue string) error {
	pattern, err := jsonpath.Compile(path)
	if err != nil {
		return err
	}

	var actualData interface{}
	json.Unmarshal(tc.ResponseBody, &actualData)
	foundValue, _ := pattern.Lookup(actualData)

	if foundValue != expectedValue {
		return fmt.Errorf("actual %s is not equal to expected %s", foundValue, expectedValue)
	}

	return nil
}

func (tc *TestExecutor) iSendARequestOn(method, path string) error {
	var body *strings.Reader
	if method == http.MethodPost || method == http.MethodPut || method != http.MethodPatch {
		body = strings.NewReader(tc.RequestBody)
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

	tc.Response = res
	tc.ResponseBody, err = ioutil.ReadAll(tc.Response.Body)
	if err != nil {
		return err
	}

	return nil
}

func FeatureContext(s *godog.Suite) {
	tc := &TestExecutor{}
	s.Step(`^I use this body$`, tc.iUseThisBody)
	s.Step(`^"([^"]*)" has string value "([^"]*)"$`, tc.hasStringValue)
	s.Step(`^I send a "([^"]*)" request on "([^"]*)"$`, tc.iSendARequestOn)
}
