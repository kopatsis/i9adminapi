package main

import (
	"context"
	"encoding/base64"
	"i9-adminapi/platform"
	"i9-adminapi/shared"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func main() {
	if err := godotenv.Load(); err != nil {
		if os.Getenv("APP_ENV") != "production" {
			log.Fatalf("Failed to load the env vars: %v", err)
		}
	}

	client, database, err := shared.ConnectDB()
	if err != nil {
		log.Fatalf("Error while connecting to mongoDB: %s.\nExiting.", err)
	}
	defer shared.DisConnectDB(client)

	firebaseConfigBase64 := os.Getenv("FIREBASE_CONFIG_BASE64")
	if firebaseConfigBase64 == "" {
		log.Fatal("FIREBASE_CONFIG_BASE64 environment variable is not set.")
	}

	configJSON, err := base64.StdEncoding.DecodeString(firebaseConfigBase64)
	if err != nil {
		log.Fatalf("Error decoding FIREBASE_CONFIG_BASE64: %v", err)
	}

	sa := option.WithCredentialsJSON(configJSON)
	firebase, err := firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	rtr := platform.New(database, firebase)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := http.ListenAndServe(":"+port, rtr); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}
