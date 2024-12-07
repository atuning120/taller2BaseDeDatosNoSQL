package controllers

import (
	"go-API/services"
	"net/http"

	"go-API/models"

	"github.com/gin-gonic/gin"
)

// UnidadControlador maneja las rutas relacionadas con las unidades.
type UsuarioControlador struct {
	servicio *services.UsuarioService
}

// NewUnidadControlador crea un nuevo controlador para las unidades.
func NewUsuarioControlador(servicio *services.UsuarioService) *UsuarioControlador {
	return &UsuarioControlador{servicio: servicio}
}

// ObtenerUsuarios obtiene todos los usuarios.
// @Summary Devuelve todos los usuarios
// @Description Devuelve la lista completa de usuarios registrados
// @Tags Usuarios
// @Accept json
// @Produce json
// @Success 200 {array} models.Usuario
// @Failure 500 {object} response.ErrorResponse
// @Router /api/usuarios [get]
func (uc *UsuarioControlador) ObtenerUsuarios(c *gin.Context) {
	usuarios, err := uc.servicio.ObtenerUsuarios()
	if err != nil {
		c.JSON(
			500,
			gin.H{"error": err.Error()},
		)
		return
	}
	c.JSON(200, usuarios)
}

// CrearUsuario maneja la creación de un nuevo usuario.
// @Summary Crear un nuevo usuario
// @Description Agrega un usuario a la base de datos
// @Tags Usuarios
// @Accept json
// @Produce json
// @Param usuario body request.CreateUsuarioRequest true "Usuario a crear"
// @Success 201 {object} models.Usuario
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/usuarios [post]
func (uc *UsuarioControlador) CrearUsuario(c *gin.Context) {
	var usuario models.Usuario

	// Validar los datos enviados en la solicitud
	if err := c.ShouldBindJSON(&usuario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	// Llamar al servicio para crear el usuario
	result, err := uc.servicio.CrearUsuario(&usuario)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"inserted_id": result.InsertedID})
}

// ObtenerUsuarioPorID obtiene un usuario por su ID.
// @Summary Obtener un usuario por ID
// @Description Devuelve un usuario en específico por su ID
// @Tags Usuarios
// @Accept json
// @Produce json
// @Param id path string true "ID del usuario"
// @Success 200 {object} models.Usuario
// @Failure 500 {object} response.ErrorResponse
// @Router /api/usuarios/{id} [get]
func (uc *UsuarioControlador) ObtenerUsuarioPorID(c *gin.Context) {
	id := c.Param("id")

	usuario, err := uc.servicio.ObtenerUsuarioPorID(id)
	if err != nil {
		c.JSON(
			500,
			gin.H{"error": err.Error()},
		)
		return
	}
	c.JSON(200, usuario)
}

// InscribirseACurso permite que un usuario se inscriba en un curso.
// @Summary Inscribir un usuario en un curso
// @Description Inscribe a un usuario en un curso específico
// @Tags Usuarios
// @Accept json
// @Produce json
// @Param inscripcion body request.InscripcionRequest true "Datos de inscripción"
// @Success 200 {object} response.InscripcionResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/usuarios/inscripcion [post]
func (uc *UsuarioControlador) InscribirseACurso(c *gin.Context) {
	var inscripcion struct {
		UsuarioID string `json:"usuario_id"`
		CursoID   string `json:"curso_id"`
	}

	// Validar los datos enviados en la solicitud
	if err := c.ShouldBindJSON(&inscripcion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	// Llamar al servicio para inscribir al usuario en el curso
	err := uc.servicio.InscribirseACurso(inscripcion.UsuarioID, inscripcion.CursoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inscripción exitosa"})
}

// ObtenerCursosInscritos obtiene los cursos en los que un usuario está inscrito.
// @Summary Obtener cursos inscritos de un usuario
// @Description Devuelve la lista de cursos en los que un usuario está inscrito
// @Tags Usuarios
// @Accept json
// @Produce json
// @Param id path string true "ID del usuario"
// @Success 200 {array} models.Curso
// @Failure 500 {object} response.ErrorResponse
// @Router /api/usuarios/{id}/cursos [get]
func (uc *UsuarioControlador) ObtenerCursosInscritos(c *gin.Context) {
	id := c.Param("id")

	cursos, err := uc.servicio.ObtenerCursosInscritos(id)
	if err != nil {
		c.JSON(
			500,
			gin.H{"error": err.Error()},
		)
		return
	}
	c.JSON(200, cursos)
}
