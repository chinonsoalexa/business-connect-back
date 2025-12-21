package server

import (
	"fmt"
	"log"
	"os"

	"business-connect/router"

	"github.com/joho/godotenv"

	myjwt "business-connect/middleware/myjwt"
)

func StartServer() {
	envErr := godotenv.Load(".env")

	if envErr != nil {
		fmt.Println("Failed to load .env file: %w", envErr)
	}

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
