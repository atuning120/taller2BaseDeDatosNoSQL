package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)

// Usuario representa un usuario que puede inscribirse en cursos
type Usuario struct {
    Nombre           string               `bson:"nombre" json:"nombre"`
    Password         string               `bson:"password" json:"password"`
    Email            string               `bson:"email" json:"email"`
    Inscritos        []primitive.ObjectID `bson:"inscritos" json:"inscritos"` // IDs de cursos inscritos
    FechaInscripcion []time.Time          `bson:"fecha_inscripcion" json:"fecha_inscripcion"`
}

// NewUsuario crea una nueva instancia de Usuario con listas vacías
func NewUsuario(nombre, password, email string) *Usuario {
    return &Usuario{
        Nombre:           nombre,
        Password:         password,
        Email:            email,
        Inscritos:        []primitive.ObjectID{},
        FechaInscripcion: []time.Time{},
    }
}