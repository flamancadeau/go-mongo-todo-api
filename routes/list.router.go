// routes/list.router.go
package routes

import (
    "net/http"
    "todo-list/controller"
)

func RegisterListRoutes() {
    http.HandleFunc("/lists", controller.CreateList)
}
