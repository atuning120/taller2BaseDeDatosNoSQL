package controllers

import (
    "net/http"

    "go-API/services"

    "github.com/gin-gonic/gin"
)

type PuntuacionesControlador struct {
    servicio *services.PuntuacionService
}

func NewPuntuacionesControlador(servicio *services.PuntuacionService) *PuntuacionesControlador {
    return &PuntuacionesControlador{servicio: servicio}
}

// CrearPuntuacionParaCurso crea una puntuación para un curso.
// @Summary Crear una puntuación para un curso
// @Description Agrega una puntuación a un curso por su ID. El usuario se identifica por email y password, se verifica que esté inscrito en el curso. Solo se guarda el nombre del usuario en la relación.
// @Tags Puntuaciones
// @Accept json
// @Produce json
// @Param id path string true "ID del curso"
// @Param puntuacion body models.Puntuacion true "Puntuación a crear (usuario (email), password, valor)"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/puntuaciones/cursos/{id} [post]
func (ctrl *PuntuacionesControlador) CrearPuntuacionParaCurso(c *gin.Context) {
    id := c.Param("id")

    var request struct {
        Email    string  `json:"email"`
        Password string  `json:"password"`
        Valor    float32 `json:"valor"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err := ctrl.servicio.CrearPuntuacionParaCurso(request.Email, request.Password, id, request.Valor)
    if err != nil {
        if err.Error() == "usuario no encontrado" || err.Error() == "el usuario no está inscrito en este curso" {
            c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Puntuación creada exitosamente"})
}

// ObtenerPromedioPuntuacion obtiene el promedio de puntuaciones de un curso.
// @Summary Obtener el promedio de puntuaciones de un curso
// @Description Devuelve el promedio de puntuaciones de un curso por su ID
// @Tags Puntuaciones
// @Accept json
// @Produce json
// @Param id path string true "ID del curso"
// @Success 200 {object} gin.H{"promedio": float64}
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/puntuaciones/cursos/{id}/promedio [get]
func (ctrl *PuntuacionesControlador) ObtenerPromedioPuntuacion(c *gin.Context) {
    id := c.Param("id")

    promedio, err := ctrl.servicio.ObtenerPromedioPuntuacion(id)
    if err != nil {
        if err.Error() == "no se encontraron puntuaciones" {
            c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{"promedio": promedio})
}


