package users

import (
	"net/http"
	middlewares "oauth2/core/middlewares"

	"github.com/gorilla/mux"
)

func SetupUsersRoutes(router *mux.Router) {
	router.Use(middlewares.AuthMiddleware)
	router.HandleFunc("/", GetUser).Methods(http.MethodGet)

}
