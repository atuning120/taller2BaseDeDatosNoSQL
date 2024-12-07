package controllers

import (
	"net/http"

	"go-API/models"
	"go-API/services"

	"github.com/gin-gonic/gin"
)

// ClaseControlador gestiona las rutas relacionadas con las clases.
type ClaseControlador struct {
	servicio *services.ClaseService
}

// NewClaseControlador crea un nuevo controlador para las clases.
func NewClaseControlador(servicio *services.ClaseService) *ClaseControlador {
	return &ClaseControlador{servicio: servicio}
}

// ObtenerClasesPorUnidad obtiene todas las clases de una unidad.
// @Summary Devuelve las clases de una unidad
// @Description Devuelve todas las clases asociadas a una unidad
// @Tags Clases
// @Accept json
// @Produce json
// @Param id path string true "ID de la unidad"
// @Success 200 {array} response.ClaseResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/unidades/{id}/clases [get]
func (cc *ClaseControlador) ObtenerClasesPorUnidad(c *gin.Context) {
	id := c.Param("id") // ID de la unidad

	clases, err := cc.servicio.ObtenerClasesPorUnidad(id)
	if err != nil {
		if err.Error() == "unidad no encontrada" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, clases)
}

// CrearClaseParaUnidad crea una nueva clase asociada a una unidad.
// @Summary Crear una clase para una unidad
// @Description Agrega una clase a la base de datos asociada a una unidad
// @Tags Clases
// @Param id path string true "ID de la unidad"
// @Param clase body request.CreateClaseRequest true "Clase a crear"
// @Accept json
// @Produce json
// @Success 200 {object} response.CrearClase
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/unidades/{id}/clases [post]
func (cc *ClaseControlador) CrearClaseParaUnidad(c *gin.Context) {
	unidadID := c.Param("id") // ID de la unidad

	var clase models.Clase
	if err := c.ShouldBindJSON(&clase); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	// Llamar al servicio para crear la clase
	result, err := cc.servicio.CrearClaseParaUnidad(unidadID, &clase)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"inserted_id": result.InsertedID})
}
