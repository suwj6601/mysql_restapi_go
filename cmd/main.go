package main

import (
	"fmt"
	"go-mysql-restapi/api"
	db "go-mysql-restapi/database"
	"go-mysql-restapi/service/auth"

	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize the database
	db := db.InitDB()

	// Setup app routers
	r := mux.NewRouter()
	r.Use(auth.AuthMiddleware)
	api := &api.Api{Router: r, DB: db}
	api.SetUpRoutes()

	// Start server
	fmt.Println("Server running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))

}
