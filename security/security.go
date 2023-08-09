package security

import (
	"context"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rabobank/credhub-service-broker/conf"
)

func Initialize() {
	basicChallenge = `Basic realm="` + conf.BasicAuthRealm + `"`
	basicCredentials.username = conf.BrokerUser
	basicCredentials.password = conf.BrokerPassword

	initializeUaa()
}

func MatchPrefix(pathPrefix string) MatchBuilder {
	builder := &middlewareBuilder{}
	return builder.MatchPrefix(pathPrefix)
}

// Authenticator an authenticator function signature, it should return if a request is authenticated and if not, which challenges, if any, to return to the client
type Authenticator func(*http.Request) (authenticated bool, challenges *string)

// Anonymous default authenticator for anonymous access
func Anonymous(*http.Request) (bool, *string) {
	return true, nil
}

// authenticatorMatchers holds the list of matching strings that an authenticator should be applied to.
type authenticatorMatchers struct {
	matchers     []string
	authenticate Authenticator
}

func Middleware(authenticators []authenticatorMatchers) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			authorized := false
			var challenge *string

		check:
			for _, authenticator := range authenticators {
				for _, prefix := range authenticator.matchers {
					if strings.HasPrefix(request.URL.Path, prefix) {
						request = request.WithContext(context.WithValue(request.Context(), "authentication", make(map[string]interface{})))
						authorized, challenge = authenticator.authenticate(request)
						break check
					}
				}
			}

			if !authorized {
				if challenge != nil {
					writer.Header().Set("WWW-Authenticate", *challenge)
				}
				writer.WriteHeader(401)
				_, _ = writer.Write([]byte("Unauthorised.\n"))
			} else {
				next.ServeHTTP(writer, request)
			}
		})
	}
}
