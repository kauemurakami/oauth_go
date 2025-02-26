package routes

import (
	"oauth2/modules/auth"
	"oauth2/modules/users"

	"github.com/gorilla/mux"
)

func SetupAppRoutes() *mux.Router {
	router := mux.NewRouter()
	//test prefix
	authRouter := router.PathPrefix("/auth").Subrouter()
	auth.SetupAuthRoutes(authRouter)
	usersRouter := router.PathPrefix("/users").Subrouter()
	users.SetupUsersRoutes(usersRouter)
	return router
}
