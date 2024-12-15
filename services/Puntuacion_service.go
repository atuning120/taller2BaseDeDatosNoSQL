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

// CrearPuntuacionParaCurso crea una puntuación para un curso y actualiza la valoración promedio.
func (s *PuntuacionService) CrearPuntuacionParaCurso(email, password, cursoID string, valor float32) error {
    // Verificar si el usuario está inscrito en el curso
    key := "usuario:" + email + ":" + password
    val, err := s.RedisClient.Get(context.TODO(), key).Result()
    if err == redis.Nil {
        return errors.New("usuario no encontrado")
    } else if err != nil {
        return err
    }

    var usuario models.Usuario
    if err := json.Unmarshal([]byte(val), &usuario); err != nil {
        return err
    }

    cursoObjectID, err := primitive.ObjectIDFromHex(cursoID)
    if err != nil {
        return errors.New("ID de curso inválido")
    }

    inscrito := false
    for _, inscritoID := range usuario.Inscritos {
        if inscritoID == cursoObjectID {
            inscrito = true
            break
        }
    }

    if !inscrito {
        return errors.New("el usuario no está inscrito en este curso")
    }

    // Crear la relación de puntuación en Neo4j
    session := s.Driver.NewSession(context.TODO(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
    defer session.Close(context.TODO())

    _, err = session.ExecuteWrite(context.TODO(), func(tx neo4j.ManagedTransaction) (interface{}, error) {
        query := `
            MATCH (u:Usuario {email: $email}), (c:Curso {id: $cursoID})
            MERGE (u)-[r:PUNTUO]->(c)
            ON CREATE SET r.valor = $valor
            RETURN r
        `
        params := map[string]interface{}{
            "email":   email,
            "cursoID": cursoID,
            "valor":   valor,
        }
        _, err := tx.Run(context.TODO(), query, params)
        return nil, err
    })

    if err != nil {
        return err
    }

    // Calcular el promedio de puntuaciones y actualizar la valoración del curso
    return s.ActualizarValoracionCurso(cursoID)
}

// ActualizarValoracionCurso actualiza la valoración promedio de un curso en MongoDB y Neo4j.
func (s *PuntuacionService) ActualizarValoracionCurso(cursoID string) error {
    session := s.Driver.NewSession(context.TODO(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
    defer session.Close(context.TODO())

    result, err := session.ExecuteRead(context.TODO(), func(tx neo4j.ManagedTransaction) (interface{}, error) {
        query := `
            MATCH (c:Curso {id: $cursoID})<-[r:PUNTUO]-()
            RETURN AVG(r.valor) AS promedio
        `
        params := map[string]interface{}{
            "cursoID": cursoID,
        }
        res, err := tx.Run(context.TODO(), query, params)
        if err != nil {
            return nil, err
        }
        if res.Next(context.TODO()) {
            promedio, _ := res.Record().Get("promedio")
            return promedio, nil
        }
        return nil, errors.New("no se encontraron puntuaciones")
    })

    if err != nil {
        return err
    }

    promedio, ok := result.(float64)
    if !ok {
        return errors.New("error al calcular el promedio")
    }

    // Actualizar la valoración del curso en MongoDB
    objectID, err := primitive.ObjectIDFromHex(cursoID)
    if err != nil {
        return errors.New("ID de curso inválido")
    }

    _, err = s.CursoCollection.UpdateOne(
        context.TODO(),
        bson.M{"_id": objectID},
        bson.M{"$set": bson.M{"valoracion": promedio}},
    )
    if err != nil {
        return err
    }

    return nil
}

// ObtenerPromedioPuntuacion obtiene el promedio de puntuaciones de un curso.
func (s *PuntuacionService) ObtenerPromedioPuntuacion(cursoID string) (float64, error) {
    session := s.Driver.NewSession(context.TODO(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
    defer session.Close(context.TODO())

    result, err := session.ExecuteRead(context.TODO(), func(tx neo4j.ManagedTransaction) (interface{}, error) {
        query := `
            MATCH (c:Curso {id: $cursoID})<-[r:PUNTUO]-()
            RETURN AVG(r.valor) AS promedio
        `
        params := map[string]interface{}{
            "cursoID": cursoID,
        }
        res, err := tx.Run(context.TODO(), query, params)
        if err != nil {
            return nil, err
        }
        if res.Next(context.TODO()) {
            promedio, _ := res.Record().Get("promedio")
            return promedio, nil
        }
        return nil, errors.New("no se encontraron puntuaciones")
    })

    if err != nil {
        return 0, err
    }

    promedio, ok := result.(float64)
    if !ok {
        return 0, errors.New("error al calcular el promedio")
    }

    return promedio, nil
}