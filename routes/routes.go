package routes

import (
	"github.com/gorilla/mux"
)

func SetupAppRoutes() *mux.Router {
	router := mux.NewRouter()
	// posts.SetupPostsRoutes(router)
	return router
}
