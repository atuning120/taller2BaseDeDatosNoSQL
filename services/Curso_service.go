package services

import (
    "context"
    "errors"
    "go-API/models"
    "log"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"

    "github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// CursoService gestiona la lógica relacionada con los cursos.
type CursoService struct {
    CursoCollection  *mongo.Collection
    UnidadCollection *mongo.Collection
    ClaseCollection  *mongo.Collection
    Driver           neo4j.DriverWithContext
}

// NewCursoService crea un nuevo servicio para los cursos.
func NewCursoService(db *mongo.Database, driver neo4j.DriverWithContext) *CursoService {
    return &CursoService{
        CursoCollection:  db.Collection("cursos"),
        UnidadCollection: db.Collection("unidades"),
        ClaseCollection:  db.Collection("clases"),
        Driver:           driver,
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

// CrearCurso agrega un nuevo curso a la base de datos y crea el nodo Course en Neo4j.
func (s *CursoService) CrearCurso(curso models.Curso) (*mongo.InsertOneResult, error) {
    // Verificar si las listas son nulas e inicializarlas como vacías
    if curso.Unidades == nil {
        curso.Unidades = []primitive.ObjectID{}
    }
    if curso.Comentarios == nil {
        curso.Comentarios = []primitive.ObjectID{}
    }

    // Insertar el curso en MongoDB
    result, err := s.CursoCollection.InsertOne(context.TODO(), curso)
    if err != nil {
        return nil, err
    }

    // Convertir el ID a ObjectID para obtener el hex
    insertedID := result.InsertedID.(primitive.ObjectID)
    cursoIDHex := insertedID.Hex()

    // Crear el nodo Course en Neo4j
    session := s.Driver.NewSession(context.TODO(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
    defer session.Close(context.TODO())

    _, err = session.ExecuteWrite(context.TODO(), func(tx neo4j.ManagedTransaction) (interface{}, error) {
        createQuery := `
            CREATE (c:Curso {
                id: $id,
                nombre: $nombre
            })
        `

        _, err := tx.Run(context.TODO(), createQuery, map[string]interface{}{
            "id":     cursoIDHex,
            "nombre": curso.Nombre,
        })
        return nil, err
    })

    if err != nil {
        // Si falla crear el nodo en Neo4j, revertimos la inserción en MongoDB
        deleteResult, delErr := s.CursoCollection.DeleteOne(context.TODO(), bson.M{"_id": insertedID})
        if delErr != nil {
            log.Printf("No se pudo eliminar el curso de MongoDB tras fallo en Neo4j: %v", delErr)
        } else {
            log.Printf("Curso eliminado de MongoDB tras fallo en Neo4j, borrados: %d", deleteResult.DeletedCount)
        }

        return nil, err
    }

    return result, nil
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

    // Actualizar la valoración del curso en MongoDB
    _, err = s.CursoCollection.UpdateOne(
        context.TODO(),
        bson.M{"_id": objectID},
        bson.M{"$set": bson.M{"valoracion": valoracion}},
    )
    if err != nil {
        return err
    }

    // Opcional: también podrías actualizar el nodo Course en Neo4j con la nueva valoracion,
    // si deseas mantenerlo sincronizado.
    // Podrías implementar algo similar a CrearCurso en este método.

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
