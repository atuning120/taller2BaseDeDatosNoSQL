package response

import (
	"go-API/models"
	"time"
)

// CursoResponse define la estructura de la respuesta para un curso.
type CursoResponse struct {
	ID          string   `json:"id"`
	Nombre      string   `json:"nombre"`
	Descripcion string   `json:"descripcion"`
	Imagen      string   `json:"imagen_url"`
	Valoracion  float32  `json:"valoracion"`
	Unidades    []string `json:"unidades"` // IDs de las unidades
	Usuarios    int      `json:"cant_usuarios"`
	Comentarios []string `json:"comentarios"` // IDs de los comentarios
}

// NewCursoResponse convierte un modelo Curso en una respuesta CursoResponse.
func NewCursoResponse(curso models.Curso) CursoResponse {
	unidades := make([]string, len(curso.Unidades))
	for i, unidad := range curso.Unidades {
		unidades[i] = unidad.Hex()
	}

	comentarios := make([]string, len(curso.Comentarios))
	for i, comentario := range curso.Comentarios {
		comentarios[i] = comentario.Hex()
	}

	return CursoResponse{
		ID:          curso.ID.Hex(),
		Nombre:      curso.Nombre,
		Descripcion: curso.Descripcion,
		Imagen:      curso.Imagen,
		Valoracion:  curso.Valoracion,
		Unidades:    unidades,
		Usuarios:    curso.Usuarios,
		Comentarios: comentarios,
	}
}

// ErrorResponse define la estructura de las respuestas de error.
type ErrorResponse struct {
	Message string `json:"message"`
}

// CrearCurso define la respuesta cuando se crea un curso exitosamente.
type CrearCurso struct {
	InsertedID string `json:"inserted_id"`
}

// UpdateValoracionResponse define la estructura de la respuesta para la actualización de la valoración.
type UpdateValoracionResponse struct {
	Message               string  `json:"message"`
	ValoracionActualizada float32 `json:"valoracion_actualizada"`
}

// UnidadResponse define la estructura de la respuesta para una unidad.
type UnidadResponse struct {
	ID      string   `json:"id"`
	IDcurso string   `json:"idcurso"`
	Nombre  string   `json:"nombre"`
	Clases  []string `json:"clases"` // IDs de las clases en formato string
}

// NewUnidadResponse convierte un modelo Unidad en una respuesta UnidadResponse.
func NewUnidadResponse(unidad models.Unidad) UnidadResponse {
	clases := make([]string, len(unidad.Clases))
	for i, clase := range unidad.Clases {
		clases[i] = clase.Hex()
	}

	return UnidadResponse{
		ID:      unidad.ID.Hex(),
		IDcurso: unidad.IDcurso.Hex(),
		Nombre:  unidad.Nombre,
		Clases:  clases,
	}
}

// CrearUnidad define la estructura de la respuesta al crear una unidad.
type CrearUnidad struct {
	InsertedID string `json:"inserted_id"`
}

// ClaseResponse define la estructura de la respuesta para una clase.
type ClaseResponse struct {
	ID           string   `json:"id"`
	UnidadID     string   `json:"unidad_id"`
	Nombre       string   `json:"nombre"`
	Descripcion  string   `json:"descripcion"`
	VideoURL     string   `json:"video_url"`
	Adjuntos_url []string `json:"adjuntos_url"`
	MeGusta      int      `json:"me_gusta"`
	NoMeGusta    int      `json:"no_me_gusta"`
	Comentarios  []string `json:"comentarios"`
}

// NewClaseResponse convierte un modelo Clase en una respuesta ClaseResponse.
func NewClaseResponse(clase models.Clase) ClaseResponse {
	comentarios := make([]string, len(clase.Comentarios))
	for i, comentario := range clase.Comentarios {
		comentarios[i] = comentario.Hex()
	}

	return ClaseResponse{
		ID:           clase.ID.Hex(),
		UnidadID:     clase.UnidadID.Hex(),
		Nombre:       clase.Nombre,
		Descripcion:  clase.Descripcion,
		VideoURL:     clase.VideoURL,
		Adjuntos_url: clase.Adjuntos_url, // Correct field name from models.Clase
		MeGusta:      clase.MeGusta,
		NoMeGusta:    clase.NoMeGusta,
		Comentarios:  comentarios,
	}
}

// CrearClase define la estructura de la respuesta al crear una clase.
type CrearClase struct {
	InsertedID string `json:"inserted_id"`
}

// ComentarioResponse define la estructura de la respuesta para un comentario.
type ComentarioResponse struct {
	ID        string    `json:"id"`
	ClaseID   string    `json:"clase_id"`
	Autor     string    `json:"autor"`
	Fecha     time.Time `json:"fecha"`
	Titulo    string    `json:"titulo"`
	Detalle   string    `json:"detalle"`
	MeGusta   int       `json:"me_gusta"`
	NoMeGusta int       `json:"no_me_gusta"`
}

// NewComentarioResponse convierte un modelo Comentario en una respuesta ComentarioResponse.
func NewComentarioResponse(comentario models.Comentario) ComentarioResponse {
	return ComentarioResponse{
		ID:        comentario.ID.Hex(),
		ClaseID:   comentario.ClaseID.Hex(),
		Autor:     comentario.Autor,
		Fecha:     comentario.Fecha,
		Titulo:    comentario.Titulo,
		Detalle:   comentario.Detalle,
		MeGusta:   comentario.MeGusta,
		NoMeGusta: comentario.NoMeGusta,
	}
}

// InscripcionResponse define la estructura de la respuesta al inscribir a un usuario.
type InscripcionResponse struct {
	Message string `json:"message"`
}

// VerClaseResponse define la estructura de la respuesta al ver una clase.
type VerClaseResponse struct {
    Message string `json:"message"`
    Estado  string `json:"estado"`
}
