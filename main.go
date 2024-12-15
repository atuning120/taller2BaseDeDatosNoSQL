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
    "go-API/neo4j"

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
var redisClient *redis.Client

func init() {
    if err := loadEnv(); err != nil {
        log.Fatal("Error al cargar las variables de entorno:", err)
    }
    if err := connectToMongoDB(); err != nil {
        log.Fatal("No se pudo conectar a MongoDB:", err)
    }
    // Inicializar Neo4j
    neo4j.InitNeo4j()
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
    migrationService := services.NewMigrationService(redisClient, db, neo4j.Driver)
    
    cursoService := services.NewCursoService(db, neo4j.Driver)
    cursoControlador := controllers.NewCursoControlador(cursoService)

    unidadService := services.NewUnidadService(db)
    unidadControlador := controllers.NewUnidadControlador(unidadService)

    claseService := services.NewClaseService(db)
    claseControlador := controllers.NewClaseControlador(claseService)

    usuarioService := services.NewUsuarioService(redisClient,db.Collection("cursos"),db.Collection("unidades"),db.Collection("clases"),neo4j.Driver)
    usuarioControlador := controllers.NewUsuarioControlador(usuarioService)

    comentarioService := services.NewComentarioService(neo4j.Driver, redisClient)
    comentarioControlador := controllers.NewComentarioControlador(comentarioService)

    comentarioCursoService := services.NewComentarioCursoService(neo4j.Driver)
    comentarioCursoControlador := controllers.NewComentarioCursoControlador(comentarioCursoService)

    puntuacionService := services.NewPuntuacionService(neo4j.Driver, db.Collection("cursos"),redisClient)
    puntuacionesControlador := controllers.NewPuntuacionesControlador(puntuacionService)

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
    router.POST("/api/usuarios/:email/:password/clases/:clase_id", usuarioControlador.VerClase)
    router.GET("/api/usuarios/progreso", usuarioControlador.ObtenerProgresoCursos)

    // Puntuaciones
    router.POST("/api/puntuaciones/cursos/:id", puntuacionesControlador.CrearPuntuacionParaCurso)
    router.GET("/api/puntuaciones/cursos/:id/promedio", puntuacionesControlador.ObtenerPromedioPuntuacion)
    router.GET("/api/puntuaciones/usuarios/:email", puntuacionesControlador.ObtenerPuntuacionesPorUsuario)	

    // Comentarios de Curso
    router.POST("/api/comentarios_curso", comentarioCursoControlador.CrearComentarioCurso)
    router.GET("/api/comentarios_curso/usuarios/:email", comentarioCursoControlador.ObtenerComentariosCursoPorUsuario)

    // Migraciones de usuarios y cursos a nodos en el grafo de Neo4j [hacer en postman]
    router.POST("/api/migrate", func(c *gin.Context) {
        if err := migrationService.MigrateUsuariosYCursos(context.Background()); err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        c.JSON(200, gin.H{"message": "Migración completada exitosamente"})
    })

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
    log.Println("Cerrando la conexión con Neo4j...")
    neo4j.CloseNeo4j()
    log.Println("Cerrando la conexión con Redis...")
    // Cerrar Redis
    redisClient.Close()
    log.Println("Conexión con MongoDB, Neo4j y Redis cerrada. Apagando servidor.")
    os.Exit(0)
}