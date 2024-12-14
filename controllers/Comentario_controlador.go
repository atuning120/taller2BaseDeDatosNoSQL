package controllers

import (
    "go-API/models"
    "go-API/services"
    "go-API/response" // Importar el paquete donde está ErrorResponse
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

type ComentarioControlador struct {
    servicio *services.ComentarioService
}

// NewComentarioControlador crea un nuevo controlador para los comentarios.
func NewComentarioControlador(servicio *services.ComentarioService) *ComentarioControlador {
    return &ComentarioControlador{servicio: servicio}
}

// ObtenerComentariosPorClase obtiene todos los comentarios de una clase.
// @Summary Devuelve los comentarios de una clase
// @Description Devuelve todos los comentarios asociados a una clase por su ID
// @Tags Comentarios
// @Accept json
// @Produce json
// @Param id path string true "ID de la clase"
// @Success 200 {array} models.Comentario "Lista de comentarios"
// @Failure 500 {object} response.ErrorResponse "Error interno del servidor"
// @Router /api/clases/{id}/comentarios [get]
func (c *ComentarioControlador) ObtenerComentariosPorClase(ctx *gin.Context) {
    claseID := ctx.Param("id")
    comentarios, err := c.servicio.ObtenerComentariosPorClase(ctx.Request.Context(), claseID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, comentarios)
}

// CrearComentarioParaClase crea un nuevo comentario para una clase.
// @Summary Crear un comentario para una clase
// @Description Agrega un comentario a una clase por su ID
// @Tags Comentarios
// @Accept json
// @Produce json
// @Param id path string true "ID de la clase"
// @Param comentario body models.Comentario true "Comentario a crear"
// @Success 201 {object} models.Comentario "Comentario creado exitosamente"
// @Failure 400 {object} response.ErrorResponse "Datos inválidos o faltan campos requeridos"
// @Failure 404 {object} response.ErrorResponse "Clase no encontrada"
// @Failure 500 {object} response.ErrorResponse "Error interno del servidor"
// @Router /api/clases/{id}/comentarios [post]
func (c *ComentarioControlador) CrearComentarioParaClase(ctx *gin.Context) {
    claseID := ctx.Param("id")
    var comentario models.Comentario

    // Validar los datos enviados en la solicitud
    if err := ctx.ShouldBindJSON(&comentario); err != nil {
        ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Datos inválidos: " + err.Error()})
        return
    }

    // Validar campos requeridos (por ejemplo, autor, titulo, detalle)
    if comentario.Autor == "" || comentario.Titulo == "" || comentario.Detalle == "" {
        ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Los campos autor, titulo y detalle son requeridos"})
        return
    }

    // Asignar fecha actual antes de crear (el servicio ya asigna, pero aquí puedes controlar)
    comentario.Fecha = time.Now()

    creado, err := c.servicio.CrearComentarioParaClase(ctx.Request.Context(), claseID, &comentario)
    if err != nil {
        if err.Error() == "clase no encontrada" {
            ctx.JSON(http.StatusNotFound, response.ErrorResponse{Message: err.Error()})
        } else {
            ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
        }
        return
    }

    ctx.JSON(http.StatusCreated, creado)
}
