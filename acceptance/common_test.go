package acceptance

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/crounch-me/back/internal/common/server"
	"github.com/crounch-me/back/internal/common/utils"
	listingPorts "github.com/crounch-me/back/internal/listing/ports"
	productsPorts "github.com/crounch-me/back/internal/products/ports"
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
	Generation   utils.GenerationLibrary
}

const (
	BaseURL = "http://localhost:3000"
)

func (te *TestExecutor) iUseThisBody(body *messages.PickleStepArgument_PickleDocString) error {
	te.RequestBody = strings.TrimSpace(body.Content)
	return nil
}

func (te *TestExecutor) iUseAnInvalidBody() error {
	te.RequestBody = ""
	return nil
}

func (te *TestExecutor) iUseAnEmptyValidBody() error {
	te.RequestBody = "{}"
	return nil
}

func (te *TestExecutor) getValueFromHeader(header string) string {
	return te.Response.Header.Get(header)
}

func (te *TestExecutor) getValueFromBody(path string) (interface{}, error) {
	pattern, err := jsonpath.Compile(path)

	if err != nil {
		return nil, err
	}

	var actualData interface{}

	json.Unmarshal(te.ResponseBody, &actualData)
	foundValue, err := pattern.Lookup(actualData)

	if err != nil {
		return nil, fmt.Errorf("%s, actual value %s", err, actualData)
	}

	return foundValue, nil
}

func (te *TestExecutor) getValueFromDataTableRow(row *messages.PickleStepArgument_PickleTable_PickleTableRow, index int) string {
	return strings.TrimSpace(row.Cells[index].Value)
}

func (te *TestExecutor) getBoolFromDataTableRow(row *messages.PickleStepArgument_PickleTable_PickleTableRow, index int) bool {
	rowString := te.getValueFromDataTableRow(row, index)

	return strings.ToUpper(rowString) == "YES"
}

func (te *TestExecutor) getValueFromVariables(toParse string) (string, error) {
	var result strings.Builder
	tmpl, err := template.New("template").Parse(toParse)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(&result, te.Variables)
	if err != nil {
		return "", err
	}
	return result.String(), nil
}

func (te *TestExecutor) iSendARequestOn(method, path string) error {
	if method != http.MethodPost &&
		method != http.MethodPut &&
		method != http.MethodPatch &&
		method != http.MethodGet &&
		method != http.MethodDelete &&
		method != http.MethodOptions {
		return fmt.Errorf("unknown http method %s", method)
	}

	replacedBody, err := te.getValueFromVariables(te.RequestBody)
	if err != nil {
		return err
	}

	body := *strings.NewReader(replacedBody)

	url, err := te.getValueFromVariables(path)
	if err != nil {
		return err
	}

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

func (te *TestExecutor) imAuthenticatedWithThisRandomUser() error {
	te.RequestBody = fmt.Sprintf(`
    {
      "email": "%s",
      "password": "%s"
    }
  `,
		te.UserEmail,
		te.UserPassword)
	err := te.iSendARequestOn(http.MethodPost, "/account/login")
	if err != nil {
		return err
	}

	err = te.theStatusCodeIs(http.StatusOK)
	if err != nil {
		return err
	}

	accessToken, err := te.getValueFromBody("$.data.token")
	if err != nil {
		return err
	}

	te.UserToken = accessToken.(string)

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
			err := te.iSendARequestOn(http.MethodPost, "/account/signup")
			if err != nil {
				return err
			}

			err = te.theStatusCodeIs(http.StatusNoContent)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (te *TestExecutor) iCreateARandomUser() error {
	email, err := te.randomEmail()
	if err != nil {
		return err
	}

	password, err := te.randomString()
	if err != nil {
		return err
	}

	te.RequestBody = fmt.Sprintf(`
    {
      "email": "%s",
      "password": "%s"
    }
  `,
		email,
		password)

	err = te.iSendARequestOn(http.MethodPost, "/account/signup")
	if err != nil {
		return err
	}

	err = te.theStatusCodeIs(http.StatusNoContent)
	if err != nil {
		return err
	}

	te.UserEmail = email
	te.UserPassword = password

	return nil
}

func (te *TestExecutor) iCreateAndAuthenticateWithARandomUser() error {
	err := te.iCreateARandomUser()
	if err != nil {
		return err
	}

	return te.imAuthenticatedWithThisRandomUser()
}

func (te *TestExecutor) createList(l *listingPorts.CreateListRequest) error {
	te.RequestBody = fmt.Sprintf(`
    {
      "name": "%s"
    }
  `, l.Name)

	err := te.iSendARequestOn(http.MethodPost, "/listing/lists")
	if err != nil {
		return err
	}

	err = te.theStatusCodeIs(http.StatusCreated)
	if err != nil {
		return err
	}

	listURLInHeader := te.getValueFromHeader(server.HeaderContentLocation)
	if listURLInHeader == "" {
		return errors.New("list URL not set in response header")
	}

	splittedURL := strings.Split(listURLInHeader, "/")
	te.Variables.ListID = splittedURL[len(splittedURL)-1]

	return nil
}

func (te *TestExecutor) iCreateTheseLists(listDataTable *messages.PickleStepArgument_PickleTable) error {
	for i, row := range listDataTable.Rows {
		if i != 0 {
			name := strings.TrimSpace(row.Cells[0].Value)
			l := &listingPorts.CreateListRequest{
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
			p := &productsPorts.CreateProductRequest{
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

func (te *TestExecutor) createProduct(p *productsPorts.CreateProductRequest) error {
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

	productURLInHeader := te.getValueFromHeader(server.HeaderContentLocation)
	if productURLInHeader == "" {
		return errors.New("product URL not set in response header")
	}
	splittedURL := strings.Split(productURLInHeader, "/")
	te.Variables.ProductID = splittedURL[len(splittedURL)-1]

	return nil
}

func (te *TestExecutor) iCreateARandomList() error {
	name, err := te.randomString()
	if err != nil {
		return err
	}

	l := &listingPorts.CreateListRequest{
		Name: name,
	}

	err = te.createList(l)
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

func (te *TestExecutor) theStatusCodeIs(code int) error {
	if te.Response.StatusCode != code {
		return fmt.Errorf("status codes are not the same: actual %d expected %d", te.Response.StatusCode, code)
	}
	return nil
}

func (te *TestExecutor) isAStringEqualTo(path string, expected string) error {
	value, err := te.getValueFromBody(path)
	if err != nil {
		return err
	}

	expectedValue, err := te.getValueFromVariables(expected)
	if err != nil {
		return err
	}

	if value != expectedValue {
		return fmt.Errorf("actual %s should be equal to expected %s", value, expectedValue)
	}

	return nil
}

func (te *TestExecutor) theBodyIsAnEmptyArray() error {
	if string(te.ResponseBody) != "{\"data\":[]}" {
		return fmt.Errorf("the body is not an empty array, actual value \"%s\"", string(te.ResponseBody))
	}
	return nil
}

func (te *TestExecutor) isANonEmptyString(path string) error {
	value, err := te.getValueFromBody(path)
	if err != nil {
		return err
	}

	if value == "" {
		return fmt.Errorf("should not be empty")
	}

	return nil
}

func (te *TestExecutor) hasStringValue(path, expectedValue string) error {
	foundValue, err := te.getValueFromBody(path)

	if err != nil {
		return err
	}

	expectedValue, err = te.getValueFromVariables(expectedValue)
	if err != nil {
		return err
	}

	if foundValue != expectedValue {
		return fmt.Errorf("actual %s is not equal to expected %s", foundValue, expectedValue)
	}

	return nil
}

func (te *TestExecutor) hasBoolValue(path, expectedValue string) error {
	value, err := te.getValueFromBody(path)
	if err != nil {
		return err
	}

	expectedBoolValue, err := strconv.ParseBool(expectedValue)
	if err != nil {
		return err
	}

	if value.(bool) != expectedBoolValue {
		return fmt.Errorf("actual %s is not equal to expected %s for path %s", value, expectedValue, path)
	}

	return nil
}

// func (te *TestExecutor) theReturnedProductsFromListAre(productsDataTable *messages.PickleStepArgument_PickleTable) error {
// 	var list *builders.GetListResponse
// 	err := json.Unmarshal(te.ResponseBody, &list)
// 	if err != nil {
// 		return err
// 	}

// 	productsInListMap := make(map[string]*builders.ProductInGetListResponse)

// 	for _, categoryInList := range list.Categories {
// 		for _, productInList := range categoryInList.Products {
// 			productsInListMap[productInList.ID] = productInList
// 		}
// 	}

// 	expectedProductsInListLength := len(productsDataTable.Rows) - 1
// 	actualProductsInListLength := len(productsInListMap)

// 	if expectedProductsInListLength != actualProductsInListLength {
// 		return fmt.Errorf("list must contains %d products, actually contains %d", expectedProductsInListLength, actualProductsInListLength)
// 	}

// 	for i, row := range productsDataTable.Rows {
// 		if i != 0 {
// 			expectedID := te.getValueFromDataTableRow(row, 0)
// 			expectedName := te.getValueFromDataTableRow(row, 1)
// 			expectedCategoryName := te.getValueFromDataTableRow(row, 2)
// 			expectedBought := te.getBoolFromDataTableRow(row, 3)

// 			expectedID, err = te.getValueFromVariables(expectedID)
// 			if err != nil {
// 				return err
// 			}

// 			productInList, ok := productsInListMap[expectedID]
// 			if !ok {
// 				return fmt.Errorf("product %s was not found", expectedID)
// 			}

// 			if productInList.Name != expectedName {
// 				return fmt.Errorf("product name %s was not expected %s", productInList.Name, expectedName)
// 			}

// 			productMessage := fmt.Sprintf("for product %s", expectedName)

// 			if productInList.Bought != expectedBought {
// 				return fmt.Errorf("product bought %t was not expected %t %s", productInList.Bought, expectedBought, productMessage)
// 			}

// 			if productInList.Category != nil && productInList.Category.Name != expectedCategoryName {
// 				return fmt.Errorf("product category name %s was not expected %s %s", productInList.Category.Name, expectedCategoryName, productMessage)
// 			}
// 		}
// 	}

// 	return nil
// }

func (te *TestExecutor) randomEmail() (string, error) {
	email, err := te.Generation.UUID()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s@crounch.me", email), nil
}

func (te *TestExecutor) randomString() (string, error) {
	return te.Generation.UUID()
}

func FeatureContext(s *godog.Suite) {
	te := &TestExecutor{
		Variables: ExecutorVariables{
			ListID:    "",
			ProductID: "",
		},
		Generation: utils.NewGeneration(),
	}

	// Requests
	s.Step(`^I use this body$`, te.iUseThisBody)
	s.Step(`^I use an invalid body$`, te.iUseAnInvalidBody)
	s.Step(`^I use an empty valid body$`, te.iUseAnEmptyValidBody)
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
	s.Step(`^I\'m authenticated with this random user$`, te.imAuthenticatedWithThisRandomUser)
	s.Step(`^I authenticate with a random user$`, te.iCreateAndAuthenticateWithARandomUser)
	s.Step(`^I create a random user$`, te.iCreateARandomUser)

	// Users
	s.Step(`^I create these users?$`, te.iCreateTheseUsers)

	// Lists
	s.Step(`^I create these lists$`, te.iCreateTheseLists)
	s.Step(`^I create a random list$`, te.iCreateARandomList)
	// s.Step(`^the returned products from list are$`, te.theReturnedProductsFromListAre)

	// Products
	s.Step(`^I create these products$`, te.iCreateTheseProducts)
}
