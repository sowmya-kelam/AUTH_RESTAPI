package main

import (
	"restapi/db"

	"restapi/routes"
	
	 "github.com/ukautz/clif"

	//"database/sql"

	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gin-gonic/gin"
)


//	@title			Auth Rest API's
//	@version		1.0
//	@description	Auth REST API documentation
func main() {
    // Create a new CLI application
    c := clif.New("MyApp", "1.0", "An example CLI app using clif")

    // Add a command named "start" with a simple function to run the server
    c.Add(clif.NewCommand("start", "Start the server", func() {
        mainServerLogic()
    }))

    // Run the CLI
    c.Run()
}

// server startup logic
func mainServerLogic() {
    database, err := db.InitDB()
    if err != nil {
        panic(fmt.Sprintf("Failed to connect to the database: %v", err))
    }
    defer database.Close()

    server := gin.Default()
    routes.RegisterRoutes(server, database)
    server.Run(":8080")
}