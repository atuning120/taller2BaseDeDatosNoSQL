package services

import (
    "context"
    "encoding/json"
    "errors"
    "go-API/models"

    "github.com/go-redis/redis/v8"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type UsuarioService struct {
    RedisClient *redis.Client
}

func NewUsuarioService(redisClient *redis.Client) *UsuarioService {
    return &UsuarioService{
        RedisClient: redisClient,
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
    usuario.Inscritos = append(usuario.Inscritos, cursoObjectID)

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

    // Aquí deberías implementar la lógica para obtener los cursos inscritos
    // Por ahora, devolvemos un error indicando que no está implementado
    return nil, errors.New("no implementado")
}