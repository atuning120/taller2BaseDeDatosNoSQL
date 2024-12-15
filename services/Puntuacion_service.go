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

    // Verificar que el curso existe en MongoDB
    var curso models.Curso
    err = s.CursoCollection.FindOne(ctx, bson.M{"_id": objCursoID}).Decode(&curso)
    if err != nil {
        return errors.New("curso no encontrado en MongoDB")
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
            // Si el curso no existe en Neo4j, crearlo
            createCourseQuery := `
                CREATE (c:Course {id: $cursoID, nombre: $nombre, descripcion: $descripcion})
            `
            _, err := tx.Run(ctx, createCourseQuery, map[string]interface{}{
                "cursoID":     cursoID,
                "nombre":      curso.Nombre,
                "descripcion": curso.Descripcion,
            })
            if err != nil {
                return nil, err
            }
        }

        // Verificar que el usuario existe en Neo4j
        checkUserQuery := `
            MATCH (u:User {email: $usuario})
            RETURN COUNT(*) > 0 AS exists
        `
        checkUserResult, err := tx.Run(ctx, checkUserQuery, map[string]interface{}{
            "usuario": email,
        })
        if err != nil {
            return nil, err
        }
        if !checkUserResult.Next(ctx) {
            // Si el usuario no existe en Neo4j, crearlo
            createUserQuery := `
                CREATE (u:User {email: $usuario, nombre: $nombre})
            `
            _, err := tx.Run(ctx, createUserQuery, map[string]interface{}{
                "usuario": email,
                "nombre":  u.Nombre,
            })
            if err != nil {
                return nil, err
            }
        }

        // Verificar si el usuario ya ha puntuado el curso
        puntuacionQuery := `
            MATCH (u:User {email: $usuario})-[r:CALIFICÓ]->(c:Course {id: $cursoID})
            RETURN COUNT(r) > 0 AS hasRated
        `
        puntuacionResult, err := tx.Run(ctx, puntuacionQuery, map[string]interface{}{
            "usuario": email,
            "cursoID": cursoID,
        })
        if err != nil {
            return nil, err
        }
        if puntuacionResult.Next(ctx) {
            hasRated, ok := puntuacionResult.Record().Get("hasRated")
            if ok && hasRated.(bool) {
                return nil, errors.New("el usuario ya ha puntuado este curso")
            }
        }

        // Crear la relación de calificación con el nombre del usuario
        createQuery := `
            MATCH (u:User {email: $usuario}), (c:Course {id: $cursoID})
            CREATE (u)-[r:CALIFICÓ]->(c)
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

// ObtenerPromedioPuntuacionesCurso obtiene el promedio de puntuaciones de un curso desde Neo4j
func (s *PuntuacionService) ObtenerPromedioPuntuacionesCurso(ctx context.Context, cursoID string) (float32, error) {
    session := s.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
    defer session.Close(ctx)

    promedio, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
        return s.calcularPromedioNeo4j(ctx, tx, cursoID)
    })
    if err != nil {
        return 0, err
    }

    return promedio.(float32), nil
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

func parseObjectID(id string) (primitive.ObjectID, error) {
    return primitive.ObjectIDFromHex(id)
}