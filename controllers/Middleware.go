package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/rabobank/credhub-service-broker/conf"
	"github.com/rabobank/credhub-service-broker/util"
	"net/http"
	"strings"
)

const IdentityHeader = "X-Broker-Api-Originating-Identity"

func BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if util.BasicAuth(w, r, conf.BrokerUser, conf.BrokerPassword) {
			// Call the next handler, which can be another middleware in the chain, or the final handler.
			next.ServeHTTP(w, r)
		}
	})
}

func DebugMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		util.DumpRequest(r)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func AddHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// AuditLogMiddleware - We are looking for the X-Broker-Api-Request-Identity header, see https://github.com/openservicebrokerapi/servicebroker/blob/v2.16/spec.md#originating-identity
func AuditLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origIdentity := "UNKNOWN"
		var jsonObject = OrigIdentity{}
		if identityHeaders := r.Header[IdentityHeader]; identityHeaders != nil && len(identityHeaders) > 0 {
			identityHeader := identityHeaders[0]
			if words := strings.Split(identityHeader, " "); len(words) == 2 {
				if decodedString, err := base64.StdEncoding.DecodeString(words[1]); decodedString != nil && err == nil {
					if err = json.Unmarshal(decodedString, &jsonObject); err == nil {
						if cfUser, err := util.CfClient.GetUserByGUID(jsonObject.UserID); err == nil {
							origIdentity = cfUser.Username
						} else {
							fmt.Printf("failed to cf lookup user with guid %s: %s\n", jsonObject.UserID, err)
						}
					} else {
						fmt.Printf("failed to parse %s header: %s\n", IdentityHeader, err)
					}
				} else {
					fmt.Printf("failed to base64 decode %s header: %s\n", IdentityHeader, err)
				}
			}
		}

		fmt.Printf("%s request on path %s by user %s (guid:%s)\n", r.Method, r.RequestURI, origIdentity, jsonObject.UserID)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

type OrigIdentity struct {
	UserID string `json:"user_id"`
}
