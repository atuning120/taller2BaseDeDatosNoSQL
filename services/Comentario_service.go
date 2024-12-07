package services

import (
	"context"
	"errors"
	"go-API/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ClaseService gestiona la lógica relacionada con las clases.
type ComentarioService struct {
	CursoCollection      *mongo.Collection
	UnidadCollection     *mongo.Collection
	ClaseCollection      *mongo.Collection
	ComentarioCollection *mongo.Collection
}

// NewClaseService crea un nuevo servicio para las clases.
func NewComentarioService(db *mongo.Database) *ComentarioService {
	return &ComentarioService{
		CursoCollection:      db.Collection("cursos"),
		ClaseCollection:      db.Collection("clases"), // Asegúrate de asignar la colección de clases aquí
		ComentarioCollection: db.Collection("comentarios"),
	}
}

// ObtenerClasesPorUnidad obtiene todas las clases de una unidad.
func (s *ComentarioService) ObtenerComentariosPorClase(id string) ([]models.Comentario, error) {
	// Convertir el ID a ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("ID inválido")
	}

	// Verificar si la clase existe
	var clase models.Clase
	err = s.ClaseCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&clase)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("clase no encontrada")
		}
		return nil, err
	}

	// Buscar los comentarios asociados a la clase
	cursor, err := s.ComentarioCollection.Find(context.TODO(), bson.M{"clase_id": objectID})
	if err != nil {
		return nil, err
	}

	var comentarios []models.Comentario
	if err = cursor.All(context.TODO(), &comentarios); err != nil {
		return nil, err
	}

	return comentarios, nil
}

// CrearComentarioParaClase crea un nuevo comentario para una clase.
func (s *ComentarioService) CrearComentarioParaClase(id string, ctx *gin.Context) (*models.Comentario, error) {
	// Convertir el ID a ObjectID
	claseID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("ID inválido")
	}

	// Verificar si la clase existe
	var clase models.Clase
	err = s.ClaseCollection.FindOne(context.TODO(), bson.M{"_id": claseID}).Decode(&clase)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("clase no encontrada")
		}
		return nil, err
	}

	// Crear un nuevo comentario
	var comentario models.Comentario
	if err := ctx.BindJSON(&comentario); err != nil {
		return nil, err
	}
	comentario.ClaseID = claseID

	// Insertar el comentario en la base de datos
	result, err := s.ComentarioCollection.InsertOne(context.TODO(), comentario)
	if err != nil {
		return nil, err
	}

	// Asignar el ID del comentario
	comentario.ID = result.InsertedID.(primitive.ObjectID)

	return &comentario, nil
}
