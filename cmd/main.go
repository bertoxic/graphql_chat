package main

import (
	"context"
	"fmt"
	"github.com/bertoxic/graphqlChat/internal/app"
	errorx "github.com/bertoxic/graphqlChat/internal/error"
	"github.com/bertoxic/graphqlChat/pkg/config"
	"github.com/bertoxic/graphqlChat/router"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()
	// err := app.LoadEnv()
	err := app.LoadEnvInProd()
	if err != nil {
		log.Fatal("unable to load .env file from any location 1")
	}

	newConfig, err := config.NewConfig(".env", ":80")
	if err != nil {
		log.Fatalf("unable to create new config: %v", err)
	}

	app, err := app.NewApp(ctx, newConfig)
	if err != nil {
		log.Fatalf("unable to create new app: %v", err)
	}

	// err = app.Repo.Migrate()
	// if err != nil {
	//	log.Fatalf("unable to run database migration: %v", err)
	//}
	fmt.Println("Application started successfully")
	// http.HandleFunc("/", handlers.Repo.HomePage)
	mux := router.Routes(app)
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		err = errorx.New(errorx.ErrInternal.Code, "", err)
		// fmt.Printf("%v", err.(*errorx.AppError).Details)
	}

}
