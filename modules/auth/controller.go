package auth

import (
	"net/http"
	functions "oauth2/modules/auth/functions"
)

func Signin(w http.ResponseWriter, r *http.Request) {
	functions.Signin(w, r)
}
func Signup(w http.ResponseWriter, r *http.Request) {
	functions.Signup(w, r)
}
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	functions.RefreshToken(w, r)
}
func RevokeToken(w http.ResponseWriter, r *http.Request) {
	functions.RevokeToken(w, r)
}
