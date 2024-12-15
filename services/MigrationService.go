package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type MigrationService struct {
	Redis   *redis.Client
	MongoDB *mongo.Database
	Neo4j   neo4j.DriverWithContext
}

func NewMigrationService(redisClient *redis.Client, mongoDB *mongo.Database, neo4jDriver neo4j.DriverWithContext) *MigrationService {
	return &MigrationService{Redis: redisClient, MongoDB: mongoDB, Neo4j: neo4jDriver}
}

// MigrateUsuariosYCursos migra usuarios desde Redis y cursos desde MongoDB a Neo4j
func (ms *MigrationService) MigrateUsuariosYCursos(ctx context.Context) error {
	if err := ms.migrateUsuarios(ctx); err != nil {
		return fmt.Errorf("error al migrar usuarios: %v", err)
	}
	if err := ms.migrateCursos(ctx); err != nil {
		return fmt.Errorf("error al migrar cursos: %v", err)
	}
	log.Println("Migración de usuarios y cursos completada")
	return nil
}

func (ms *MigrationService) migrateUsuarios(ctx context.Context) error {
	// Obtener todas las claves de usuarios
	keys, err := ms.Redis.Keys(ctx, "usuario:*:*").Result()
	if err != nil {
		return fmt.Errorf("error al obtener claves de usuarios en Redis: %v", err)
	}

	session := ms.Neo4j.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	for _, key := range keys {
		// Extraer email y password de la clave
		parts := strings.Split(key, ":")
		if len(parts) < 3 {
			log.Printf("Clave de usuario mal formada: %v", key)
			continue
		}
		email := parts[1]
		password := parts[2]

		// Verificar si el usuario ya existe en Neo4j
		exists, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
			query := "MATCH (u:Usuario {email: $email}) RETURN COUNT(u) > 0 AS exists"
			params := map[string]interface{}{"email": email}
			result, err := tx.Run(ctx, query, params)
			if err != nil {
				return false, err
			}
			if result.Next(ctx) {
				return result.Record().Values[0].(bool), nil
			}
			return false, nil
		})
		if err != nil {
			log.Printf("Error al verificar existencia de usuario en Neo4j: %v", err)
			continue
		}

		if exists.(bool) {
			log.Printf("Usuario con email %s ya existe en Neo4j, omitiendo migración", email)
			continue
		}

		// Obtener los datos del usuario
		usuarioJSON, err := ms.Redis.Get(ctx, key).Result()
		if err != nil {
			log.Printf("Error al obtener usuario de Redis: %v", err)
			continue
		}

		var usuario struct {
			Nombre string `json:"nombre"`
		}
		if err := json.Unmarshal([]byte(usuarioJSON), &usuario); err != nil {
			log.Printf("Error al deserializar usuario: %v", err)
			continue
		}

		// Crear nodo de usuario en Neo4j con email, nombre y password
		_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
			query := "CREATE (:Usuario {email: $email, nombre: $nombre, password: $password})"
			params := map[string]interface{}{
				"email":    email,
				"nombre":   usuario.Nombre,
				"password": password,
			}
			_, err := tx.Run(ctx, query, params)
			return nil, err
		})
		if err != nil {
			log.Printf("Error al crear nodo Usuario en Neo4j: %v", err)
		}
	}
	return nil
}


func (ms *MigrationService) migrateCursos(ctx context.Context) error {
	cursosCollection := ms.MongoDB.Collection("cursos")
	cursor, err := cursosCollection.Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("error al obtener cursos de MongoDB: %v", err)
	}
	defer cursor.Close(ctx)

	session := ms.Neo4j.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	for cursor.Next(ctx) {
		var curso struct {
			ID     string `bson:"_id"`
			Nombre string `bson:"nombre"`
		}
		if err := cursor.Decode(&curso); err != nil {
			log.Printf("Error al decodificar curso: %v", err)
			continue
		}

		// Verificar si el curso ya existe en Neo4j
		exists, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
			query := "MATCH (c:Curso {id: $id}) RETURN COUNT(c) > 0 AS exists"
			params := map[string]interface{}{"id": curso.ID}
			result, err := tx.Run(ctx, query, params)
			if err != nil {
				return false, err
			}
			if result.Next(ctx) {
				return result.Record().Values[0].(bool), nil
			}
			return false, nil
		})
		if err != nil {
			log.Printf("Error al verificar existencia de curso en Neo4j: %v", err)
			continue
		}

		if exists.(bool) {
			log.Printf("Curso con ID %s ya existe en Neo4j, omitiendo migración", curso.ID)
			continue
		}

		// Crear nodo de curso en Neo4j
		_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
			query := "CREATE (:Curso {id: $id, nombre: $nombre})"
			params := map[string]interface{}{"id": curso.ID, "nombre": curso.Nombre}
			_, err := tx.Run(ctx, query, params)
			return nil, err
		})
		if err != nil {
			log.Printf("Error al crear nodo Curso en Neo4j: %v", err)
		}
	}
	return nil
}
