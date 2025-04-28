// routes/list.router.go
package routes

import (
    "net/http"
    "todo-list/controller"
)

func RegisterListRoutes() {
    http.HandleFunc("/api/lists", controller.CreateList)
    http.HandleFunc("/api/getlist", controller.CreateList)
}
