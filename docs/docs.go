// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
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
        "/ltp": {
            "get": {
                "description": "ltp",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "LTPHandler"
                ],
                "summary": "ltp",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.LTPResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.LTP": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "string"
                },
                "pair": {
                    "type": "string"
                }
            }
        },
        "dto.LTPResponse": {
            "type": "object",
            "properties": {
                "ltp": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.LTP"
                    }
                }
            }
        }
    },
    "securityDefinitions": {
        "Authorization": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        },
        "X-API-Key": {
            "type": "apiKey",
            "name": "X-API-Key",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "LTP API",
	Description:      "LTP API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
