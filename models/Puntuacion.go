// models/puntuacion.go
package models

import "time"

// Puntuacion representa una valoración que un usuario da a un curso.
type Puntuacion struct {
    Usuario  string    `json:"usuario"`  // email del usuario
    Password string    `json:"password"` // contraseña del usuario
    Valor    float32   `json:"valor"`
    Fecha    time.Time `json:"fecha"`
}
