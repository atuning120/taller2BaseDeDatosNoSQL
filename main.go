package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-API/controllers"
	_ "go-API/docs" // Importar los documentos de Swagger
	"go-API/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/go-redis/redis/v8"
)

// @title API de Cursos y Usuarios
// @version 1.0
// @description Esta es una API para gestionar cursos y usuarios.

// @host localhost:8080
// @BasePath /

var mongoClient *mongo.Client

func init() {
	if err := loadEnv(); err != nil {
		log.Fatal("Error al cargar las variables de entorno:", err)
	}
	if err := connectToMongoDB(); err != nil {
		log.Fatal("No se pudo conectar a MongoDB:", err)
	}
}

func main() {
	router := gin.Default()

	redisClient := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

	// Enlace con Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Inicializar servicios y controladores
	db := mongoClient.Database("miBaseDeDatos")
	cursoService := services.NewCursoService(db)
	cursoControlador := controllers.NewCursoControlador(cursoService)

	unidadService := services.NewUnidadService(db)
	unidadControlador := controllers.NewUnidadControlador(unidadService)

	claseService := services.NewClaseService(db)
	claseControlador := controllers.NewClaseControlador(claseService)

	usuarioService := services.NewUsuarioService(redisClient,db.Collection("cursos"))
    usuarioControlador := controllers.NewUsuarioControlador(usuarioService)

	comentarioService := services.NewComentarioService(db)
	comentarioControlador := controllers.NewComentarioControlador(comentarioService)

	// Rutas de la API
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Conexión exitosa"})
	})

	// Cursos
	router.GET("/api/cursos", cursoControlador.ObtenerCursos)
	router.GET("/api/cursos/:id", cursoControlador.ObtenerCursoPorID)
	router.PATCH("/api/cursos/:id/valoracion", cursoControlador.ActualizarValoracion)
	router.POST("/api/cursos", cursoControlador.CrearCurso)
	router.GET("/api/cursos/:id/clases", cursoControlador.ObtenerClasesPorCurso)

	// Unidades
	router.GET("/api/cursos/:id/unidades", unidadControlador.ObtenerUnidadesPorCurso)
	router.POST("/api/cursos/:id/unidades", unidadControlador.CrearUnidad)

	// Clases
	router.GET("/api/unidades/:id/clases", claseControlador.ObtenerClasesPorUnidad)
	router.POST("/api/unidades/:id/clases", claseControlador.CrearClaseParaUnidad)

	// Comentarios
	router.GET("/api/clases/:id/comentarios", comentarioControlador.ObtenerComentariosPorClase)
	router.POST("/api/clases/:id/comentarios", comentarioControlador.CrearComentarioParaClase)

	// Usuarios
	router.GET("/api/usuarios", usuarioControlador.ObtenerUsuarios)
	router.GET("/api/usuarios/usuario", usuarioControlador.ObtenerUsuarioPorCorreoYContrasena)
	router.GET("/api/usuarios/cursos", usuarioControlador.ObtenerCursosInscritos)
	router.POST("/api/usuarios", usuarioControlador.CrearUsuario)
	router.POST("/api/usuarios/inscripcion", usuarioControlador.InscribirseACurso)
	router.POST("/api/usuarios/ver_clase/:clase_id", usuarioControlador.VerClase)
	router.GET("/api/usuarios/progreso", usuarioControlador.ObtenerProgresoCursos)

	// Iniciar el servidor
	go func() {
		if err := router.Run(); err != nil {
			log.Fatal("Error al iniciar el servidor:", err)
		}
	}()

	gracefulShutdown()
}

func connectToMongoDB() error {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return fmt.Errorf("la URI de MongoDB no está configurada")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return err
	}

	if err = client.Ping(context.TODO(), nil); err != nil {
		return err
	}

	log.Println("Conectado a MongoDB")
	mongoClient = client
	return nil
}

func loadEnv() error {
	return godotenv.Load()
}

func gracefulShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Cerrando la conexión con MongoDB...")
	if err := mongoClient.Disconnect(context.TODO()); err != nil {
		log.Fatal("Error al desconectar MongoDB:", err)
	}
	log.Println("Conexión con MongoDB cerrada. Apagando servidor.")
	os.Exit(0)
}
