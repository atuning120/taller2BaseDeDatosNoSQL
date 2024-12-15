// models/puntuacion.go
package models

// Puntuacion representa una valoración que un usuario da a un curso.
type Puntuacion struct {
	Email  string  `json:"email"`  // email del usuario
	Password string  `json:"password"` // contraseña del usuario
	Valor    float32 `json:"valor"`
}
