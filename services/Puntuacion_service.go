// services/puntuacion_service.go
package services

import (
    "context"
    "encoding/json"
    "errors"
    "go-API/models"

    "github.com/go-redis/redis/v8"
    "github.com/neo4j/neo4j-go-driver/v5/neo4j"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

type PuntuacionService struct {
    Driver          neo4j.DriverWithContext
    CursoCollection *mongo.Collection
    RedisClient     *redis.Client
}

func NewPuntuacionService(driver neo4j.DriverWithContext, cursoCollection *mongo.Collection, redisClient *redis.Client) *PuntuacionService {
    return &PuntuacionService{
        Driver:          driver,
        CursoCollection: cursoCollection,
        RedisClient:     redisClient,
    }
}

// CrearPuntuacionParaCurso crea una nueva puntuación de un usuario para un curso.
// Calcula el nuevo promedio desde Neo4j y actualiza la valoración en MongoDB.
// No actualiza la valoración en el nodo Course de Neo4j.
func (s *PuntuacionService) CrearPuntuacionParaCurso(ctx context.Context, cursoID string, email string, password string, valor float32) error {
    if valor < 0 || valor > 5 {
        return errors.New("la puntuación debe estar entre 0 y 5")
    }

    // Obtener el usuario desde Redis usando email y password
    key := "usuario:" + email + ":" + password
    val, err := s.RedisClient.Get(ctx, key).Result()
    if err == redis.Nil {
        return errors.New("usuario no encontrado o credenciales inválidas")
    } else if err != nil {
        return err
    }

    var u models.Usuario
    if err := json.Unmarshal([]byte(val), &u); err != nil {
        return err
    }

    // Verificar si el usuario está inscrito en el curso
    objCursoID, err := parseObjectID(cursoID)
    if err != nil {
        return errors.New("cursoID inválido")
    }

    inscrito := false
    for _, cID := range u.Inscritos {
        if cID == objCursoID {
            inscrito = true
            break
        }
    }

    if !inscrito {
        return errors.New("el usuario no está inscrito en este curso")
    }

    session := s.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
    defer session.Close(ctx)

    _, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
        // Verificar que el curso existe en Neo4j
        checkQuery := `
            MATCH (c:Course {id: $cursoID})
            RETURN COUNT(*) > 0 AS exists
        `
        checkResult, err := tx.Run(ctx, checkQuery, map[string]interface{}{
            "cursoID": cursoID,
        })
        if err != nil {
            return nil, err
        }
        if !checkResult.Next(ctx) {
            return nil, errors.New("curso no encontrado")
        }
        exists, ok := checkResult.Record().Get("exists")
        if !ok || !exists.(bool) {
            return nil, errors.New("curso no encontrado")
        }

        // Crear o actualizar la relación de calificación con el nombre del usuario
        createQuery := `
            MATCH (u:User {email: $usuario}), (c:Course {id: $cursoID})
            MERGE (u)-[r:CALIFICÓ]->(c)
            SET r.valor = $valor, r.fecha = datetime(), r.nombre = $nombre
            RETURN r.valor AS valor
        `
        result, err := tx.Run(ctx, createQuery, map[string]interface{}{
            "usuario": email,
            "cursoID": cursoID,
            "valor":   valor,
            "nombre":  u.Nombre,
        })
        if err != nil {
            return nil, err
        }
        _, err = result.Consume(ctx)
        if err != nil {
            return nil, err
        }

        // Obtener el nuevo promedio desde Neo4j
        nuevoPromedio, err := s.calcularPromedioNeo4j(ctx, tx, cursoID)
        if err != nil {
            return nil, err
        }

        // Actualizar la valoración en MongoDB
        update := bson.M{
            "$set": bson.M{
                "valoracion": nuevoPromedio,
            },
        }

        _, err = s.CursoCollection.UpdateOne(ctx, bson.M{"_id": objCursoID}, update)
        if err != nil {
            return nil, err
        }

        return nil, nil
    })

    return err
}

// calcularPromedioNeo4j obtiene el promedio directamente desde Neo4j
func (s *PuntuacionService) calcularPromedioNeo4j(ctx context.Context, tx neo4j.ManagedTransaction, cursoID string) (float32, error) {
    query := `
        MATCH (:User)-[r:CALIFICÓ]->(:Course {id: $cursoID})
        RETURN avg(r.valor) AS promedio
    `
    records, err := tx.Run(ctx, query, map[string]interface{}{
        "cursoID": cursoID,
    })
    if err != nil {
        return 0, err
    }

    if records.Next(ctx) {
        record := records.Record()
        promedioValue, ok := record.Get("promedio")
        if !ok || promedioValue == nil {
            return 0, nil
        }

        switch v := promedioValue.(type) {
        case float64:
            return float32(v), nil
        case float32:
            return v, nil
        default:
            return 0, nil
        }
    }

    if err = records.Err(); err != nil {
        return 0, err
    }

    return 0, nil
}

// ObtenerPromedioPuntuacionesCurso obtiene la valoración del curso desde MongoDB (ya no se calcula en Neo4j)
func (s *PuntuacionService) ObtenerPromedioPuntuacionesCurso(ctx context.Context, cursoID string) (float32, error) {
    objCursoID, err := parseObjectID(cursoID)
    if err != nil {
        return 0, errors.New("cursoID inválido")
    }

    // Obtener el curso desde MongoDB y retornar el campo valoracion
    var curso models.Curso
    err = s.CursoCollection.FindOne(ctx, bson.M{"_id": objCursoID}).Decode(&curso)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return 0, errors.New("curso no encontrado")
        }
        return 0, err
    }

    // Retornar la valoracion del curso directamente
    return curso.Valoracion, nil
}

func parseObjectID(id string) (primitive.ObjectID, error) {
    return primitive.ObjectIDFromHex(id)
}
