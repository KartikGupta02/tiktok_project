package route

import (
	"net/http"

	"github.com/kartik/tiktok_project/internal/controller"
	"github.com/kartik/tiktok_project/middleware"
)

// RegisterVideoRoutes registers video endpoints and uses the standalone middleware.
func RegisterVideoRoutes() {
	// Protected: require Authorization: Bearer <token>
	http.HandleFunc("/videos/create", middleware.AuthMiddleware(controller.CreateVideo))
	http.HandleFunc("/videos/update", middleware.AuthMiddleware(controller.UpdateVideo))
	http.HandleFunc("/videos/delete", middleware.AuthMiddleware(controller.DeleteVideo))

	// Public: no auth
	http.HandleFunc("/videos/get", controller.GetVideo)
	http.HandleFunc("/videos", controller.ListVideos)
}
