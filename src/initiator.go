package initiator

import (
	"context"
	"cloud.google.com/go/firestore"
	"doctors_care/src/handlers"
	"doctors_care/src/utils"
	// "github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"

	"log"
)

func Initialize() {

	e := echo.New()
	e.Use(middleware.CORS())

	var err error
	utils.FBClient, err = initializeFireStore()

	if err != nil {
		log.Fatalln(err)
	}

	defer utils.FBClient.Close()

	// err = godotenv.Load(".env")

	// if err != nil {
	// 	log.Fatalf("Error loading .env file")
	// }

	e.Static("/", "src/build/")

	api := e.Group("/api")
	// api.Use(middleware.JWT([]byte(os.Getenv("HASH_SECRET"))))

	api.GET("/getADoctor", handlers.GetADoctor)
	e.Logger.Fatal(e.Start(":1323"))
}

func initializeFireStore() (*firestore.Client, error) {
	utils.Ctx = context.Background()
	sa := option.WithCredentialsFile("./serviceAccountKey.json")
	app, err := firebase.NewApp(utils.Ctx, nil, sa)

	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(utils.Ctx)
	if err != nil {
		return nil, err
	}

	return client, nil

}
