package services

import (
	"context"
	"errors"

	"go-API/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ClaseService gestiona la lógica relacionada con las clases.
type ClaseService struct {
	CursoCollection  *mongo.Collection
	UnidadCollection *mongo.Collection
	ClaseCollection  *mongo.Collection
}

// NewClaseService crea un nuevo servicio para las clases.
func NewClaseService(db *mongo.Database) *ClaseService {
	return &ClaseService{
		CursoCollection:  db.Collection("cursos"),
		UnidadCollection: db.Collection("unidades"),
		ClaseCollection:  db.Collection("clases"), // Asegúrate de asignar la colección de clases aquí
	}
}

// ObtenerClasesPorUnidad obtiene todas las clases de una unidad.
func (s *ClaseService) ObtenerClasesPorUnidad(id string) ([]models.Clase, error) {
	// Convertir el ID a ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("ID inválido")
	}

	// Verificar si la unidad existe
	var unidad models.Unidad
	err = s.UnidadCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&unidad)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("unidad no encontrada")
		}
		return nil, err
	}

	// Buscar las clases asociadas a la unidad
	cursor, err := s.ClaseCollection.Find(context.TODO(), bson.M{"unidad_id": objectID})
	if err != nil {
		return nil, err
	}

	var clases []models.Clase
	if err = cursor.All(context.TODO(), &clases); err != nil {
		return nil, err
	}

	return clases, nil
}

// CrearClaseParaUnidad crea una nueva clase y la asocia a una unidad.
func (s *ClaseService) CrearClaseParaUnidad(unidadID string, clase *models.Clase) (*mongo.InsertOneResult, error) {
	// Convertir el ID de la unidad a ObjectID
	objectID, err := primitive.ObjectIDFromHex(unidadID)
	if err != nil {
		return nil, errors.New("ID de unidad inválido")
	}

	// Verificar si la unidad existe
	var unidad models.Unidad
	err = s.UnidadCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&unidad)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("unidad no encontrada")
		}
		return nil, err
	}

	// Asegurarse de que las listas estén inicializadas
	if clase.Adjuntos_url == nil {
		clase.Adjuntos_url = []string{}
	}
	if clase.Comentarios == nil {
		clase.Comentarios = []primitive.ObjectID{}
	}

	// Asignar el ID de la unidad a la clase
	clase.UnidadID = objectID
	clase.ID = primitive.NewObjectID() // Generar un nuevo ObjectID para la clase

	// Insertar la clase en la colección de clases
	result, err := s.ClaseCollection.InsertOne(context.TODO(), clase)
	if err != nil {
		return nil, err
	}

	// Actualizar la unidad con el ID de la nueva clase
	_, err = s.UnidadCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectID},
		bson.M{"$push": bson.M{"clases": clase.ID}},
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
