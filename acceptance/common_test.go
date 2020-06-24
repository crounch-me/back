package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/crounch-me/back/domain/lists"
	"github.com/crounch-me/back/domain/products"
	"github.com/crounch-me/back/util"
	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
	"github.com/oliveagle/jsonpath"
)

type ExecutorVariables struct {
	ListID    string
	ProductID string
}

type TestExecutor struct {
	RequestBody  string
	Response     *http.Response
	ResponseBody []byte
	UserEmail    string
	UserPassword string
	UserToken    string
	Variables    ExecutorVariables
}

const (
	BaseURL = "http://localhost:3000"
)

func (te *TestExecutor) iUseThisBody(body *messages.PickleStepArgument_PickleDocString) error {
	te.RequestBody = strings.TrimSpace(body.Content)
	return nil
}

func (te *TestExecutor) getValue(path string) (interface{}, error) {
	pattern, err := jsonpath.Compile(path)

	if err != nil {
		return nil, err
	}

	var actualData interface{}

	json.Unmarshal(te.ResponseBody, &actualData)
	foundValue, err := pattern.Lookup(actualData)

	if err != nil {
		return nil, err
	}

	return foundValue, nil
}

func (te *TestExecutor) hasStringValue(path, expectedValue string) error {
	foundValue, err := te.getValue(path)

	if err != nil {
		return err
	}

	if foundValue != expectedValue {
		return fmt.Errorf("actual %s is not equal to expected %s", foundValue, expectedValue)
	}

	return nil
}

func (te *TestExecutor) hasBoolValue(path, expectedValue string) error {
	foundValue, err := te.getValue(path)

	if err != nil {
		return err
	}

	expectedBoolValue, err := strconv.ParseBool(expectedValue)
	if err != nil {
		return err
	}

	if foundValue.(bool) != expectedBoolValue {
		return fmt.Errorf("actual %s is not equal to expected %s for path %s", foundValue, expectedValue, path)
	}

	return nil
}

func (te *TestExecutor) iSendARequestOn(method, path string) error {
	var b strings.Builder
	var u strings.Builder
	if method != http.MethodPost &&
		method != http.MethodPut &&
		method != http.MethodPatch &&
		method != http.MethodGet &&
		method != http.MethodOptions {
		return fmt.Errorf("unknown http method %s", method)
	}

	tmpl, err := template.New("body").Parse(te.RequestBody)

	if err != nil {
		return err
	}

	err = tmpl.Execute(&b, te.Variables)

	if err != nil {
		return err
	}

	replacedBody := b.String()
	body := *strings.NewReader(replacedBody)

	tmpl, err = template.New("url").Parse(path)

	if err != nil {
		return err
	}

	err = tmpl.Execute(&u, te.Variables)

	if err != nil {
		return err
	}

	url := u.String()

	req, err := http.NewRequest(method, BaseURL+url, &body)

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

func (te *TestExecutor) iCreateTheseUsers(userDataTable *messages.PickleStepArgument_PickleTable) error {
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
	password := randomString()
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
	var e strings.Builder
	pattern, err := jsonpath.Compile(path)

	if err != nil {
		return err
	}

	var actualData interface{}
	json.Unmarshal(te.ResponseBody, &actualData)
	foundValue, _ := pattern.Lookup(actualData)

	tmpl, err := template.New("body-string").Parse(expected)

	if err != nil {
		return err
	}

	err = tmpl.Execute(&e, te.Variables)

	if err != nil {
		return err
	}

	realExpected := e.String()

	if foundValue != realExpected {
		return fmt.Errorf("actual %s should be equal to expected %s", foundValue, realExpected)
	}

	return nil
}

func (te *TestExecutor) theBodyIsAnEmptyArray() error {
	if string(te.ResponseBody) == "[]" {
		return fmt.Errorf("the body is not empty, actual value %s", te.ResponseBody)
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
		return fmt.Errorf("should not be empty")
	}

	return nil
}

func (te *TestExecutor) iCreateAndAuthenticateWithARandomUser() error {
	err := te.iCreateARandomUser()

	if err != nil {
		return err
	}

	err = te.imAuthenticatedWithThisRandomUSer()

	return err
}

func (te *TestExecutor) createList(l *lists.List) error {
	te.RequestBody = fmt.Sprintf(`
    {
      "name": "%s"
    }
  `, l.Name)
	err := te.iSendARequestOn(http.MethodPost, "/lists")
	if err != nil {
		return err
	}

	err = te.theStatusCodeIs(http.StatusCreated)

	if err != nil {
		return err
	}

	id, err := te.getValue("$.id")

	if err != nil {
		return err
	}

	te.Variables.ListID = id.(string)

	return nil
}

func (te *TestExecutor) iCreateTheseLists(listDataTable *messages.PickleStepArgument_PickleTable) error {
	for i, row := range listDataTable.Rows {
		if i != 0 {
			name := strings.TrimSpace(row.Cells[0].Value)
			l := &lists.List{
				Name: name,
			}
			err := te.createList(l)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (te *TestExecutor) iCreateTheseProducts(productDataTable *messages.PickleStepArgument_PickleTable) error {
	for i, row := range productDataTable.Rows {
		if i != 0 {
			name := strings.TrimSpace(row.Cells[0].Value)
			p := &products.Product{
				Name: name,
			}
			err := te.createProduct(p)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (te *TestExecutor) createProduct(p *products.Product) error {
	te.RequestBody = fmt.Sprintf(`
    {
      "name": "%s"
    }
  `, p.Name)
	err := te.iSendARequestOn(http.MethodPost, "/products")
	if err != nil {
		return err
	}

	err = te.theStatusCodeIs(http.StatusCreated)

	if err != nil {
		return err
	}

	id, err := te.getValue("$.id")

	if err != nil {
		return err
	}

	te.Variables.ProductID = id.(string)

	return nil
}

func (te *TestExecutor) iCreateARandomList() error {
	l := &lists.List{
		Name: randomString(),
	}

	err := te.createList(l)

	if err != nil {
		return err
	}

	return nil
}

func (te *TestExecutor) theHeaderEquals(header, value string) error {
	headerValue := te.Response.Header.Get(header)

	if headerValue != value {
		return fmt.Errorf("actual value \"%s\" for header \"%s\" should equal \"%s\"", headerValue, header, value)
	}

	return nil
}

func randomEmail() string {
	return fmt.Sprintf("%s@crounch.me", util.RandString(10))
}

func randomString() string {
	return util.RandString(10)
}

func FeatureContext(s *godog.Suite) {
	te := &TestExecutor{
		Variables: ExecutorVariables{
			ListID:    "",
			ProductID: "",
		},
	}

	// Requests
	s.Step(`^I use this body$`, te.iUseThisBody)
	s.Step(`^I send a "([^"]*)" request on "([^"]*)"$`, te.iSendARequestOn)

	// Assertions
	s.Step(`^the header "([^"]*)" equals "([^"]*)"$`, te.theHeaderEquals)
	s.Step(`^the status code is (\d+)$`, te.theStatusCodeIs)
	s.Step(`^"([^"]*)" has string value "([^"]*)"$`, te.hasStringValue)
	s.Step(`^"([^"]*)" has bool value "([^"]*)"$`, te.hasBoolValue)
	s.Step(`^"([^"]*)" is a string equal to "([^"]*)"$`, te.isAStringEqualTo)
	s.Step(`^"([^"]*)" is a non empty string$`, te.isANonEmptyString)
	s.Step(`^the body is an empty array$`, te.theBodyIsAnEmptyArray)

	// Authentication
	s.Step(`^I\'m authenticated with this random user$`, te.imAuthenticatedWithThisRandomUSer)
	s.Step(`^I create and authenticate with a random user$`, te.iCreateAndAuthenticateWithARandomUser)
	s.Step(`^I create a random user$`, te.iCreateARandomUser)

	// Users
	s.Step(`^I create these users?$`, te.iCreateTheseUsers)

	// Lists
	s.Step(`^I create these lists$`, te.iCreateTheseLists)
	s.Step(`^I create a random list$`, te.iCreateARandomList)

	// Products
	s.Step(`^I create these products$`, te.iCreateTheseProducts)
}
