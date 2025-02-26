package auth

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetupAuthRoutes(router *mux.Router) {

	// router.HandleFunc("/auth/signin", Signin).Methods(http.MethodPost)
	router.HandleFunc("/signup", Signup).Methods(http.MethodPost)
}
