package routes

import (
	"oauth2/modules/auth"

	"github.com/gorilla/mux"
)

func SetupAppRoutes() *mux.Router {
	router := mux.NewRouter()
	//test prefix
	authRouter := router.PathPrefix("/auth").Subrouter()
	auth.SetupAuthRoutes(authRouter)

	// auth.SetupAuthRoutes(router)
	// normal acima
	return router
}
