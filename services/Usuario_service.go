package services

import (
	"context"
	"errors"
	"go-API/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UnidadService maneja la lógica relacionada con las unidades.
type UsuarioService struct {
	UsuarioCollection *mongo.Collection
	CursoCollection   *mongo.Collection
}

// NewUnidadService crea un nuevo servicio para las unidades.
func NewUsuarioService(db *mongo.Database) *UsuarioService {
	return &UsuarioService{
		UsuarioCollection: db.Collection("usuarios"),
		CursoCollection:   db.Collection("cursos"),
	}
}

// ObtenerUsuarios obtiene todos los usuarios.
func (us *UsuarioService) ObtenerUsuarios() ([]models.Usuario, error) {
	cursor, err := us.UsuarioCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	var usuarios []models.Usuario
	if err = cursor.All(context.TODO(), &usuarios); err != nil {
		return nil, err
	}

	return usuarios, nil
}

// CrearUsuario agrega un nuevo usuario a la base de datos.
func (s *UsuarioService) CrearUsuario(usuario *models.Usuario) (*mongo.InsertOneResult, error) {
	// Asegurarse de que la lista de inscritos esté inicializada
	if usuario.Inscritos == nil {
		usuario.Inscritos = []primitive.ObjectID{}
	}

	usuario.ID = primitive.NewObjectID() // Generar un nuevo ID para el usuario

	// Insertar el usuario en la base de datos
	result, err := s.UsuarioCollection.InsertOne(context.TODO(), usuario)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ObtenerUsuarioPorID obtiene un usuario por su ID.
func (us *UsuarioService) ObtenerUsuarioPorID(id string) (*models.Usuario, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var usuario models.Usuario
	err = us.UsuarioCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&usuario)
	if err != nil {
		return nil, err
	}

	return &usuario, nil
}

// InscribirseACurso permite que un usuario se inscriba en un curso.
func (s *UsuarioService) InscribirseACurso(usuarioID, cursoID string) error {
	// Convertir IDs a ObjectID
	usuarioObjectID, err := primitive.ObjectIDFromHex(usuarioID)
	if err != nil {
		return errors.New("ID de usuario inválido")
	}

	cursoObjectID, err := primitive.ObjectIDFromHex(cursoID)
	if err != nil {
		return errors.New("ID de curso inválido")
	}

	// Verificar si el usuario existe
	var usuario models.Usuario
	err = s.UsuarioCollection.FindOne(context.TODO(), bson.M{"_id": usuarioObjectID}).Decode(&usuario)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("usuario no encontrado")
		}
		return err
	}

	// Verificar si el curso existe
	var curso models.Curso
	err = s.CursoCollection.FindOne(context.TODO(), bson.M{"_id": cursoObjectID}).Decode(&curso)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("curso no encontrado")
		}
		return err
	}

	// Actualizar el usuario agregando el ID del curso a la lista de inscritos
	_, err = s.UsuarioCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": usuarioObjectID},
		bson.M{"$addToSet": bson.M{"inscritos": cursoObjectID}}, // Evita duplicados
	)
	if err != nil {
		return err
	}

	// Incrementar el contador de usuarios inscritos en el curso
	_, err = s.CursoCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": cursoObjectID},
		bson.M{"$inc": bson.M{"cant_usuarios": 1}},
	)
	if err != nil {
		return err
	}

	return nil
}

// ObtenerCursosInscritos obtiene los cursos en los que un usuario está inscrito.
func (s *UsuarioService) ObtenerCursosInscritos(usuarioID string) ([]models.Curso, error) {
	// Convertir el ID a ObjectID
	usuarioObjectID, err := primitive.ObjectIDFromHex(usuarioID)
	if err != nil {
		return nil, errors.New("ID de usuario inválido")
	}

	// Obtener el usuario
	var usuario models.Usuario
	err = s.UsuarioCollection.FindOne(context.TODO(), bson.M{"_id": usuarioObjectID}).Decode(&usuario)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	// Obtener los cursos en los que está inscrito el usuario
	cursor, err := s.CursoCollection.Find(
		context.TODO(),
		bson.M{"_id": bson.M{"$in": usuario.Inscritos}},
	)
	if err != nil {
		return nil, err
	}

	var cursos []models.Curso
	if err = cursor.All(context.TODO(), &cursos); err != nil {
		return nil, err
	}

	return cursos, nil
}
