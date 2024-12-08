package services

import (
	"context"
	"errors"

	"go-API/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UnidadService maneja la lógica relacionada con las unidades.
type UnidadService struct {
	UnidadCollection *mongo.Collection
	CursoCollection  *mongo.Collection
}

// NewUnidadService crea un nuevo servicio para las unidades.
func NewUnidadService(db *mongo.Database) *UnidadService {
	return &UnidadService{
		UnidadCollection: db.Collection("unidades"),
		CursoCollection:  db.Collection("cursos"),
	}
}

// ObtenerUnidadesPorCurso obtiene todas las unidades asociadas a un curso.
func (s *UnidadService) ObtenerUnidadesPorCurso(id string) ([]models.Unidad, error) {
	objectID, err := primitive.ObjectIDFromHex(id) // Convertir a ObjectID
	if err != nil {
		return nil, errors.New("ID inválido")
	}

	// Obtener el curso por su ID
	var curso models.Curso
	err = s.CursoCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&curso)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("curso no encontrado")
		}
		return nil, err
	}

	// Buscar las unidades cuyo ID esté en la lista del curso
	cursor, err := s.UnidadCollection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": curso.Unidades}})
	if err != nil {
		return nil, err
	}

	var unidades []models.Unidad
	if err = cursor.All(context.TODO(), &unidades); err != nil {
		return nil, err
	}

	return unidades, nil
}

// CrearUnidad crea una nueva unidad y la asocia a un curso.
func (s *UnidadService) CrearUnidad(id string, unidad models.Unidad) (*mongo.InsertOneResult, error) {
    objectID, err := primitive.ObjectIDFromHex(id) // Convertir a ObjectID
    if err != nil {
        return nil, errors.New("ID inválido")
    }

    // Verificar si el curso existe
    var curso models.Curso
    err = s.CursoCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&curso)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, errors.New("curso no encontrado")
        }
        return nil, err
    }

    // Crear la nueva unidad con el ID del curso
    nuevaUnidad := models.Unidad{
        ID:      primitive.NewObjectID(),
        IDcurso: objectID,
        Nombre:  unidad.Nombre,
        Clases:  []primitive.ObjectID{},
    }

    // Insertar la unidad en la colección de unidades
    result, err := s.UnidadCollection.InsertOne(context.TODO(), nuevaUnidad)
    if err != nil {
        return nil, err
    }

    // Agregar el ID de la nueva unidad al curso
    _, err = s.CursoCollection.UpdateOne(
		context.TODO(),
        bson.M{"_id": objectID},
        bson.M{"$push": bson.M{"unidades": result.InsertedID}},
    )
    if err != nil {
        return nil, err
    }

    return result, nil
}
