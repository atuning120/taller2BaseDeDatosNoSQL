package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Clase representa una clase específica dentro de una unidad
type Clase struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	UnidadID     primitive.ObjectID   `bson:"unidad_id" json:"unidad_id"`
	Nombre       string               `bson:"nombre" json:"nombre"`
	Descripcion  string               `bson:"descripcion" json:"descripcion"`
	VideoURL     string               `bson:"video_url" json:"video_url"`
	Adjuntos_url []string             `bson:"adjuntos_url" json:"adjuntos_url"`
	Comentarios  []primitive.ObjectID `bson:"comentarios" json:"comentarios"`
	MeGusta      int                  `bson:"me_gusta" json:"me_gusta"`
	NoMeGusta    int                  `bson:"no_me_gusta" json:"no_me_gusta"`
}
