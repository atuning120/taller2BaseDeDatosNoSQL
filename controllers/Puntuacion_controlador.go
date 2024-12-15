package controllers

import (
	"go-API/response"
	"go-API/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PuntuacionesControlador struct {
	Servicio *services.PuntuacionService
}

func NewPuntuacionesControlador(servicio *services.PuntuacionService) *PuntuacionesControlador {
	return &PuntuacionesControlador{Servicio: servicio}
}

type PromedioResponse struct {
	Promedio float32 `json:"promedio"`
}

// CrearPuntuacionParaCurso crea una nueva puntuación de un usuario para un curso.
// Ahora requiere usuario (email) y password en el body además del valor.
//
// @Summary Crear una puntuación para un curso
// @Description Agrega una puntuación a un curso por su ID. El usuario se identifica por email y password, se verifica que esté inscrito en el curso. Solo se guarda el nombre del usuario en la relación.
// @Tags Puntuaciones
// @Accept json
// @Produce json
// @Param id path string true "ID del curso (ObjectID en hex)"
// @Param puntuacion body models.Puntuacion true "Puntuación a crear (usuario (email), password, valor)"
// @Success 200 {object} response.MessageResponse "Puntuación creada exitosamente"
// @Failure 400 {object} response.ErrorResponse "Error en la solicitud"
// @Failure 404 {object} response.ErrorResponse "Curso no encontrado"
// @Failure 500 {object} response.ErrorResponse "Error interno del servidor"
// @Router /api/puntuaciones/cursos/{id} [post]
func (ctrl *PuntuacionesControlador) CrearPuntuacionParaCurso(c *gin.Context) {
	var request struct {
		Email    string  `json:"email"`
		Password string  `json:"password"`
		Valor    float32 `json:"valor"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cursoID := c.Param("id")
	err := ctrl.Servicio.CrearPuntuacionParaCurso(c.Request.Context(), cursoID, request.Email, request.Password, request.Valor)
	if err != nil {
		switch err.Error() {
		case "curso no encontrado", "el usuario no está inscrito en este curso", "usuario no encontrado o credenciales inválidas", "el usuario ya ha puntuado este curso":
			c.JSON(http.StatusNotFound, response.ErrorResponse{Message: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{Message: "Puntuación creada exitosamente"})
}

// ObtenerPromedioPuntuacion obtiene el promedio de puntuaciones de un curso.
//
// @Summary Obtener el promedio de puntuaciones de un curso
// @Description Obtiene el promedio de puntuaciones de un curso por su ID
// @Tags Puntuaciones
// @Accept json
// @Produce json
// @Param id path string true "ID del curso (ObjectID en hex)"
// @Success 200 {object} PromedioResponse "Devuelve el promedio en un campo 'promedio'"
// @Failure 500 {object} response.ErrorResponse "Error interno del servidor"
// @Router /api/puntuaciones/cursos/{id}/promedio [get]
func (p *PuntuacionesControlador) ObtenerPromedioPuntuacion(ctx *gin.Context) {
	cursoID := ctx.Param("id")

	promedio, err := p.Servicio.ObtenerPromedioPuntuacionesCurso(ctx.Request.Context(), cursoID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, PromedioResponse{Promedio: promedio})
}
