package security

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudfoundry-community/go-uaa"
	"github.com/gorilla/mux"
	"github.com/rabobank/credhub-service-broker/conf"
	"github.com/rabobank/credhub-service-broker/httpHelper"
	"github.com/rabobank/credhub-service-broker/model"
)

const PermissionsUrl = "%s/v3/service_instances/%s/permissions"

var api *uaa.API

func initializeUaa() error {
	if a, e := uaa.New(conf.UaaApiURL, uaa.WithClientCredentials(conf.ClientId, conf.ClientSecret, uaa.JSONWebToken), uaa.WithSkipSSLValidation(true), uaa.WithVerbosity(false)); e != nil {
		return e
	} else {
		api = a
	}
	return nil
}

func UAA(request *http.Request) (bool, *string) {
	header := request.Header.Get("Authorization")
	if strings.HasPrefix(strings.ToLower(header), "bearer ") {
		token := header[7:]
		// Request carries a bearer token, let's check if it's valid
		body, e := httpHelper.Request(conf.UaaApiURL, "introspect").
			WithHeader("Content-Type", "application/x-www-form-urlencoded").
			PostContentWithClient(api.Client, ([]byte)("token="+token))
		if e != nil {
			fmt.Printf("[AUTH] Error while validating client token: %v\n", e)
			return false, nil
		}
		var tokenIntrospection model.UaaTokenIntrospection
		if e = json.Unmarshal(body, &tokenIntrospection); e != nil {
			fmt.Printf("Failed to introspect token: %v\n", e)
			return false, nil
		}
		if !tokenIntrospection.Active {
			fmt.Printf("[AUTH] User %s trying to use an invalid token\n", tokenIntrospection.UserName)
			return false, nil
		}
		if authentication, isType := request.Context().Value("authentication").(map[string]string); isType {
			authentication["user"] = tokenIntrospection.UserName
		}

		// let's see if there's a service instance id to authroize the user
		if serviceInstanceId, found := mux.Vars(request)["service_instance_guid"]; !found {
			fmt.Printf("[AUTH] Invalid endpoint requested by %s: %s\n", tokenIntrospection.UserName, request.URL.Path)
		} else {
			body, e = httpHelper.
				Request(fmt.Sprintf(PermissionsUrl, conf.CfApiURL, serviceInstanceId)).
				WithBearerToken(token).
				Get()
			if e != nil {
				fmt.Printf("Error while trying to check permissions of user %s for service %s: %v\n", tokenIntrospection.UserName, serviceInstanceId, e)
			} else {
				var permissions model.CfServiceInstancePermissions
				if e = json.Unmarshal(body, &permissions); e != nil {
					fmt.Printf("Unable to check user %s permissions for service %s: %v\n", tokenIntrospection.UserName, serviceInstanceId, e)
				} else if permissions.Manage {
					return true, nil
				} else {
					fmt.Printf("[AUTH] User %s has no permissions to manage requested service %s\n", tokenIntrospection.UserName, serviceInstanceId)
				}
			}
		}
	} else {
		fmt.Println("No bearer token provided")
	}

	return false, nil
}
