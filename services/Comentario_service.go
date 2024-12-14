// services/comentario_service.go
package services

import (
    "context"
    "errors"
    "go-API/models"
    "time"

    "github.com/google/uuid"
    "github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type ComentarioService struct {
    Driver neo4j.DriverWithContext
}

func NewComentarioService(driver neo4j.DriverWithContext) *ComentarioService {
    return &ComentarioService{Driver: driver}
}

// ObtenerComentariosPorClase obtiene todos los comentarios asociados a una clase.
func (s *ComentarioService) ObtenerComentariosPorClase(ctx context.Context, claseID string) ([]models.Comentario, error) {
    session := s.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
    defer session.Close(ctx)

    result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
        query := `
            MATCH (u:User)-[:COMENTÓ]->(c:Comment)-[:PERTENECE_A]->(:Course)-[:CONTENEDOR_DE]->(:Clase {id: $claseID})
            RETURN c.id AS id, c.autor AS autor, c.fecha AS fecha, c.titulo AS titulo, c.detalle AS detalle, c.meGusta AS meGusta, c.noMeGusta AS noMeGusta
            ORDER BY c.fecha DESC
        `
        records, err := tx.Run(ctx, query, map[string]interface{}{
            "claseID": claseID,
        })
        if err != nil {
            return nil, err
        }

        var comentarios []models.Comentario
        for records.Next(ctx) {
            record := records.Record()
            comentario := models.Comentario{
                ID:        record.Values[0].(string),
                Autor:     record.Values[1].(string),
                Fecha:     record.Values[2].(time.Time),
                Titulo:    record.Values[3].(string),
                Detalle:   record.Values[4].(string),
                MeGusta:   int(record.Values[5].(int64)),
                NoMeGusta: int(record.Values[6].(int64)),
            }
            comentarios = append(comentarios, comentario)
        }

        if err = records.Err(); err != nil {
            return nil, err
        }

        return comentarios, nil
    })

    if err != nil {
        return nil, err
    }

    return result.([]models.Comentario), nil
}

// CrearComentarioParaClase crea un nuevo comentario asociado a una clase.
func (s *ComentarioService) CrearComentarioParaClase(ctx context.Context, claseID string, comentario *models.Comentario) (*models.Comentario, error) {
    session := s.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
    defer session.Close(ctx)

    result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
        // Verificar que la clase existe
        checkQuery := `
            MATCH (clase:Clase {id: $claseID})
            RETURN COUNT(*) > 0 AS exists
        `
        checkResult, err := tx.Run(ctx, checkQuery, map[string]interface{}{
            "claseID": claseID,
        })
        if err != nil {
            return nil, err
        }
        if !checkResult.Next(ctx) {
            return nil, errors.New("clase no encontrada")
        }
        exists, ok := checkResult.Record().Get("exists")
        if !ok || !exists.(bool) {
            return nil, errors.New("clase no encontrada")
        }

        // Generar un ID único para el comentario
        comentarioID := uuid.New().String()
        comentario.ID = comentarioID
        comentario.Fecha = time.Now()

        // Crear el comentario y las relaciones
        createQuery := `
            MATCH (u:User {email: $autor}), (clase:Clase {id: $claseID})
            CREATE (u)-[:COMENTÓ]->(c:Comment {
                id: $id,
                autor: $autor,
                fecha: datetime(),
                titulo: $titulo,
                detalle: $detalle,
                meGusta: $meGusta,
                noMeGusta: $noMeGusta
            })-[:PERTENECE_A]->(:Course)-[:CONTENEDOR_DE]->(clase)
            RETURN c.id AS id, c.autor AS autor, c.fecha AS fecha, c.titulo AS titulo, c.detalle AS detalle, c.meGusta AS meGusta, c.noMeGusta AS noMeGusta
        `
        // Primero obtenemos el ResultWithContext y el error
        res, err := tx.Run(ctx, createQuery, map[string]interface{}{
            "id":        comentario.ID,
            "autor":     comentario.Autor,
            "titulo":    comentario.Titulo,
            "detalle":   comentario.Detalle,
            "meGusta":   comentario.MeGusta,
            "noMeGusta": comentario.NoMeGusta,
            "claseID":   claseID,
        })
        if err != nil {
            return nil, err
        }

        // Ahora consumimos el resultado
        _, err = res.Consume(ctx)
        if err != nil {
            return nil, err
        }

        return comentario, nil
    })

    if err != nil {
        return nil, err
    }

    return result.(*models.Comentario), nil
}

