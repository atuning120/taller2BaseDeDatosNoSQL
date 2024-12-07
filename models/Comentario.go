package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Comentario representa un comentario realizado sobre un curso o clase
type Comentario struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ClaseID   primitive.ObjectID `bson:"clase_id" json:"clase_id"`
	Autor     string             `bson:"autor" json:"autor"`
	Fecha     time.Time          `bson:"fecha" json:"fecha"`
	Titulo    string             `bson:"titulo" json:"titulo"`
	Detalle   string             `bson:"detalle" json:"detalle"`
	MeGusta   int                `bson:"me_gusta" json:"me_gusta"`
	NoMeGusta int                `bson:"no_me_gusta" json:"no_me_gusta"`
}
