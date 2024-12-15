package controllers

import (
	"net/http"

	"go-API/services"

	"github.com/gin-gonic/gin"
)

type ComentarioCursoControlador struct {
	servicio *services.ComentarioCursoService
}

func NewComentarioCursoControlador(servicio *services.ComentarioCursoService) *ComentarioCursoControlador {
	return &ComentarioCursoControlador{servicio: servicio}
}

// CrearComentarioCurso crea un comentario para un curso.
// @Summary Crear un comentario para un curso
// @Description Agrega un comentario a un curso por su ID. El usuario se identifica por email.
// @Tags ComentariosCurso
// @Accept json
// @Produce json
// @Param comentario body models.ComentarioCurso true "Comentario a crear (usuario (email), cursoID, texto)"
// @Success 200 {object} map[string]string "message: Comentario creado exitosamente"
// @Failure 400 {object} map[string]string "error: Bad Request"
// @Failure 500 {object} map[string]string "error: Internal Server Error"
// @Router /api/comentarios_curso [post]
func (ctrl *ComentarioCursoControlador) CrearComentarioCurso(c *gin.Context) {
	var request struct {
		Email   string `json:"email"`
		CursoID string `json:"curso_id"`
		Texto   string `json:"texto"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.servicio.CrearComentarioCurso(request.Email, request.CursoID, request.Texto)
	if err != nil {
		if err.Error() == "el comentario debe tener al menos 15 caracteres" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comentario creado exitosamente"})
}

// ObtenerComentariosCursoPorUsuario obtiene todos los comentarios hechos por un usuario.
// @Summary Obtener todos los comentarios hechos por un usuario
// @Description Devuelve todos los comentarios hechos por un usuario por su email
// @Tags ComentariosCurso
// @Accept json
// @Produce json
// @Param email path string true "Email del usuario"
// @Success 200 {array} map[string]interface{} "curso: nombre del curso, comentario: texto del comentario"
// @Failure 400 {object} map[string]string "error: Bad Request"
// @Failure 500 {object} map[string]string "error: Internal Server Error"
// @Router /api/comentarios_curso/usuarios/{email} [get]
func (ctrl *ComentarioCursoControlador) ObtenerComentariosCursoPorUsuario(c *gin.Context) {
	email := c.Param("email")

	comentarios, err := ctrl.servicio.ObtenerComentariosCursoPorUsuario(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comentarios)
}
