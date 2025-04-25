package routes

import (
	"net/http"
	"strings"
	"todo-list/controller"
)

func RegisterUserRoutes() {
	// Signup
	http.HandleFunc("/auth/signup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			controller.CreateUser(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Login
	http.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			controller.Login(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})


	// Get, update, or delete user by ID
	http.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "User ID is missing", http.StatusBadRequest)
			return
		}
		userID := parts[len(parts)-1]

		switch r.Method {
		case http.MethodGet:
			controller.GetUserByID(w, r, userID)
		case http.MethodPut:
			controller.UpdateUser(w, r, userID)
		case http.MethodDelete:
			controller.DeleteUser(w, r, userID)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
