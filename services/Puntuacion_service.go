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
	// Verificar si la puntuación está en el rango de 0 a 5
	if valor < 0 || valor > 5 {
		return errors.New("la puntuación debe estar entre 0 y 5")
	}

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

	// Verificar si el usuario ya ha puntuado el curso
	session := s.Driver.NewSession(context.TODO(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(context.TODO())

	exists, err := session.ExecuteRead(context.TODO(), func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
            MATCH (u:Usuario {email: $email})-[r:PUNTUO]->(c:Curso {id: $cursoID})
            RETURN COUNT(r) > 0 AS exists
        `
		params := map[string]interface{}{
			"email":   email,
			"cursoID": cursoID,
		}
		res, err := tx.Run(context.TODO(), query, params)
		if err != nil {
			return false, err
		}
		if res.Next(context.TODO()) {
			return res.Record().Values[0].(bool), nil
		}
		return false, nil
	})

	if err != nil {
		return err
	}

	if exists.(bool) {
		return errors.New("curso ya puntuado, intente con otro curso")
	}

	// Crear la relación de puntuación en Neo4j
	session = s.Driver.NewSession(context.TODO(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.TODO())

	_, err = session.ExecuteWrite(context.TODO(), func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
            MATCH (u:Usuario {email: $email}), (c:Curso {id: $cursoID})
            CREATE (u)-[r:PUNTUO {valor: $valor}]->(c)
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

// ObtenerPuntuacionesPorUsuario obtiene todas las puntuaciones hechas por un usuario.
func (s *PuntuacionService) ObtenerPuntuacionesPorUsuario(email string) ([]map[string]interface{}, error) {
	session := s.Driver.NewSession(context.TODO(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(context.TODO())

	result, err := session.ExecuteRead(context.TODO(), func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
            MATCH (u:Usuario {email: $email})-[r:PUNTUO]->(c:Curso)
            RETURN c.nombre AS curso, r.valor AS valoracion
        `
		params := map[string]interface{}{
			"email": email,
		}
		res, err := tx.Run(context.TODO(), query, params)
		if err != nil {
			return nil, err
		}

		var puntuaciones []map[string]interface{}
		for res.Next(context.TODO()) {
			record := res.Record()
			curso, _ := record.Get("curso")
			valoracion, _ := record.Get("valoracion")
			puntuacion := map[string]interface{}{
				"curso":      curso,
				"valoracion": valoracion,
			}
			puntuaciones = append(puntuaciones, puntuacion)
		}
		return puntuaciones, nil
	})

	if err != nil {
		return nil, err
	}

	puntuaciones, ok := result.([]map[string]interface{})
	if !ok {
		return nil, errors.New("error al obtener las puntuaciones")
	}

	return puntuaciones, nil
}
