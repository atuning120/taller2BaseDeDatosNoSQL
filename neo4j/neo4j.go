// neo4j/neo4j.go
package neo4j

import (
    "context"
    "log"
    "os"
    "time"

    "github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Driver es el driver global de Neo4j
var Driver neo4j.DriverWithContext

// InitNeo4j inicializa la conexión a Neo4j.
func InitNeo4j() {
    uri := os.Getenv("NEO4J_URI")           // e.g., bolt://localhost:7687
    username := os.Getenv("NEO4J_USER")     // e.g., neo4j
    password := os.Getenv("NEO4J_PASSWORD") // e.g., password

    var err error
    Driver, err = neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
    if err != nil {
        log.Fatalf("Error al crear el driver de Neo4j: %v", err)
    }

    // Verificar la conexión con un timeout de 5 segundos
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    err = Driver.VerifyConnectivity(ctx)
    if err != nil {
        log.Fatalf("Error al conectar con Neo4j: %v", err)
    }

    log.Println("Conectado exitosamente a Neo4j")
}

// CloseNeo4j cierra la conexión a Neo4j.
func CloseNeo4j() {
    if Driver != nil {
        Driver.Close(context.Background())
    }
}
