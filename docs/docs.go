// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "vangxitrum",
            "url": "http://www.swagger.io/support",
            "email": "19522482@gm.uit.edu.vn"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/sign-in": {
            "post": {
                "description": "Login and set access token to header",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Login",
                "parameters": [
                    {
                        "description": "User's credentials",
                        "name": "SignInModel",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.SignInModel"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/auth/sign-up": {
            "post": {
                "description": "Endpoint to allow a user to sign up with their details",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "sign up user",
                "parameters": [
                    {
                        "description": "User's profile",
                        "name": "SignUpModel",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.SignUpModel"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "success",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "fail",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/auth/verify": {
            "get": {
                "description": "Verify user email using verification code",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Verify user email",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "unique token",
                        "name": "uniqueToken",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/cart": {
            "put": {
                "description": "add item to cart, return info of added item",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Cart"
                ],
                "summary": "add item to cart",
                "parameters": [
                    {
                        "description": "Add item request",
                        "name": "CartRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.AddItemRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.CartItem"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "Update cart item by delete all old cart items, then add received item to customer card",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Cart"
                ],
                "summary": "update cart of customer by received cart item",
                "parameters": [
                    {
                        "description": "access token received after login",
                        "name": "CartRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.UpdateCartRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete cart item by product id of customer",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Cart"
                ],
                "summary": "delete card item of customer cart",
                "parameters": [
                    {
                        "type": "string",
                        "description": "customer's id",
                        "name": "customer_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "product's id",
                        "name": "product_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/cart/checkout": {
            "get": {
                "description": "Check out customer cart then delete invalid item(sold out item) return list of sold-out Items' ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Cart"
                ],
                "summary": "Use to validate items of cart then modify it if invalid",
                "parameters": [
                    {
                        "type": "string",
                        "description": "customer's id",
                        "name": "customer_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/product/{id}": {
            "get": {
                "description": "get the product info",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "product"
                ],
                "summary": "get product info",
                "parameters": [
                    {
                        "type": "string",
                        "description": "product's id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Product"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/products/": {
            "get": {
                "description": "list of the product",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "product"
                ],
                "summary": "list of product",
                "parameters": [
                    {
                        "type": "string",
                        "description": "a list of brand name separated by commas",
                        "name": "brands",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "a list of color name separated by commas (FULL UPPERCASE format)",
                        "name": "colors",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "a list of tag name ['HOT','NEW','SALE'] separated by commas",
                        "name": "tags",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "a list of gender type ['KID','WOMEN','MEN'] separated by commas",
                        "name": "genders",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "a list of type name separated by commas",
                        "name": "types",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Minimum of avg rate of product",
                        "name": "rate",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Range of values in format 'min_value,max_value' ",
                        "name": "price",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Key work relate to products' name ",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "current page's number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Length per page from '1' to '10000'",
                        "name": "page_size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.ListProductResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/users/me": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "get the current user info",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "get user info",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.CartItem": {
            "type": "object",
            "properties": {
                "color": {
                    "type": "string"
                },
                "customer_id": {
                    "type": "string"
                },
                "product_detail": {
                    "$ref": "#/definitions/model.Product"
                },
                "product_id": {
                    "type": "string"
                },
                "quantity": {
                    "type": "integer"
                },
                "size": {
                    "type": "string"
                }
            }
        },
        "model.Product": {
            "type": "object",
            "properties": {
                "avr_rate": {
                    "type": "number"
                },
                "brand": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "discount_amount": {
                    "type": "number"
                },
                "discount_percent": {
                    "type": "number"
                },
                "gender": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "photos": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.ProductPhoto"
                    }
                },
                "price": {
                    "type": "integer"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "types": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "model.ProductPhoto": {
            "type": "object",
            "properties": {
                "color": {
                    "type": "string"
                },
                "mainPhoto": {
                    "type": "string"
                },
                "productId": {
                    "type": "string"
                },
                "subPhotos": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "model.SignInModel": {
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
                    "type": "string",
                    "minLength": 8
                }
            }
        },
        "model.SignUpModel": {
            "type": "object",
            "required": [
                "email",
                "fullname",
                "password",
                "password_confirm"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "fullname": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "minLength": 8
                },
                "password_confirm": {
                    "type": "string"
                }
            }
        },
        "model.User": {
            "type": "object",
            "properties": {
                "avatar": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "fullname": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "phoneNumber": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                },
                "verified": {
                    "type": "boolean"
                }
            }
        },
        "request.AddItemRequest": {
            "type": "object",
            "required": [
                "customer_id",
                "productId",
                "quantity"
            ],
            "properties": {
                "customer_id": {
                    "type": "string"
                },
                "productId": {
                    "type": "string"
                },
                "quantity": {
                    "type": "integer"
                }
            }
        },
        "request.UpdateCartRequest": {
            "type": "object",
            "required": [
                "customer_id",
                "items"
            ],
            "properties": {
                "customer_id": {
                    "type": "string"
                },
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/request.UpdateItem"
                    }
                }
            }
        },
        "request.UpdateItem": {
            "type": "object",
            "required": [
                "productId",
                "quantity"
            ],
            "properties": {
                "productId": {
                    "type": "string"
                },
                "quantity": {
                    "type": "integer"
                }
            }
        },
        "response.ListProductResponse": {
            "type": "object",
            "properties": {
                "page": {
                    "type": "integer"
                },
                "products": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Product"
                    }
                },
                "total_page": {
                    "type": "integer"
                }
            }
        }
    },
    "securityDefinitions": {
        "BasicAuth": {
            "type": "basic"
        }
    },
    "externalDocs": {
        "description": "OpenAPI",
        "url": "https://swagger.io/resources/open-api/"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "127.0.0.1:8081",
	BasePath:         "/api/",
	Schemes:          []string{},
	Title:            "Online Shop API",
	Description:      "Online Shop API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
