package services

import (
    "context"
    "errors"

    "github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type ComentarioCursoService struct {
    Driver neo4j.DriverWithContext
}

func NewComentarioCursoService(driver neo4j.DriverWithContext) *ComentarioCursoService {
    return &ComentarioCursoService{
        Driver: driver,
    }
}

// CrearComentarioCurso crea un nuevo comentario para un curso.
func (s *ComentarioCursoService) CrearComentarioCurso(email, cursoID, texto string) error {
    if len(texto) < 15 {
        return errors.New("el comentario debe tener al menos 15 caracteres")
    }

    session := s.Driver.NewSession(context.TODO(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
    defer session.Close(context.TODO())

    _, err := session.ExecuteWrite(context.TODO(), func(tx neo4j.ManagedTransaction) (interface{}, error) {
        query := `
            MATCH (u:Usuario {email: $email}), (c:Curso {id: $cursoID})
            CREATE (u)-[r:REALIZO_COMENTARIO {texto: $texto}]->(c)
            RETURN r
        `
        params := map[string]interface{}{
            "email":   email,
            "cursoID": cursoID,
            "texto":   texto,
        }
        _, err := tx.Run(context.TODO(), query, params)
        return nil, err
    })

    return err
}

// ObtenerComentariosCursoPorUsuario obtiene todos los comentarios hechos por un usuario.
func (s *ComentarioCursoService) ObtenerComentariosCursoPorUsuario(email string) ([]map[string]interface{}, error) {
    session := s.Driver.NewSession(context.TODO(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
    defer session.Close(context.TODO())

    result, err := session.ExecuteRead(context.TODO(), func(tx neo4j.ManagedTransaction) (interface{}, error) {
        query := `
            MATCH (u:Usuario {email: $email})-[r:REALIZO_COMENTARIO]->(c:Curso)
            RETURN c.nombre AS curso, r.texto AS comentario
        `
        params := map[string]interface{}{
            "email": email,
        }
        res, err := tx.Run(context.TODO(), query, params)
        if err != nil {
            return nil, err
        }

        var comentarios []map[string]interface{}
        for res.Next(context.TODO()) {
            record := res.Record()
            curso, _ := record.Get("curso")
            comentario, _ := record.Get("comentario")
            comentarioCurso := map[string]interface{}{
                "curso":      curso,
                "comentario": comentario,
            }
            comentarios = append(comentarios, comentarioCurso)
        }
        return comentarios, nil
    })

    if err != nil {
        return nil, err
    }

    comentarios, ok := result.([]map[string]interface{})
    if !ok {
        return nil, errors.New("error al obtener los comentarios")
    }

    return comentarios, nil
}