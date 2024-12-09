package controllers

import (
    "go-API/services"
    "go-API/models"
    "go-API/request"
    "net/http"

    "github.com/gin-gonic/gin"
)

// UsuarioControlador maneja las rutas relacionadas con los usuarios.
type UsuarioControlador struct {
    servicio *services.UsuarioService
}

// NewUsuarioControlador crea un nuevo controlador para los usuarios.
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
    var input request.CreateUsuarioRequest

    // Validar los datos enviados en la solicitud
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
        return
    }

    // Crear el objeto Usuario
    usuario := models.Usuario{
        Nombre:   input.Nombre,
        Email:    input.Email,
        Password: input.Password,
    }

    // Llamar al servicio para crear el usuario
    result, err := uc.servicio.CrearUsuario(&usuario)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"inserted_id": result})
}

// ObtenerUsuarioPorCorreoYContrasena obtiene un usuario por su correo y contraseña.
// @Summary Obtener un usuario por correo y contraseña
// @Description Devuelve un usuario en específico por su correo y contraseña
// @Tags Usuarios
// @Accept json
// @Produce json
// @Param email query string true "Correo del usuario"
// @Param password query string true "Contraseña del usuario"
// @Success 200 {object} models.Usuario
// @Failure 500 {object} response.ErrorResponse
// @Router /api/usuarios/usuario [get]
func (uc *UsuarioControlador) ObtenerUsuarioPorCorreoYContrasena(c *gin.Context) {
    email := c.Query("email")
    password := c.Query("password")

    usuario, err := uc.servicio.ObtenerUsuarioPorCorreoYContrasena(email, password)
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
    var inscripcion request.InscripcionRequest

    // Validar los datos enviados en la solicitud
    if err := c.ShouldBindJSON(&inscripcion); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
        return
    }

    // Llamar al servicio para inscribir al usuario en el curso
    err := uc.servicio.InscribirseACurso(inscripcion.Email, inscripcion.Password, inscripcion.CursoID)
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
// @Param email query string true "Correo del usuario"
// @Param password query string true "Contraseña del usuario"
// @Success 200 {array} models.Curso
// @Failure 500 {object} response.ErrorResponse
// @Router /api/usuarios/cursos [get]
func (uc *UsuarioControlador) ObtenerCursosInscritos(c *gin.Context) {
    email := c.Query("email")
    password := c.Query("password")

    if email == "" || password == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Email y password son requeridos"})
        return
    }

    cursos, err := uc.servicio.ObtenerCursosInscritos(email, password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, cursos)
}


// VerClase permite que un usuario vea una clase y actualiza su progreso en el curso.
// @Summary Ver una clase
// @Description Permite que un usuario vea una clase y actualiza su progreso en el curso
// @Tags Usuarios
// @Accept json
// @Produce json
// @Param email path string true "Correo del usuario"
// @Param password path string true "Contraseña del usuario"
// @Param clase_id path string true "ID de la clase"
// @Success 200 {object} response.VerClaseResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/usuarios/{email}/{password}/clases/{clase_id} [post]
func (uc *UsuarioControlador) VerClase(c *gin.Context) {
    email := c.Param("email")
    password := c.Param("password")
    claseID := c.Param("clase_id")

    err := uc.servicio.VerClase(email, password, claseID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Clase vista exitosamente"})
}

// ObtenerProgresoCursos obtiene el progreso de los cursos en los que un usuario está inscrito.
// @Summary Devuelve el progreso de los cursos de un usuario
// @Description Devuelve el progreso de los cursos en los que un usuario está inscrito
// @Tags Usuarios
// @Accept json
// @Produce json
// @Param email query string true "Email del usuario"
// @Param password query string true "Contraseña del usuario"
// @Success 200 {array} models.ProgresoCurso
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/usuarios/progreso [get]
func (uc *UsuarioControlador) ObtenerProgresoCursos(c *gin.Context) {
    email := c.Query("email")
    password := c.Query("password")

    progresos, err := uc.servicio.ObtenerProgresoCursos(email, password)
    if err != nil {
        if err.Error() == "usuario no encontrado" {
            c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    c.JSON(http.StatusOK, progresos)
}