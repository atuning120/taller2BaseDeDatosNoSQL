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
// @Success 200 {object} map[string]string "message: Puntuación creada exitosamente"
// @Failure 400 {object} map[string]string "error: Bad Request"
// @Failure 404 {object} map[string]string "error: Not Found"
// @Failure 500 {object} map[string]string "error: Internal Server Error"
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
        } else if err.Error() == "curso ya puntuado, intente con otro curso" || err.Error() == "la puntuación debe estar entre 0 y 5" {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
// @Success 200 {object} map[string]float64 "promedio: 0.0"
// @Failure 400 {object} map[string]string "error: Bad Request"
// @Failure 404 {object} map[string]string "error: Not Found"
// @Failure 500 {object} map[string]string "error: Internal Server Error"
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

// ObtenerPuntuacionesPorUsuario obtiene todas las puntuaciones hechas por un usuario.
// @Summary Obtener todas las puntuaciones hechas por un usuario
// @Description Devuelve todas las puntuaciones hechas por un usuario por su email
// @Tags Puntuaciones
// @Accept json
// @Produce json
// @Param email path string true "Email del usuario"
// @Success 200 {array} map[string]interface{} "curso: nombre del curso, valoracion: valor de la puntuación"
// @Failure 400 {object} map[string]string "error: Bad Request"
// @Failure 404 {object} map[string]string "error: Not Found"
// @Failure 500 {object} map[string]string "error: Internal Server Error"
// @Router /api/puntuaciones/usuarios/{email} [get]
func (ctrl *PuntuacionesControlador) ObtenerPuntuacionesPorUsuario(c *gin.Context) {
    email := c.Param("email")

    puntuaciones, err := ctrl.servicio.ObtenerPuntuacionesPorUsuario(email)
    if err != nil {
        if err.Error() == "error al obtener las puntuaciones" {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    c.JSON(http.StatusOK, puntuaciones)
}