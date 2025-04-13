// controller/list_controller.go
package controller

import (
    "context"
    "encoding/json"
    "net/http"
    "todo-list/config"
    "todo-list/model"
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
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
