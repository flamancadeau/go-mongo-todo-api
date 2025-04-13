package main

import (
    "fmt"
    "log"
    "net/http"
    "todo-list/config"
    "todo-list/routes"
)

func main() {
    config.ConnectDB()

    // Register all routes
    routes.RegisterUserRoutes()
    routes.RegisterListRoutes()

    fmt.Println("🚀 Server running at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
