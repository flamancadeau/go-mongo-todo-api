package config

import (
    "context"
    "log"
    "os"
    "time"

    "github.com/joho/godotenv"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
    if err != nil {
        log.Fatal(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }

    DB = client.Database(os.Getenv("DB_NAME"))
    log.Println("Connected to MongoDB!")
}

func GetCollection(collectionName string) *mongo.Collection {
    return DB.Collection(collectionName)
}
