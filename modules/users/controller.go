package users

import (
	"net/http"
	functions "oauth2/modules/users/functions"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	functions.GetUser(w, r)
}
