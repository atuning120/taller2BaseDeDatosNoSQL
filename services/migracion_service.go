// services/migracion_service.go
package services

import (
    "context"
    "encoding/json"
    "fmt"
    "log"

    "go-API/models"

    "github.com/go-redis/redis/v8"
    "github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// MigracionService se encarga de migrar usuarios de Redis a Neo4j
type MigracionService struct {
    RedisClient *redis.Client
    Neo4jDriver neo4j.DriverWithContext
    LogPrefix   string
}

// NewMigracionService crea una nueva instancia de MigracionService
func NewMigracionService(redisClient *redis.Client, neo4jDriver neo4j.DriverWithContext, logPrefix string) *MigracionService {
    return &MigracionService{
        RedisClient: redisClient,
        Neo4jDriver: neo4jDriver,
        LogPrefix:   logPrefix,
    }
}

// MigrarUsuarios migra todos los usuarios de Redis a Neo4j
func (ms *MigracionService) MigrarUsuarios(ctx context.Context) error {
    // 1. Obtener todas las claves de usuarios en Redis
    patron := "usuario:*"
    keys, err := ms.RedisClient.Keys(ctx, patron).Result()
    if err != nil {
        return fmt.Errorf("error al obtener claves de Redis: %v", err)
    }

    log.Printf("[%s] Encontradas %d claves de usuarios en Redis.", ms.LogPrefix, len(keys))

    // 2. Iterar sobre cada clave y migrar el usuario
    for _, key := range keys {
        // Obtener el valor almacenado en la clave
        val, err := ms.RedisClient.Get(ctx, key).Result()
        if err != nil {
            log.Printf("[%s] Error al obtener el valor de la clave %s: %v", ms.LogPrefix, key, err)
            continue
        }

        // Deserializar el usuario
        var usuario models.Usuario
        if err := json.Unmarshal([]byte(val), &usuario); err != nil {
            log.Printf("[%s] Error al deserializar el usuario de la clave %s: %v", ms.LogPrefix, key, err)
            continue
        }

        // Crear el nodo de usuario en Neo4j
        err = ms.crearUsuarioEnNeo4j(ctx, &usuario)
        if err != nil {
            log.Printf("[%s] Error al crear el usuario %s en Neo4j: %v", ms.LogPrefix, usuario.Email, err)
            continue
        }

        log.Printf("[%s] Usuario %s migrado exitosamente.", ms.LogPrefix, usuario.Email)
    }

    return nil
}

// crearUsuarioEnNeo4j crea un nodo User en Neo4j si no existe
func (ms *MigracionService) crearUsuarioEnNeo4j(ctx context.Context, usuario *models.Usuario) error {
    session := ms.Neo4jDriver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
    defer session.Close(ctx)

    _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
        // Verificar si el usuario ya existe
        checkQuery := `
            MATCH (u:User {email: $email})
            RETURN u
        `
        result, err := tx.Run(ctx, checkQuery, map[string]interface{}{
            "email": usuario.Email,
        })
        if err != nil {
            return nil, err
        }

        if result.Next(ctx) {
            // Usuario ya existe, no hacer nada
            return nil, nil
        }

        if err = result.Err(); err != nil {
            return nil, err
        }

        // Crear el nodo User
        createQuery := `
            CREATE (u:User {
                email: $email,
                nombre: $nombre,
                fecha_creacion: datetime()
            })
            RETURN u
        `
        _, err = tx.Run(ctx, createQuery, map[string]interface{}{
            "email":  usuario.Email,
            "nombre": usuario.Nombre,
        })
        if err != nil {
            return nil, err
        }

        return nil, nil
    })

    return err
}
