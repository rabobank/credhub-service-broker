package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rabobank/credhub-service-broker/credhub"
	"github.com/rabobank/credhub-service-broker/model"
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

func updateMap(originalMap map[string]interface{}, updatedValues map[string]interface{}) {
	for k, v := range updatedValues {
		if originalValue, found := originalMap[k]; !found {
			// the original map doesn't have the key, just set it with whatever value is given
			originalMap[k] = v
		} else if updateSubMap, isMap := v.(map[string]interface{}); !isMap {
			// the updated value for the key is not a map, simply overwrite the value
			originalMap[k] = v
		} else if originalSubMap, isMap := originalValue.(map[string]interface{}); isMap {
			// the updated value is a map and the original value for the key is also a map, let's merge it recursively
			updateMap(originalSubMap, updateSubMap)
		} else {
			// the original value is not a map, simply set it with whatever value is provided
			originalMap[k] = v
		}
	}
}

func ListServiceKeys(w http.ResponseWriter, r *http.Request) {
	if username, serviceInstanceId, ok := userAndService(w, r); ok {
		fmt.Printf("[API] %s updating keys of service %s\n", username, serviceInstanceId)
		if data, e := credhub.GetCredhubData(credentialsPath.ServiceInstanceId(serviceInstanceId)); e != nil {
			fmt.Printf("Unable to get service %s credentials, deeming it a bad request.\n", serviceInstanceId)
			util.WriteHttpResponse(w, http.StatusBadRequest, nil)
		} else if len(data.Data) == 0 {
			fmt.Printf("[API] %s trying to update keys for non-existing credhub service %s\n", username, serviceInstanceId)
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
	if username, serviceInstanceId, ok := userAndService(w, r); ok {
		fmt.Printf("[API] %s updating keys of service %s\n", username, serviceInstanceId)
		if data, e := credhub.GetCredhubData(credentialsPath.ServiceInstanceId(serviceInstanceId)); e != nil {
			fmt.Printf("Unable to get service %s credentials, deeming it a bad request.\n", serviceInstanceId)
			util.WriteHttpResponse(w, http.StatusBadRequest, nil)
		} else if len(data.Data) == 0 {
			fmt.Printf("[API] %s trying to list keys for non-existing credhub service %s\n", username, serviceInstanceId)
			util.WriteHttpResponse(w, http.StatusNotFound, "Not Found")
		} else if credentials, isType := data.Data[0].Value.(map[string]interface{}); isType {
			updatedValues := make(map[string]interface{})
			if e = util.ProvisionObjectFromRequest(r, &updatedValues); e != nil {
				fmt.Printf("Error when processing provided updating object: %v\n", e)
				util.WriteHttpResponse(w, http.StatusBadRequest, "Unable to process json object")
			} else {
				updateMap(credentials, updatedValues)
				if e = credhub.SetCredhubJson(model.CredhubJsonRequest{Type: "json", Name: credentialsPath.ServiceInstanceId(serviceInstanceId), Value: credentials}); e != nil {
					fmt.Printf("Failed to submit credentials update to credhub: %v\n", e)
					util.WriteHttpResponse(w, http.StatusInternalServerError, "Failed to update service")
				} else {
					fmt.Printf("[API] %s has updated credentials in service %s\n", username, serviceInstanceId)
					util.WriteHttpResponse(w, http.StatusAccepted, "Credentials Updated")
				}
			}
		} else {
			fmt.Printf("credentials for service %s do not have a json object map as a value\n", serviceInstanceId)
			util.WriteHttpResponse(w, http.StatusBadRequest, "credentials are not a json object")
		}
	}
}

func DeleteServiceKeys(w http.ResponseWriter, r *http.Request) {
}

func GetServiceHistory(w http.ResponseWriter, r *http.Request) {
}
