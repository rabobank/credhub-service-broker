package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rabobank/credhub-service-broker/credhub"
	"github.com/rabobank/credhub-service-broker/util"
)

func userAndService(w http.ResponseWriter, r *http.Request) (username string, serviceInstanceId string, ok bool) {
	ok = false
	authenticationContext, isType := r.Context().Value("authentication").(map[string]interface{})
	if !isType {
		fmt.Println("API reached without an authentication context. Failing")
		util.WriteHttpResponse(w, http.StatusInternalServerError, "Authorization Failure")
		return
	}
	entry, found := authenticationContext["user"]
	if !found {
		fmt.Println("API reached without an authentication context user. Failing")
		util.WriteHttpResponse(w, http.StatusInternalServerError, "Authorization Failure")
		return
	}
	username = entry.(string)
	serviceInstanceId = mux.Vars(r)["service_instance_guid"]
	ok = true
	return
}

func listMapKeys(m map[string]interface{}, keys []string, prefix string) []string {
	for k, v := range m {
		if subMap, isType := v.(map[string]interface{}); isType {
			keys = listMapKeys(subMap, keys, prefix+k+".")
		} else {
			keys = append(keys, prefix+k)
		}
	}
	return keys
}

func ListServiceKeys(w http.ResponseWriter, r *http.Request) {
	if username, serviceInstanceId, ok := userAndService(w, r); ok {
		fmt.Printf("[API] %s listing keys of service %s\n", username, serviceInstanceId)
		if data, e := credhub.GetCredhubData(credentialsPath.ServiceInstanceId(serviceInstanceId)); e != nil {
			fmt.Printf("Unable to get service %s credentials, deeming it a bad request.\n", serviceInstanceId)
			util.WriteHttpResponse(w, http.StatusBadRequest, nil)
		} else if len(data.Data) == 0 {
			fmt.Printf("[API] %s trying to list keys for non-existing credhub service %s\n", username, serviceInstanceId)
			util.WriteHttpResponse(w, http.StatusNotFound, "Not Found")
		} else if credentials, isType := data.Data[0].Value.(map[string]interface{}); isType {
			keys := listMapKeys(credentials, nil, "")
			util.WriteHttpResponse(w, http.StatusOK, keys)
		} else {
			fmt.Printf("credentials for service %s do not have a json object map as a value\n", serviceInstanceId)
			util.WriteHttpResponse(w, http.StatusBadRequest, "credentials are not a json object")
		}
	}
}

func UpdateServiceKeys(w http.ResponseWriter, r *http.Request) {
}

func DeleteServiceKeys(w http.ResponseWriter, r *http.Request) {
}

func GetServiceHistory(w http.ResponseWriter, r *http.Request) {
}
