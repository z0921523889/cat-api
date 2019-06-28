// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2019-06-28 14:19:34.6699179 +0800 CST m=+0.959974101

package docs

import (
	"bytes"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/cat": {
            "post": {
                "description": "Add a new cat",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Add new cat to the database",
                "parameters": [
                    {
                        "type": "string",
                        "description": "貓的名稱",
                        "name": "name",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "貓的級別",
                        "name": "level",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "貓的價格",
                        "name": "price",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "貓的pet幣",
                        "name": "pet_coin",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "貓的預約價格",
                        "name": "reservation_price",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "貓的即搶價格",
                        "name": "adoption_price",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "貓的合約時間",
                        "name": "contract_days",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "貓的合約增益",
                        "name": "contract_benefit",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/controller.Message"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controller.Message": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "message"
                }
            }
        },
        "httputil.HTTPError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 400
                },
                "message": {
                    "type": "string",
                    "example": "status bad request"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo swaggerInfo

type s struct{}

func (s *s) ReadDoc() string {
	t, err := template.New("swagger_info").Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, SwaggerInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}