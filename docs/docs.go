// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Gabriel Teixeira",
            "url": "https://github.com/gfteix"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/books": {
            "get": {
                "description": "Retrieves a list of books with optional filters",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "books"
                ],
                "summary": "Get books with filters",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Filter by title",
                        "name": "title",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by author",
                        "name": "author",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by ISBN",
                        "name": "isbn",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/types.Book"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.APIError"
                        }
                    }
                }
            },
            "post": {
                "description": "Adds a new book to the library system",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "books"
                ],
                "summary": "Create a new book",
                "parameters": [
                    {
                        "description": "Book details",
                        "name": "book",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.CreateBookPayload"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/types.APIError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.APIError"
                        }
                    }
                }
            }
        },
        "/books/{bookId}/items/{itemId}": {
            "get": {
                "description": "Retrieves a specific book item by its ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "books"
                ],
                "summary": "Get a book item by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Book ID",
                        "name": "bookId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Book Item ID",
                        "name": "itemId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.BookCopy"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/types.APIError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.APIError"
                        }
                    }
                }
            }
        },
        "/books/{id}": {
            "get": {
                "description": "Retrieves a book by its ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "books"
                ],
                "summary": "Get a book by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Book ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.Book"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/types.APIError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.APIError"
                        }
                    }
                }
            }
        },
        "/books/{id}/items": {
            "get": {
                "description": "Retrieves all items belonging to a book by its ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "books"
                ],
                "summary": "Get items of a book",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Book ID",
                        "name": "id",
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
                                "$ref": "#/definitions/types.BookCopy"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.APIError"
                        }
                    }
                }
            },
            "post": {
                "description": "Adds a new book item to a book",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "books"
                ],
                "summary": "Create a book item",
                "parameters": [
                    {
                        "description": "Book item details",
                        "name": "bookCopy",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.CreateBookCopyPayload"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Book ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/types.APIError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.APIError"
                        }
                    }
                }
            }
        },
        "/loans": {
            "get": {
                "description": "Retrieves loans with optional filters",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "loans"
                ],
                "summary": "Get Loans",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Filter by User ID",
                        "name": "userId",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by Loan Status",
                        "name": "status",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by Book Item ID",
                        "name": "bookCopyId",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/types.Loan"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/types.APIError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.APIError"
                        }
                    }
                }
            },
            "post": {
                "description": "Creates a book loan",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "loans"
                ],
                "summary": "Creates a Loan",
                "parameters": [
                    {
                        "description": "Loan that needs to be created",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.CreateLoanPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/types.APIError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.APIError"
                        }
                    }
                }
            }
        },
        "/loans/{id}": {
            "get": {
                "description": "Retrieves a loan by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "loans"
                ],
                "summary": "Get a Loan",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Loan ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.Loan"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/types.APIError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.APIError"
                        }
                    }
                }
            }
        },
        "/users": {
            "post": {
                "description": "Creates an User",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Creates an User",
                "parameters": [
                    {
                        "description": "User object that needs to be created",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.CreateUserPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/types.APIError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.APIError"
                        }
                    }
                }
            }
        },
        "/users/": {
            "get": {
                "description": "Retrieves users",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get users",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/types.User"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/types.APIError"
                        }
                    }
                }
            }
        },
        "/users/{id}": {
            "get": {
                "description": "Retrieves user details by their unique ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get a user by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/types.APIError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/types.APIError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "types.APIError": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "types.Book": {
            "type": "object",
            "properties": {
                "author": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "isbn": {
                    "type": "string"
                },
                "numberOfPages": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "types.BookCopy": {
            "type": "object",
            "properties": {
                "bookId": {
                    "type": "string"
                },
                "condition": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "location": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "types.CreateBookCopyPayload": {
            "type": "object",
            "properties": {
                "bookId": {
                    "type": "string"
                },
                "condition": {
                    "type": "string"
                },
                "location": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "types.CreateBookPayload": {
            "type": "object",
            "properties": {
                "author": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "isbn": {
                    "type": "string"
                },
                "numberOfPages": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "types.CreateLoanPayload": {
            "type": "object",
            "properties": {
                "bookCopyId": {
                    "type": "string"
                },
                "expiringDate": {
                    "type": "string"
                },
                "loanDate": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "userId": {
                    "type": "string"
                }
            }
        },
        "types.CreateUserPayload": {
            "type": "object",
            "required": [
                "email",
                "name"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "types.Loan": {
            "type": "object",
            "properties": {
                "bookCopyId": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "expiringDate": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "loanDate": {
                    "type": "string"
                },
                "returnDate": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "userId": {
                    "type": "string"
                }
            }
        },
        "types.User": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Book Loan API",
	Description:      "API to manage book loans",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
