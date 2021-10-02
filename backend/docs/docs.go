// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Sergey \"hackfeed\" Kononenko",
            "url": "https://hackfeed.github.io",
            "email": "hackfeed@yandex.ru"
        },
        "license": {
            "name": "GPL-3.0 License",
            "url": "http://www.gnu.org/licenses/gpl-3.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/analytics/mrr": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Getting MRR analytics data with all components for given period",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "analytics"
                ],
                "summary": "Get MRR analytics data",
                "parameters": [
                    {
                        "description": "Parameters for MRR analytics",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Period"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ResponseSuccessAnalytics"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Response"
                        }
                    }
                }
            }
        },
        "/auth/login": {
            "post": {
                "description": "Logging user in by retrieving his data from the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Logging user in",
                "parameters": [
                    {
                        "description": "User's email and password",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ResponseSuccessAuth"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Response"
                        }
                    }
                }
            }
        },
        "/auth/signup": {
            "post": {
                "description": "Signing user up by adding him to the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Signing user up",
                "parameters": [
                    {
                        "description": "User's email and password",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ResponseSuccessAuth"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Response"
                        }
                    }
                }
            }
        },
        "/files/delete/{filename}": {
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Deleting file and cleaning database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "files"
                ],
                "summary": "Deleting user's file",
                "parameters": [
                    {
                        "type": "string",
                        "description": "File to delete",
                        "name": "filename",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Response"
                        }
                    }
                }
            }
        },
        "/files/load": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Loading files' names, uploaded by user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "files"
                ],
                "summary": "Loading user's files",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ResponseSuccessLoadFiles"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Response"
                        }
                    }
                }
            }
        },
        "/files/upload": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Saving file locally on the server and parsing its content to database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "files"
                ],
                "summary": "Saving user's file",
                "parameters": [
                    {
                        "type": "file",
                        "description": "File to upload",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ResponseSuccessSaveFile"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.File": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "uploaded_at": {
                    "type": "string"
                }
            }
        },
        "domain.TotalMRR": {
            "type": "object",
            "properties": {
                "churn": {
                    "type": "array",
                    "items": {
                        "type": "number"
                    }
                },
                "contraction": {
                    "type": "array",
                    "items": {
                        "type": "number"
                    }
                },
                "expansion": {
                    "type": "array",
                    "items": {
                        "type": "number"
                    }
                },
                "new": {
                    "type": "array",
                    "items": {
                        "type": "number"
                    }
                },
                "old": {
                    "type": "array",
                    "items": {
                        "type": "number"
                    }
                },
                "reactivation": {
                    "type": "array",
                    "items": {
                        "type": "number"
                    }
                },
                "total": {
                    "type": "array",
                    "items": {
                        "type": "number"
                    }
                }
            }
        },
        "models.Period": {
            "type": "object",
            "required": [
                "filename",
                "period_end",
                "period_start"
            ],
            "properties": {
                "filename": {
                    "type": "string",
                    "example": "filename.csv"
                },
                "period_end": {
                    "type": "string",
                    "example": "2021-01-01"
                },
                "period_start": {
                    "type": "string",
                    "example": "2019-01-01"
                }
            }
        },
        "models.Response": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "models.ResponseSuccessAnalytics": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Analytics is loaded"
                },
                "months": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "mrr": {
                    "$ref": "#/definitions/domain.TotalMRR"
                }
            }
        },
        "models.ResponseSuccessAuth": {
            "type": "object",
            "properties": {
                "expires_at": {
                    "type": "integer"
                },
                "id_token": {
                    "type": "string"
                },
                "local_id": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "models.ResponseSuccessLoadFiles": {
            "type": "object",
            "properties": {
                "files": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.File"
                    }
                },
                "message": {
                    "type": "string",
                    "example": "Files are loaded"
                }
            }
        },
        "models.ResponseSuccessSaveFile": {
            "type": "object",
            "properties": {
                "filename": {
                    "type": "string",
                    "example": "filename.csv"
                },
                "message": {
                    "type": "string",
                    "example": "File is uploaded"
                }
            }
        },
        "models.User": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "test@test.com"
                },
                "password": {
                    "type": "string",
                    "example": "password123"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "token",
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
	Host:        "weblabs.com:8003",
	BasePath:    "/api/v1",
	Schemes:     []string{},
	Title:       "remrratality API",
	Description: "API for getting MRR analytics of your app's money flow.",
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
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
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
