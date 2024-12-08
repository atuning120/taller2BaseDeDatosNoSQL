package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)

// ProgresoCurso representa el progreso de un usuario en un curso
type ProgresoCurso struct {
    CursoID   primitive.ObjectID `bson:"curso_id" json:"curso_id"`
    ClasesVistas []primitive.ObjectID `bson:"clases_vistas" json:"clases_vistas"`
    Estado    string             `bson:"estado" json:"estado"` // INICIADO, EN CURSO, COMPLETADO
}

// Usuario representa un usuario que puede inscribirse en cursos
type Usuario struct {
    Nombre           string               `bson:"nombre" json:"nombre"`
    Password         string               `bson:"password" json:"password"`
    Email            string               `bson:"email" json:"email"`
    Inscritos        []primitive.ObjectID `bson:"inscritos" json:"inscritos"` // IDs de cursos inscritos
    FechaInscripcion []time.Time          `bson:"fecha_inscripcion" json:"fecha_inscripcion"`
    Progresos        []ProgresoCurso      `bson:"progresos" json:"progresos"` // Progreso de los cursos
}

// NewUsuario crea una nueva instancia de Usuario con listas vacías
func NewUsuario(nombre, password, email string) *Usuario {
    return &Usuario{
        Nombre:           nombre,
        Password:         password,
        Email:            email,
        Inscritos:        []primitive.ObjectID{},
        FechaInscripcion: []time.Time{},
        Progresos:        []ProgresoCurso{},
    }
}