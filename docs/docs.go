// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/clases/{id}/comentarios": {
            "get": {
                "description": "Devuelve todos los comentarios asociados a una clase por su ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Comentarios"
                ],
                "summary": "Devuelve los comentarios de una clase",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID de la clase",
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
                                "$ref": "#/definitions/response.ComentarioResponse"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Agrega un comentario a una clase por su ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Comentarios"
                ],
                "summary": "Crear un comentario para una clase",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID de la clase",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Comentario a crear",
                        "name": "comentario",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.CreateComentarioRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/response.ComentarioResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/cursos": {
            "get": {
                "description": "Devuelve todos los cursos disponibles",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Cursos"
                ],
                "summary": "Devuelve todos los cursos",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/response.CursoResponse"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Agrega un curso a la base de datos",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Cursos"
                ],
                "summary": "Crear un curso",
                "parameters": [
                    {
                        "description": "Curso a crear",
                        "name": "curso",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.CreateCursoRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.CrearCurso"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/cursos/{id}": {
            "get": {
                "description": "Devuelve un curso en específico dado su ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Cursos"
                ],
                "summary": "Devuelve un curso según su ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID del curso",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.CursoResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/cursos/{id}/unidades": {
            "get": {
                "description": "Devuelve una unidades de un curso en específico dado su ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Unidades"
                ],
                "summary": "Devuelve unidades de un curso",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID del curso",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.CursoResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "agregar una unidad a un curso",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Unidades"
                ],
                "summary": "Crear unidad",
                "parameters": [
                    {
                        "description": "Unidad a crear",
                        "name": "unidad",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.CreateUnidadRequest"
                        }
                    },
                    {
                        "type": "string",
                        "description": "ID del curso",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.CrearUnidad"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/cursos/{id}/valoracion": {
            "patch": {
                "description": "Actualiza la valoración de un curso según la nueva valoración proporcionada",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Cursos"
                ],
                "summary": "Actualiza la valoración de un curso",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID del curso",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Nueva valoración del curso",
                        "name": "valoracion",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.UpdateValoracionRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.UpdateValoracionResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/unidades/{id}/clases": {
            "get": {
                "description": "Devuelve todas las clases asociadas a una unidad",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Clases"
                ],
                "summary": "Devuelve las clases de una unidad",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID de la unidad",
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
                                "$ref": "#/definitions/response.ClaseResponse"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Agrega una clase a la base de datos asociada a una unidad",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Clases"
                ],
                "summary": "Crear una clase para una unidad",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID de la unidad",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Clase a crear",
                        "name": "clase",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.CreateClaseRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.CrearClase"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/usuarios": {
            "get": {
                "description": "Devuelve la lista completa de usuarios registrados",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Usuarios"
                ],
                "summary": "Devuelve todos los usuarios",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Usuario"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Agrega un usuario a la base de datos",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Usuarios"
                ],
                "summary": "Crear un nuevo usuario",
                "parameters": [
                    {
                        "description": "Usuario a crear",
                        "name": "usuario",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.CreateUsuarioRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.Usuario"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/usuarios/cursos": {
            "get": {
                "description": "Devuelve la lista de cursos en los que un usuario está inscrito",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Usuarios"
                ],
                "summary": "Obtener cursos inscritos de un usuario",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Correo del usuario",
                        "name": "email",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Contraseña del usuario",
                        "name": "password",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Curso"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/usuarios/inscripcion": {
            "post": {
                "description": "Inscribe a un usuario en un curso específico",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Usuarios"
                ],
                "summary": "Inscribir un usuario en un curso",
                "parameters": [
                    {
                        "description": "Datos de inscripción",
                        "name": "inscripcion",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.InscripcionRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.InscripcionResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/usuarios/usuario": {
            "get": {
                "description": "Devuelve un usuario en específico por su correo y contraseña",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Usuarios"
                ],
                "summary": "Obtener un usuario por correo y contraseña",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Correo del usuario",
                        "name": "email",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Contraseña del usuario",
                        "name": "password",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Usuario"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/response.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Curso": {
            "type": "object",
            "properties": {
                "cant_usuarios": {
                    "type": "integer"
                },
                "comentarios": {
                    "description": "Lista de IDs de comentarios",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "descripcion": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "imagen_url": {
                    "type": "string"
                },
                "nombre": {
                    "type": "string"
                },
                "unidades": {
                    "description": "Lista de IDs de unidades",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "valoracion": {
                    "type": "number"
                }
            }
        },
        "models.Usuario": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "inscritos": {
                    "description": "IDs de cursos inscritos",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "nombre": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "request.CreateClaseRequest": {
            "type": "object",
            "required": [
                "descripcion",
                "nombre",
                "video_url"
            ],
            "properties": {
                "descripcion": {
                    "type": "string"
                },
                "nombre": {
                    "type": "string"
                },
                "video_url": {
                    "type": "string"
                }
            }
        },
        "request.CreateComentarioRequest": {
            "type": "object",
            "required": [
                "autor",
                "detalle",
                "me_gusta",
                "no_me_gusta",
                "titulo"
            ],
            "properties": {
                "autor": {
                    "type": "string"
                },
                "detalle": {
                    "type": "string"
                },
                "me_gusta": {
                    "type": "integer"
                },
                "no_me_gusta": {
                    "type": "integer"
                },
                "titulo": {
                    "type": "string"
                }
            }
        },
        "request.CreateCursoRequest": {
            "type": "object",
            "required": [
                "nombre"
            ],
            "properties": {
                "descripcion": {
                    "type": "string"
                },
                "imagen_url": {
                    "type": "string"
                },
                "nombre": {
                    "type": "string"
                }
            }
        },
        "request.CreateUnidadRequest": {
            "type": "object",
            "required": [
                "nombre"
            ],
            "properties": {
                "nombre": {
                    "type": "string"
                }
            }
        },
        "request.CreateUsuarioRequest": {
            "type": "object",
            "required": [
                "email",
                "nombre",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "nombre": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "request.InscripcionRequest": {
            "type": "object",
            "required": [
                "curso_id",
                "email",
                "password"
            ],
            "properties": {
                "curso_id": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "request.UpdateValoracionRequest": {
            "type": "object",
            "required": [
                "valoracion"
            ],
            "properties": {
                "valoracion": {
                    "type": "number"
                }
            }
        },
        "response.ClaseResponse": {
            "type": "object",
            "properties": {
                "adjuntos_url": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "comentarios": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "descripcion": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "me_gusta": {
                    "type": "integer"
                },
                "no_me_gusta": {
                    "type": "integer"
                },
                "nombre": {
                    "type": "string"
                },
                "unidad_id": {
                    "type": "string"
                },
                "video_url": {
                    "type": "string"
                }
            }
        },
        "response.ComentarioResponse": {
            "type": "object",
            "properties": {
                "autor": {
                    "type": "string"
                },
                "clase_id": {
                    "type": "string"
                },
                "detalle": {
                    "type": "string"
                },
                "fecha": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "me_gusta": {
                    "type": "integer"
                },
                "no_me_gusta": {
                    "type": "integer"
                },
                "titulo": {
                    "type": "string"
                }
            }
        },
        "response.CrearClase": {
            "type": "object",
            "properties": {
                "inserted_id": {
                    "type": "string"
                }
            }
        },
        "response.CrearCurso": {
            "type": "object",
            "properties": {
                "inserted_id": {
                    "type": "string"
                }
            }
        },
        "response.CrearUnidad": {
            "type": "object",
            "properties": {
                "inserted_id": {
                    "type": "string"
                }
            }
        },
        "response.CursoResponse": {
            "type": "object",
            "properties": {
                "cant_usuarios": {
                    "type": "integer"
                },
                "comentarios": {
                    "description": "IDs de los comentarios",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "descripcion": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "imagen_url": {
                    "type": "string"
                },
                "nombre": {
                    "type": "string"
                },
                "unidades": {
                    "description": "IDs de las unidades",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "valoracion": {
                    "type": "number"
                }
            }
        },
        "response.ErrorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "response.InscripcionResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "response.UpdateValoracionResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "valoracion_actualizada": {
                    "type": "number"
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
	Title:            "API de Cursos y Usuarios",
	Description:      "Esta es una API para gestionar cursos y usuarios.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
