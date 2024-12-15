package controllers

import (
    "go-API/models"
    "go-API/services"
    "go-API/response"
    "net/http"

    "github.com/gin-gonic/gin"
)

type ComentarioControlador struct {
    servicio *services.ComentarioService
}

func NewComentarioControlador(servicio *services.ComentarioService) *ComentarioControlador {
    return &ComentarioControlador{servicio: servicio}
}

// ObtenerComentariosPorClase
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

// CrearComentarioParaClase
// @Summary Crear un comentario para una clase
// @Description Agrega un comentario a una clase por su ID. Se requiere autor (email del usuario), password del usuario, titulo, detalle, meGusta, noMeGusta. La fecha se asigna automáticamente.
// @Tags Comentarios
// @Accept json
// @Produce json
// @Param id path string true "ID de la clase"
// @Param comentario body models.Comentario true "Comentario a crear (autor, password, titulo, detalle, meGusta, noMeGusta)"
// @Success 201 {object} models.Comentario "Comentario creado exitosamente"
// @Failure 400 {object} response.ErrorResponse "Datos inválidos o faltan campos requeridos"
// @Failure 404 {object} response.ErrorResponse "Clase no encontrada o usuario no encontrado/credenciales inválidas"
// @Failure 500 {object} response.ErrorResponse "Error interno del servidor"
// @Router /api/clases/{id}/comentarios [post]
func (c *ComentarioControlador) CrearComentarioParaClase(ctx *gin.Context) {
    claseID := ctx.Param("id")
    var comentario models.Comentario

    if err := ctx.ShouldBindJSON(&comentario); err != nil {
        ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Datos inválidos: " + err.Error()})
        return
    }

    // Validar campos requeridos
    if comentario.Autor == "" || comentario.Password == "" || comentario.Titulo == "" || comentario.Detalle == "" {
        ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Campos autor, password, titulo y detalle son requeridos"})
        return
    }

    // meGusta y noMeGusta pueden ser 0 por defecto

    creado, err := c.servicio.CrearComentarioParaClase(ctx.Request.Context(), claseID, &comentario)
    if err != nil {
        if err.Error() == "clase no encontrada" || err.Error() == "usuario no encontrado o credenciales inválidas" {
            ctx.JSON(http.StatusNotFound, response.ErrorResponse{Message: err.Error()})
        } else {
            ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
        }
        return
    }

    ctx.JSON(http.StatusCreated, creado)
}
