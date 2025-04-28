// controller/list_controller.go
package controller

import (
	"context"
	"encoding/json"
	"net/http"
	// "os"
	"time"
	"todo-list/config"
	"todo-list/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/bson"
)

func CreateList(w http.ResponseWriter, r *http.Request) {
    listCollection := config.GetCollection("lists")

    var list model.List
    err := json.NewDecoder(r.Body).Decode(&list)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    list.ID = primitive.NewObjectID()
    list.Date = time.Now()

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err = listCollection.InsertOne(ctx, list)
    if err != nil {
        http.Error(w, "Failed to create list", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(list)
}

func GetAllLists(w http.ResponseWriter, r *http.Request) {
    listCollection := config.GetCollection("lists")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cursor, err := listCollection.Find(ctx, bson.M{})
    if err != nil {
        http.Error(w, "Error fetching lists", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(ctx)

    var lists []model.List
    for cursor.Next(ctx) {
        var list model.List
        if err := cursor.Decode(&list); err != nil {
            http.Error(w, "Error decoding list", http.StatusInternalServerError)
            return
        }
        lists = append(lists, list)
    }

    // Set content-type and respond with JSON
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(lists); err != nil {
        http.Error(w, "Error encoding response", http.StatusInternalServerError)
    }
}
