package api

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	controllers "github.com/Wallruzz9114/gen_api/api/controllers"
	seed "github.com/Wallruzz9114/gen_api/api/seed"
)

var server = controllers.Server{}

// Run ...
func Run() {
	var err error
	err = godotenv.Load()

	if err != nil {
		log.Fatalf("Error getting env, not coming through %v", err)
	} else {
		fmt.Println("We are getting the env value")
	}

	server.Initialize(
		os.Getenv("DB_DRIVER"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	seed.Load(server.DB)
	server.Run(":8080")
}
