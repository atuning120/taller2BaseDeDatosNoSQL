// models/comentario.go
package models

import (
    "time"
)

// Comentario representa un comentario realizado por un usuario en una clase.
type Comentario struct {
    ID        string    `json:"id"`
    ClaseID   string    `json:"clase_id"`
    Autor     string    `json:"autor"`
    Fecha     time.Time `json:"fecha"`
    Titulo    string    `json:"titulo"`
    Detalle   string    `json:"detalle"`
    MeGusta   int       `json:"me_gusta"`
    NoMeGusta int       `json:"no_me_gusta"`
}
