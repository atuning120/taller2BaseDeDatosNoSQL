package services

import (
	"context"
	"errors"

	"go-API/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CursoService gestiona la lógica relacionada con los cursos.
type CursoService struct {
    CursoCollection  *mongo.Collection
    UnidadCollection *mongo.Collection
    ClaseCollection  *mongo.Collection
}

// NewCursoService crea un nuevo servicio para los cursos.
func NewCursoService(db *mongo.Database) *CursoService {
    return &CursoService{
        CursoCollection:  db.Collection("cursos"),
        UnidadCollection: db.Collection("unidades"),
        ClaseCollection:  db.Collection("clases"),
    }
}

// ObtenerCursos obtiene todos los cursos de la base de datos.
func (s *CursoService) ObtenerCursos() ([]models.Curso, error) {
	var cursos []models.Curso
	cursor, err := s.CursoCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &cursos); err != nil {
		return nil, err
	}
	return cursos, nil
}

// CrearCurso agrega un nuevo curso a la base de datos.
func (s *CursoService) CrearCurso(curso models.Curso) (*mongo.InsertOneResult, error) {
	// Verificar si las listas son nulas e inicializarlas como vacías
	if curso.Unidades == nil {
		curso.Unidades = []primitive.ObjectID{}
	}
	if curso.Comentarios == nil {
		curso.Comentarios = []primitive.ObjectID{}
	}

	return s.CursoCollection.InsertOne(context.TODO(), curso)
}

// ObtenerCursoPorID obtiene un curso por su ID.
func (s *CursoService) ObtenerCursoPorID(id string) (*models.Curso, error) {
	objectID, err := primitive.ObjectIDFromHex(id) // Convertir a ObjectID
	if err != nil {
		return nil, errors.New("ID inválido")
	}

	var curso models.Curso
	err = s.CursoCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&curso)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("curso no encontrado")
		}
		return nil, err
	}

	return &curso, nil
}

// ActualizarValoracion actualiza la valoración promedio de un curso.
func (s *CursoService) ActualizarValoracion(id string, valoracion float32) error {
	// Convertir el ID a ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("ID inválido")
	}

	// Actualizar la valoración del curso
	_, err = s.CursoCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"valoracion": valoracion}},
	)
	if err != nil {
		return err
	}

	return nil
}


// ObtenerClasesPorCurso obtiene todas las clases de un curso.
func (s *CursoService) ObtenerClasesPorCurso(id string) ([]models.Clase, error) {
    objectID, err := primitive.ObjectIDFromHex(id)
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

    // Obtener las unidades del curso
    var unidades []models.Unidad
    cursor, err := s.UnidadCollection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": curso.Unidades}})
    if err != nil {
        return nil, err
    }
    if err = cursor.All(context.TODO(), &unidades); err != nil {
        return nil, err
    }

    // Obtener las clases de cada unidad
    var clases []models.Clase
    for _, unidad := range unidades {
        var unidadClases []models.Clase
        cursor, err := s.ClaseCollection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": unidad.Clases}})
        if err != nil {
            return nil, err
        }
        if err = cursor.All(context.TODO(), &unidadClases); err != nil {
            return nil, err
        }
        clases = append(clases, unidadClases...)
    }

    return clases, nil
}