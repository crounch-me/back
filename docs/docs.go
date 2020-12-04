// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/health": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Return health of application",
                "operationId": "get-health",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/internal.Health"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/lists": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "list"
                ],
                "summary": "Get the lists of the owner",
                "operationId": "get-owners-lists",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/list.List"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/internal.Error"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "list"
                ],
                "summary": "Create a list",
                "operationId": "create-list",
                "parameters": [
                    {
                        "description": "List to create",
                        "name": "list",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.CreateListRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/list.List"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/internal.Error"
                        }
                    }
                }
            }
        },
        "/lists/{listID}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "list"
                ],
                "summary": "Reads a list with products in categories",
                "operationId": "get-list",
                "parameters": [
                    {
                        "type": "string",
                        "description": "List ID",
                        "name": "listID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/builders.GetListResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/internal.Error"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "list"
                ],
                "summary": "Delete the entire list with its products",
                "operationId": "delete-list",
                "parameters": [
                    {
                        "type": "string",
                        "description": "List ID",
                        "name": "listID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {},
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/internal.Error"
                        }
                    }
                }
            }
        },
        "/lists/{listID}/products/{productID}": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "product-in-list"
                ],
                "summary": "Add the product to the list",
                "operationId": "add-product-to-list",
                "parameters": [
                    {
                        "type": "string",
                        "description": "List ID",
                        "name": "listID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Product ID",
                        "name": "productID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/list.ProductInListLink"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/internal.Error"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "product-in-list"
                ],
                "summary": "Delete the product from the list",
                "operationId": "delete-product-from-list",
                "parameters": [
                    {
                        "type": "string",
                        "description": "List ID",
                        "name": "listID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Product ID",
                        "name": "productID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {},
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/internal.Error"
                        }
                    }
                }
            },
            "patch": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "product-in-list"
                ],
                "summary": "Update the product in the list partially",
                "operationId": "update-product-in-list",
                "parameters": [
                    {
                        "description": "Product in list",
                        "name": "productInList",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/list.UpdateProductInList"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Product in list",
                        "name": "listID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Product in list",
                        "name": "productID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/list.ProductInListLink"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/internal.Error"
                        }
                    }
                }
            }
        },
        "/logout": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "tags": [
                    "user"
                ],
                "summary": "Removes an user authorization",
                "operationId": "logout",
                "responses": {
                    "204": {},
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/internal.Error"
                        }
                    }
                }
            }
        },
        "/me": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Removes an user authorization",
                "operationId": "me",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/users.User"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/internal.Error"
                        }
                    }
                }
            }
        },
        "/products": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "product"
                ],
                "summary": "Create a new product",
                "operationId": "create-product",
                "parameters": [
                    {
                        "description": "Product to create",
                        "name": "product",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.CreateProductRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/products.Product"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/internal.Error"
                        }
                    }
                }
            }
        },
        "/products/search": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "product"
                ],
                "summary": "Search a product by its name in default products, it removes accentuated characters and is case insensitive",
                "operationId": "search-default-products",
                "parameters": [
                    {
                        "description": "Product search request",
                        "name": "product",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.ProductSearchRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/products.Product"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/internal.Error"
                        }
                    }
                }
            }
        },
        "/users/login": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Creates a new user authorization",
                "operationId": "login",
                "parameters": [
                    {
                        "description": "User to login with",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/authorization.Authorization"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/internal.Error"
                        }
                    }
                }
            }
        },
        "/users/signup": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Creates a new user",
                "operationId": "signup",
                "parameters": [
                    {
                        "description": "User to signup with",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.SignupRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/account.User"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/internal.Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "account.User": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "authorization.Authorization": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string"
                },
                "owner": {
                    "$ref": "#/definitions/account.User"
                }
            }
        },
        "builders.CategoryInGetListResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "products": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/builders.ProductInGetListResponse"
                    }
                }
            }
        },
        "builders.GetListResponse": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "archivationDate": {
                    "type": "string"
                },
                "categories": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/builders.CategoryInGetListResponse"
                    }
                },
                "contributors": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/account.User"
                    }
                },
                "creationDate": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "builders.ProductInGetListResponse": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "bought": {
                    "type": "boolean"
                },
                "category": {
                    "$ref": "#/definitions/categories.Category"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "owner": {
                    "$ref": "#/definitions/account.User"
                }
            }
        },
        "categories.Category": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "handler.CreateListRequest": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "handler.CreateProductRequest": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "handler.LoginRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "handler.ProductSearchRequest": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "handler.SignupRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "internal.Error": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "fields": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/internal.FieldError"
                    }
                }
            }
        },
        "internal.FieldError": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "internal.Health": {
            "type": "object",
            "properties": {
                "alive": {
                    "type": "boolean"
                }
            }
        },
        "list.List": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "archivationDate": {
                    "type": "string"
                },
                "contributors": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/account.User"
                    }
                },
                "creationDate": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "products": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/list.ProductInList"
                    }
                }
            }
        },
        "list.ProductInList": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "bought": {
                    "type": "boolean"
                },
                "category": {
                    "$ref": "#/definitions/categories.Category"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "owner": {
                    "$ref": "#/definitions/account.User"
                }
            }
        },
        "list.ProductInListLink": {
            "type": "object",
            "properties": {
                "bought": {
                    "type": "boolean"
                },
                "listId": {
                    "type": "string"
                },
                "productId": {
                    "type": "string"
                }
            }
        },
        "list.UpdateProductInList": {
            "type": "object",
            "properties": {
                "bought": {
                    "type": "boolean"
                }
            }
        },
        "products.Product": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "category": {
                    "$ref": "#/definitions/categories.Category"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "owner": {
                    "$ref": "#/definitions/account.User"
                }
            }
        },
        "users.User": {
            "type": "object"
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "localhost:3000",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "Crounch Me API",
	Description: "API serving the grocery application.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
