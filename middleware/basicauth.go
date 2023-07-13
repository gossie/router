package middleware

import (
	"net/http"

	"github.com/gossie/router"
)

type UserData struct {
	username, password string
}

func NewUserData(username, password string) *UserData {
	return &UserData{username, password}
}

func BasicAuth(users []*UserData) func(router.HttpHandler) router.HttpHandler {
	return func(in router.HttpHandler) router.HttpHandler {
		return func(w http.ResponseWriter, r *http.Request, pathVariables map[string]string) {
			if user, pass, ok := r.BasicAuth(); ok {
				for _, ud := range users {
					if user == ud.username && pass == ud.password {
						in(w, r, pathVariables)
						return
					}
				}
			}
			http.Error(w, "", http.StatusUnauthorized)
		}
	}
}
