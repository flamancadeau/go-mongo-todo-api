package main

import (
    "fmt"
    "log"
    "net/http"
    "todo-list/config"
    "todo-list/routes"
    "github.com/joho/godotenv"
)


func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome already was connected sucessfull!")
}


func main() {

       // Load .env file
       err := godotenv.Load()
       if err != nil {
           log.Fatal("Error loading .env file")
       }
   
    config.ConnectDB()

    // Register all routes
    routes.RegisterUserRoutes()
    routes.RegisterListRoutes()

    	// testing  handler for the root URL
	http.HandleFunc("/", welcomeHandler)


    fmt.Println("🚀 Server running at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
