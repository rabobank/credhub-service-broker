package security

import (
	"crypto/subtle"
	"net/http"
)

var basicCredentials struct {
	username string
	password string
}
var basicChallenge string

func BasicAuth(r *http.Request) (bool, *string) {

	user, pass, ok := r.BasicAuth()
	if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(basicCredentials.username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(basicCredentials.password)) != 1 {
		return false, &basicChallenge
	}
	return true, nil
}
