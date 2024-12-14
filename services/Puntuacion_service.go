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
// Ahora se usa email y password para construir la clave en Redis: "usuario:email:password".
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

        // Calcular el nuevo promedio y actualizar en MongoDB
        promedio, err := s.ObtenerPromedioPuntuacionesCurso(ctx, cursoID)
        if err != nil {
            return nil, err
        }

        update := bson.M{
            "$set": bson.M{
                "valoracion": promedio,
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

// ObtenerPromedioPuntuacionesCurso obtiene la valoración promedio de un curso.
func (s *PuntuacionService) ObtenerPromedioPuntuacionesCurso(ctx context.Context, cursoID string) (float32, error) {
    session := s.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
    defer session.Close(ctx)

    promedio, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
        query := `
            MATCH (:User)-[r:CALIFICÓ]->(:Course {id: $cursoID})
            RETURN avg(r.valor) AS promedio
        `
        records, err := tx.Run(ctx, query, map[string]interface{}{
            "cursoID": cursoID,
        })
        if err != nil {
            return nil, err
        }

        if records.Next(ctx) {
            record := records.Record()
            promedioValue, ok := record.Get("promedio")
            if !ok || promedioValue == nil {
                return float32(0), nil
            }
            return promedioValue.(float64), nil
        }

        if err = records.Err(); err != nil {
            return nil, err
        }

        return float32(0), nil
    })

    if err != nil {
        return 0, err
    }

    return float32(promedio.(float64)), nil
}

func parseObjectID(id string) (primitive.ObjectID, error) {
    return primitive.ObjectIDFromHex(id)
}
