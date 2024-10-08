package main

import (
	"context"
	"fmt"
	"github.com/bertoxic/graphqlChat/internal/app"
	"github.com/bertoxic/graphqlChat/pkg/config"
	"log"
)

func main() {
	ctx := context.Background()
	err := app.LoadEnv()
	if err != nil {
		log.Fatal("unable to load .env file from any location 1")
	}

	newConfig, err := config.NewConfig("bert", ":80")
	if err != nil {
		log.Fatalf("unable to create new config: %v", err)
	}

	_, err = app.NewApp(ctx, newConfig)
	if err != nil {
		log.Fatalf("unable to create new app: %v", err)
	}

	//err = app.DB.Migrate()
	//if err != nil {
	//	log.Fatalf("unable to run database migration: %v", err)
	//}

	fmt.Println("Application started successfully")
}
