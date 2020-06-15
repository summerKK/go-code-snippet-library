package api

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/summerKK/go-code-snippet-library/champion-go/api/controllers"
)

var server = &controllers.Server{}

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("sad .env file found")
	}
}

func Run() {

	server.Initialize(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	// This is for testing, when done, do well to comment
	// seed.Load(server.DB)

	apiPort := fmt.Sprintf(":%s", os.Getenv("API_PORT"))
	fmt.Printf("Listening to port %s", apiPort)

	server.Run(apiPort)
}
