package users

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetupUsersRoutes(router *mux.Router) {

	router.HandleFunc("/", GetUser).Methods(http.MethodGet)

}
