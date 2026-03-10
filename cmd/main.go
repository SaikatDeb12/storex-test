package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/SaikatDeb12/storeX/internal/database"
	"github.com/SaikatDeb12/storeX/internal/routes"
	"github.com/SaikatDeb12/storeX/internal/utils"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env")
	}

	err := database.Connect()
	if err != nil {
		log.Fatal("Error connecting to database\n", err)
	}

	router := routes.SetUpRouter()
	serverPort := utils.GetEnvVariables("SERVER_PORT")
	address := fmt.Sprintf(":%s", serverPort)
	fmt.Printf("server running on port %s\n", serverPort)
	log.Fatal(http.ListenAndServe(address, router))
}
