package server

import (
	// "fmt"
	"log"
	"os"

	"business-connect/router"

	"github.com/joho/godotenv"

	myjwt "business-connect/middleware/myjwt"
)

func LoadEnv() {
	if os.Getenv("RENDER") == "" {
		if err := godotenv.Load(".env"); err != nil {
			log.Println("No .env file found, using system env")
		}
	}
}

func StartServer() {
	LoadEnv()

	// init the JWTs
	jwtErr := myjwt.InitJWT()
	if jwtErr != nil {
		log.Println("Error initializing the JWT's!")
		log.Fatal(jwtErr)
	}

	PORT := os.Getenv("PORT")

	// running all routers in the Routers() function
	routes := router.Routers()

	// running on port "port" local host
	routes.Listen(":" + PORT)
}
