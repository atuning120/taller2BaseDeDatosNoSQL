package controllers

import (
	"net/http"

	"go-API/models"
	"go-API/services"

	"github.com/gin-gonic/gin"
)

type CursoControlador struct {
	servicio *services.CursoService
}

func NewCursoControlador(servicio *services.CursoService) *CursoControlador {
	return &CursoControlador{servicio: servicio}
}

// ObtenerCursos devuelve todos los cursos disponibles.
// @Summary Devuelve todos los cursos
// @Description Devuelve todos los cursos disponibles
// @Tags Cursos
// @Accept json
// @Produce json
// @Success 200 {array} response.CursoResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/cursos [get]
func (ctrl *CursoControlador) ObtenerCursos(c *gin.Context) {
	cursos, err := ctrl.servicio.ObtenerCursos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cursos)
}

// CrearCurso crea un nuevo curso.
// @Summary Crear un curso
// @Description Agrega un curso a la base de datos
// @Tags Cursos
// @Param curso body request.CreateCursoRequest true "Curso a crear"
// @Accept json
// @Produce json
// @Success 200 {object} response.CrearCurso
// @Failure 500 {object} response.ErrorResponse
// @Router /api/cursos [post]
func (ctrl *CursoControlador) CrearCurso(c *gin.Context) {
	var request struct {
		Nombre      string  `json:"nombre"`
		Descripcion string  `json:"descripcion"`
		Imagen      string  `json:"imagen_url"`
		Valoracion  float32 `json:"valoracion"`
	}

	// Verificar si los datos enviados son correctos
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Crear un nuevo curso usando el constructor
	curso := models.NewCurso(request.Nombre, request.Descripcion, request.Imagen, request.Valoracion)

	result, err := ctrl.servicio.CrearCurso(curso)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"inserted_id": result.InsertedID})
}

// ObtenerCursoPorID devuelve un curso específico por su ID.
// @Summary Devuelve un curso según su ID
// @Description Devuelve un curso en específico dado su ID
// @Tags Cursos
// @Accept json
// @Produce json
// @Param id path string true "ID del curso"
// @Success 200 {object} response.CursoResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/cursos/{id} [get]
func (ctrl *CursoControlador) ObtenerCursoPorID(c *gin.Context) {
	id := c.Param("id")
	curso, err := ctrl.servicio.ObtenerCursoPorID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, curso)
}

// ActualizarValoracion actualiza la valoración promedio de un curso.
// @Summary Actualiza la valoración de un curso
// @Description Actualiza la valoración de un curso según la nueva valoración proporcionada
// @Tags Cursos
// @Accept json
// @Produce json
// @Param id path string true "ID del curso"
// @Param valoracion body request.UpdateValoracionRequest true "Nueva valoración del curso"
// @Success 200 {object} response.UpdateValoracionResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/cursos/{id}/valoracion [patch]
func (cc *CursoControlador) ActualizarValoracion(c *gin.Context) {
	id := c.Param("id") // ID del curso

	// Estructura para recibir la nueva valoración
	var body struct {
		Valoracion float32 `json:"valoracion"`
	}

	// Validar el cuerpo de la solicitud
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	// Llamar al servicio para actualizar la valoración
	err := cc.servicio.ActualizarValoracion(id, body.Valoracion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Valoración actualizada exitosamente"})
}
