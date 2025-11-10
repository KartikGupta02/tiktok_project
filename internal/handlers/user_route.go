package route

import (
	"net/http"

	"github.com/kartik/tiktok_project/internal/controller"
	"github.com/kartik/tiktok_project/middleware"
)

func RegisterUserRoutes() {
	http.HandleFunc("/home", controller.GetHtmlData)

	http.HandleFunc("/create_user", controller.CreateUser)
	http.HandleFunc("/get_all_users", controller.GetUsers)
	http.HandleFunc("/users/get", controller.GetUser)
	http.HandleFunc("/users/delete", middleware.AuthMiddleware(controller.DeleteUser))
	http.HandleFunc("/users/update", middleware.AuthMiddleware(controller.UpdateUser))
	http.HandleFunc("/users/login", controller.LoginUser)
	http.HandleFunc("/users/logout", controller.LogoutUser)

}
