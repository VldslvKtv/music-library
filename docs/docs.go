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
        "/get_data/songs": {
            "get": {
                "description": "Получение данных библиотеки с фильтрацией по всем полям и пагинацией (метод GET).",
                "produces": [
                    "application/json"
                ],
                "summary": "Получение данных библиотеки",
                "operationId": "get-all-data",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Имя группы",
                        "name": "group",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Имя песни",
                        "name": "song",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Дата релиза",
                        "name": "releaseDate",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Текст песни",
                        "name": "text",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Ссылка на песню",
                        "name": "link",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Номер страницы",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Размер страницы",
                        "name": "pageSize",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/get_all_data.Response"
                        }
                    },
                    "500": {
                        "description": "failed to get songs",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/get_data/text": {
            "get": {
                "description": "Получение текста песни с пагинацией по куплетам (метод GET).",
                "produces": [
                    "application/json"
                ],
                "summary": "Получение текста песни",
                "operationId": "get-song",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Имя группы",
                        "name": "group",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Имя песни",
                        "name": "song",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Номер страницы",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Размер страницы",
                        "name": "pageSize",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/get_song.Response"
                        }
                    },
                    "400": {
                        "description": "group and song parameters are required or any other errors",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "failed to get song",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/songs/": {
            "post": {
                "description": "Добавление новой песни в формате JSON и запросе к внешнему API.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Добавление новой песни",
                "operationId": "add-song",
                "parameters": [
                    {
                        "description": "Данные песни",
                        "name": "song",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.SongAndGroup"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Data"
                        }
                    },
                    "400": {
                        "description": "failed to decode req-body or any other errors",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/songs/{id}": {
            "delete": {
                "description": "Удаление песни по ID.",
                "produces": [
                    "application/json"
                ],
                "summary": "Удаление песни",
                "operationId": "delete-song",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID песни",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "invalid ID or any other errors",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "patch": {
                "description": "Изменение данных песни по ID.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Изменение данных песни",
                "operationId": "update-song",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID песни",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Данные песни",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Data"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "failed to decode req-body or any other errors",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "get_all_data.Response": {
            "description": "Структура ответа с данными песен и информацией о пагинации.",
            "type": "object",
            "properties": {
                "currentPage": {
                    "type": "integer"
                },
                "maxPageSize": {
                    "type": "integer"
                },
                "songs": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Data"
                    }
                },
                "totalPages": {
                    "type": "integer"
                },
                "totalSongs": {
                    "type": "integer"
                }
            }
        },
        "get_song.Response": {
            "description": "Структура ответа с текстом песни и информацией о куплетах.",
            "type": "object",
            "properties": {
                "currentPage": {
                    "type": "integer"
                },
                "maxVersesPerPage": {
                    "type": "integer"
                },
                "totalPages": {
                    "type": "integer"
                },
                "verses": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "models.CustomTime": {
            "type": "object",
            "properties": {
                "time.Time": {
                    "type": "string"
                }
            }
        },
        "models.Data": {
            "type": "object",
            "properties": {
                "group": {
                    "type": "string"
                },
                "link": {
                    "type": "string"
                },
                "releaseDate": {
                    "$ref": "#/definitions/models.CustomTime"
                },
                "song": {
                    "type": "string"
                },
                "text": {
                    "type": "string"
                }
            }
        },
        "models.SongAndGroup": {
            "type": "object",
            "properties": {
                "group": {
                    "type": "string"
                },
                "song": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8002",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Music Library API",
	Description:      "This is a sample server for a music library.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
