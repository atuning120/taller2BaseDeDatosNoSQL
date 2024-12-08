package services

import (
	"context"
	"encoding/json"
	"errors"
	"go-API/models"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsuarioService struct {
	RedisClient     *redis.Client
	CursoCollection *mongo.Collection
}

func NewUsuarioService(redisClient *redis.Client, cursoCollection *mongo.Collection) *UsuarioService {
    return &UsuarioService{
        RedisClient:     redisClient,
        CursoCollection: cursoCollection,
    }
}

func (us *UsuarioService) ObtenerUsuarios() ([]models.Usuario, error) {
	var usuarios []models.Usuario

	keys, err := us.RedisClient.Keys(context.TODO(), "usuario:*").Result()
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		val, err := us.RedisClient.Get(context.TODO(), key).Result()
		if err != nil {
			return nil, err
		}

		var usuario models.Usuario
		if err := json.Unmarshal([]byte(val), &usuario); err != nil {
			return nil, err
		}

		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

func (us *UsuarioService) CrearUsuario(usuario *models.Usuario) (string, error) {
	key := "usuario:" + usuario.Email + ":" + usuario.Password

	data, err := json.Marshal(usuario)
	if err != nil {
		return "", err
	}

	err = us.RedisClient.Set(context.TODO(), key, data, 0).Err()
	if err != nil {
		return "", err
	}

	return key, nil
}

func (us *UsuarioService) ObtenerUsuarioPorCorreoYContrasena(email, password string) (*models.Usuario, error) {
	key := "usuario:" + email + ":" + password

	val, err := us.RedisClient.Get(context.TODO(), key).Result()
	if err == redis.Nil {
		return nil, errors.New("usuario no encontrado")
	} else if err != nil {
		return nil, err
	}

	var usuario models.Usuario
	if err := json.Unmarshal([]byte(val), &usuario); err != nil {
		return nil, err
	}

	return &usuario, nil
}

func (us *UsuarioService) InscribirseACurso(email, password, cursoID string) error {
	key := "usuario:" + email + ":" + password

	val, err := us.RedisClient.Get(context.TODO(), key).Result()
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
		return err
	}

	// Verificar si el usuario ya está inscrito en el curso
	for _, inscrito := range usuario.Inscritos {
		if inscrito == cursoObjectID {
			return errors.New("el usuario ya está inscrito en este curso")
		}
	}

	// Agregar el curso a Inscritos y la fecha de inscripción
	usuario.Inscritos = append(usuario.Inscritos, cursoObjectID)
	usuario.FechaInscripcion = append(usuario.FechaInscripcion, time.Now())

	// Crear y agregar el ProgresoCurso
	nuevoProgreso := models.ProgresoCurso{
		CursoID:      cursoObjectID,
		ClasesVistas: []primitive.ObjectID{},
		Estado:       "INICIADO",
	}
	usuario.Progresos = append(usuario.Progresos, nuevoProgreso)

	data, err := json.Marshal(usuario)
	if err != nil {
		return err
	}

	err = us.RedisClient.Set(context.TODO(), key, data, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (us *UsuarioService) ObtenerCursosInscritos(email, password string) ([]models.Curso, error) {
    // Construir la clave de Redis
    key := "usuario:" + email + ":" + password

    // Obtener el usuario desde Redis
    val, err := us.RedisClient.Get(context.TODO(), key).Result()
    if err == redis.Nil {
        return nil, errors.New("usuario no encontrado")
    } else if err != nil {
        return nil, err
    }

    // Deserializar el usuario
    var usuario models.Usuario
    if err := json.Unmarshal([]byte(val), &usuario); err != nil {
        return nil, err
    }

    // Verificar si el usuario tiene cursos inscritos
    if len(usuario.Inscritos) == 0 {
        return []models.Curso{}, nil // Retorna un slice vacío
    }

    // Convertir los IDs de cursos a ObjectID si es necesario
    var objectIDs []primitive.ObjectID
    for _, id := range usuario.Inscritos {
        objectIDs = append(objectIDs, id)
    }

    // Realizar una sola consulta para obtener todos los cursos inscritos
    cursor, err := us.CursoCollection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": objectIDs}})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.TODO())

    var cursos []models.Curso
    for cursor.Next(context.TODO()) {
        var curso models.Curso
        if err := cursor.Decode(&curso); err != nil {
            return nil, err
        }
        cursos = append(cursos, curso)
    }

    if err := cursor.Err(); err != nil {
        return nil, err
    }

    return cursos, nil
}

// VerClase permite que un usuario vea una clase y actualiza su progreso en el curso.
func (s *UsuarioService) VerClase(email, password, claseID string) error {
	// Obtener el usuario
	usuario, err := s.ObtenerUsuarioPorCorreoYContrasena(email, password)
	if err != nil {
		return err
	}

	// Convertir claseID a ObjectID
	claseObjectID, err := primitive.ObjectIDFromHex(claseID)
	if err != nil {
		return errors.New("ID de clase inválido")
	}

	// Verificar si la clase pertenece a algún curso inscrito por el usuario
	var cursoID primitive.ObjectID
	claseEncontrada := false
	for _, progreso := range usuario.Progresos {
		totalClases, err := s.obtenerTotalClasesPorCurso(progreso.CursoID)
		if err != nil {
			return err
		}

		for _, clase := range totalClases {
			if clase == claseObjectID {
				cursoID = progreso.CursoID
				claseEncontrada = true
				break
			}
		}

		if claseEncontrada {
			break
		}
	}

	if !claseEncontrada {
		return errors.New("la clase no pertenece a ningún curso inscrito por el usuario")
	}

	// Actualizar el progreso del curso
	for i, progreso := range usuario.Progresos {
		if progreso.CursoID == cursoID {
			if contains(progreso.ClasesVistas, claseObjectID) {
				return errors.New("clase ya vista")
			}
			progreso.ClasesVistas = append(progreso.ClasesVistas, claseObjectID)

			// Obtener el total de clases del curso
			totalClases, err := s.obtenerTotalClasesPorCurso(progreso.CursoID)
			if err != nil {
				return err
			}

			// Actualizar el estado del progreso
			if len(progreso.ClasesVistas) == 0 {
				progreso.Estado = "INICIADO"
			} else if len(progreso.ClasesVistas) < len(totalClases) {
				progreso.Estado = "EN CURSO"
			} else if len(progreso.ClasesVistas) == len(totalClases) {
				progreso.Estado = "COMPLETADO"
			}

			usuario.Progresos[i] = progreso
			break
		}
	}

	// Actualizar el usuario en la base de datos Redis
	data, err := json.Marshal(usuario)
	if err != nil {
		return err
	}

	key := "usuario:" + email + ":" + password
	err = s.RedisClient.Set(context.TODO(), key, data, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// obtenerTotalClasesPorCurso obtiene el número total de clases de un curso.
func (s *UsuarioService) obtenerTotalClasesPorCurso(cursoID primitive.ObjectID) ([]primitive.ObjectID, error) {
	var curso models.Curso
	err := s.CursoCollection.FindOne(context.TODO(), bson.M{"_id": cursoID}).Decode(&curso)
	if err != nil {
		return nil, err
	}
	return curso.Unidades, nil
}

// contains verifica si una lista de ObjectIDs contiene un ObjectID específico.
func contains(slice []primitive.ObjectID, item primitive.ObjectID) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// ObtenerProgresoCursos obtiene el progreso de los cursos en los que un usuario está inscrito.
func (s *UsuarioService) ObtenerProgresoCursos(email, password string) ([]models.ProgresoCurso, error) {
	// Obtener el usuario
	usuario, err := s.ObtenerUsuarioPorCorreoYContrasena(email, password)
	if err != nil {
		return nil, err
	}

	return usuario.Progresos, nil
}
