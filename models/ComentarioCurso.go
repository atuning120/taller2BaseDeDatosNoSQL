package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// ComentarioCurso representa un comentario que un usuario hace a un curso.
type ComentarioCurso struct {
    Email   string             `bson:"email" json:"email"`
    CursoID primitive.ObjectID `bson:"curso_id" json:"curso_id"`
    Texto   string             `bson:"texto" json:"texto"`
}