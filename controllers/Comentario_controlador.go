package controllers

import (
	"go-API/services"

	"github.com/gin-gonic/gin"
)

// UnidadControlador maneja las rutas relacionadas con las unidades.
type ComentarioControlador struct {
	servicio *services.ComentarioService
}

// NewUnidadControlador crea un nuevo controlador para las unidades.
func NewComentarioControlador(servicio *services.ComentarioService) *ComentarioControlador {
	return &ComentarioControlador{servicio: servicio}
}

// ObtenerComentariosPorClase obtiene todos los comentarios de una clase.
// ObtenerComentariosPorClase obtiene todos los comentarios de una clase.
// @Summary Devuelve los comentarios de una clase
// @Description Devuelve todos los comentarios asociados a una clase por su ID
// @Tags Comentarios
// @Accept json
// @Produce json
// @Param id path string true "ID de la clase"
// @Success 200 {array} response.ComentarioResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/clases/{id}/comentarios [get]
func (c *ComentarioControlador) ObtenerComentariosPorClase(ctx *gin.Context) {
	claseID := ctx.Param("id")
	comentarios, err := c.servicio.ObtenerComentariosPorClase(claseID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, comentarios)
}

// CrearComentarioParaClase crea un nuevo comentario para una clase.
// @Summary Crear un comentario para una clase
// @Description Agrega un comentario a una clase por su ID
// @Tags Comentarios
// @Accept json
// @Produce json
// @Param id path string true "ID de la clase"
// @Param comentario body request.CreateComentarioRequest true "Comentario a crear"
// @Success 201 {object} response.ComentarioResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/clases/{id}/comentarios [post]
func (c *ComentarioControlador) CrearComentarioParaClase(ctx *gin.Context) {
	claseID := ctx.Param("id")
	comentario, err := c.servicio.CrearComentarioParaClase(claseID, ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(201, comentario)
}
