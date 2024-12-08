package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Curso representa un curso en la base de datos.
type Curso struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Nombre      string               `bson:"nombre" json:"nombre"`
	Descripcion string               `bson:"descripcion" json:"descripcion"`
	Imagen      string               `bson:"imagen_url" json:"imagen_url"`
	Valoracion  float32              `bson:"valoracion" json:"valoracion"`
	Unidades    []primitive.ObjectID `bson:"unidades" json:"unidades"` // Lista de IDs de unidades
	Usuarios    int                  `bson:"cant_usuarios" json:"cant_usuarios"`
	Comentarios []primitive.ObjectID `bson:"comentarios" json:"comentarios"` // Lista de IDs de comentarios
	Clases int `bson:"cant_clases" json:"cant_clases"`
}

// NewCurso crea un nuevo curso con listas vacías por defecto.
func NewCurso(nombre, descripcion, imagen string, valoracion float32) Curso {
	return Curso{
		ID:          primitive.NewObjectID(),
		Nombre:      nombre,
		Descripcion: descripcion,
		Imagen:      imagen,
		Valoracion:  valoracion,
		Unidades:    []primitive.ObjectID{}, // Inicializado como lista vacía
		Comentarios: []primitive.ObjectID{}, // Inicializado como lista vacía
	}
}
