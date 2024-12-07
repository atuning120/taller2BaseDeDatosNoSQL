package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Usuario representa un usuario que puede inscribirse en cursos
type Usuario struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Nombre    string               `bson:"nombre" json:"nombre"`
	Email     string               `bson:"email" json:"email"`
	Inscritos []primitive.ObjectID `bson:"inscritos" json:"inscritos"` // IDs de cursos inscritos
}
