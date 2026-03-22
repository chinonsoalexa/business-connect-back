package server

import (
	"log"
	"os"

	"business-connect/router"

	// "github.com/joho/godotenv"

	myjwt "business-connect/middleware/myjwt"
	emailValid "business-connect/email"
)

func LoadEnv() {
	// if os.Getenv("RENDER") == "" {
	// 	if err := godotenv.Load(".env"); err != nil {
	// 		log.Println("No .env file found, using system env")
	// 	}
	// }

		// 3️⃣ Disposable email check
	err := emailValid.LoadDisposableList("fakeEmails.json")
	if err != nil {
			log.Println("Email validation service error")
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
