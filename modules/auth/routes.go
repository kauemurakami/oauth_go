package auth

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetupAuthRoutes(router *mux.Router) {

	router.HandleFunc("/signin", Signin).Methods(http.MethodPost)
	router.HandleFunc("/signup", Signup).Methods(http.MethodPost)
	router.HandleFunc("/refresh-token", RefreshToken).Methods(http.MethodPost)
	router.HandleFunc("/revoke-token", RevokeToken).Methods(http.MethodPost)
}
