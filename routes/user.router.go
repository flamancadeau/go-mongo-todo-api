// routes/user.router.go
package routes

import (
    "net/http"
    "todo-list/controller"
)

func RegisterUserRoutes() {
    http.HandleFunc("/users", controller.CreateUser)
}
