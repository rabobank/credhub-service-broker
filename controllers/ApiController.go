package controllers

import (
	"fmt"
	"net/http"
	"strings"

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

func deleteKey(credentials map[string]interface{}, parts []string) bool {
	if value, isFound := credentials[parts[0]]; isFound {
		if len(parts) > 1 {
			// the key has more parts, let's check if it's a map
			if subMap, isMap := value.(map[string]interface{}); isMap {
				if deleteKey(subMap, parts[1:]) {
					if len(subMap) == 0 {
						delete(credentials, parts[0])
					}
					return true
				}
			}
		} else {
			delete(credentials, parts[0])
			return true
		}
	}
	return false
}

func deleteKeys(credentials map[string]interface{}, keysToDelete []string) ([]string, bool) {
	var ignoredKeys []string
	var deletedKeys bool
	for _, k := range keysToDelete {
		keyParts := strings.Split(k, ".")
		if len(keyParts) == 0 {
			ignoredKeys = append(ignoredKeys, k)
		} else if !deleteKey(credentials, keyParts) {
			ignoredKeys = append(ignoredKeys, k)
		} else {
			deletedKeys = true
		}
	}
	return ignoredKeys, deletedKeys
}

func ListServiceKeys(w http.ResponseWriter, r *http.Request) {
	if username, serviceInstanceId, ok := userAndService(w, r); ok {
		fmt.Printf("[API] %s Listing keys of service %s\n", username, serviceInstanceId)
		if data, e := credhub.GetCredhubData(credentialsPath.ServiceInstanceId(serviceInstanceId), 0); e != nil {
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

func ListServiceVersions(w http.ResponseWriter, r *http.Request) {
	if username, serviceInstanceId, ok := userAndService(w, r); ok {
		fmt.Printf("[API] %s listing credential versions for service %s\n", username, serviceInstanceId)
		if data, e := credhub.GetCredhubData(credentialsPath.ServiceInstanceId(serviceInstanceId), 20); e != nil {
			fmt.Printf("Unable to get service %s versions, deeming it a bad request.\n", serviceInstanceId)
			util.WriteHttpResponse(w, http.StatusBadRequest, nil)
		} else if len(data.Data) == 0 {
			fmt.Printf("[API] %s trying to list versions for non-existing credhub service %s\n", username, serviceInstanceId)
			util.WriteHttpResponse(w, http.StatusNotFound, "Not Found")
		} else {
			versions := make([]model.SecretsVersionKeys, 0)
			for _, entry := range data.Data {
				if credentials, isType := entry.Value.(map[string]interface{}); isType {
					versions = append(versions, model.SecretsVersionKeys{
						VersionCreatedAt: entry.VersionCreatedAt,
						ID:               entry.ID,
						Keys:             listMapKeys(credentials, nil, ""),
					})
				} else {
					fmt.Printf("version %s from service %s is not a json map object. Skipping it.\n", entry.ID, serviceInstanceId)
				}
			}

			if len(versions) > 0 {
				util.WriteHttpResponse(w, http.StatusOK, versions)
			} else {
				util.WriteHttpResponse(w, http.StatusNotFound, "credentials don't have any valid versions")
			}
		}
	}
}

func UpdateServiceKeys(w http.ResponseWriter, r *http.Request) {
	if username, serviceInstanceId, ok := userAndService(w, r); ok {
		fmt.Printf("[API] %s updating keys of service %s\n", username, serviceInstanceId)
		if data, e := credhub.GetCredhubData(credentialsPath.ServiceInstanceId(serviceInstanceId), 0); e != nil {
			fmt.Printf("Unable to get service %s credentials, deeming it a bad request.\n", serviceInstanceId)
			util.WriteHttpResponse(w, http.StatusBadRequest, nil)
		} else if len(data.Data) == 0 {
			fmt.Printf("[API] %s trying to update keys for non-existing credhub service %s\n", username, serviceInstanceId)
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
	if username, serviceInstanceId, ok := userAndService(w, r); ok {
		fmt.Printf("[API] %s deleting keys from service %s\n", username, serviceInstanceId)
		if data, e := credhub.GetCredhubData(credentialsPath.ServiceInstanceId(serviceInstanceId), 0); e != nil {
			fmt.Printf("Unable to get service %s credentials, deeming it a bad request.\n", serviceInstanceId)
			util.WriteHttpResponse(w, http.StatusBadRequest, nil)
		} else if len(data.Data) == 0 {
			fmt.Printf("[API] %s trying to list keys for non-existing credhub service %s\n", username, serviceInstanceId)
			util.WriteHttpResponse(w, http.StatusNotFound, "Not Found")
		} else if credentials, isType := data.Data[0].Value.(map[string]interface{}); isType {
			keysToDelete := make([]string, 0)
			if e = util.ProvisionObjectFromRequest(r, &keysToDelete); e != nil {
				fmt.Printf("Error when processing provided deleting array: %v\n", e)
				util.WriteHttpResponse(w, http.StatusBadRequest, "Unable to process json array")
			} else {
				ignoredKeys, keysHaveBeenDeleted := deleteKeys(credentials, keysToDelete)
				response := model.DeleteResponse{IgnoredKeys: ignoredKeys}

				if keysHaveBeenDeleted {
					if e = credhub.SetCredhubJson(model.CredhubJsonRequest{Type: "json", Name: credentialsPath.ServiceInstanceId(serviceInstanceId), Value: credentials}); e != nil {
						fmt.Printf("Failed to submit credentials update to credhub: %v\n", e)
						util.WriteHttpResponse(w, http.StatusInternalServerError, "Failed to update service")
					} else {
						fmt.Printf("[API] %s has deleted keys from service %s credentials\n", username, serviceInstanceId)
						util.WriteHttpResponse(w, http.StatusAccepted, response)
					}
				} else {
					util.WriteHttpResponse(w, http.StatusNotModified, response)
				}
			}
		} else {
			fmt.Printf("credentials for service %s do not have a json object map as a value\n", serviceInstanceId)
			util.WriteHttpResponse(w, http.StatusBadRequest, "credentials are not a json object")
		}
	}
}

func ReinstateServiceVersion(w http.ResponseWriter, r *http.Request) {
	if username, serviceInstanceId, ok := userAndService(w, r); ok {
		versionId := mux.Vars(r)["version_id"]
		fmt.Printf("[API] %s reinstating credential version %s for service %s\n", username, versionId, serviceInstanceId)

		if data, e := credhub.GetCredhubDataVersion(versionId); e != nil {
			fmt.Printf("Unable to get version %s, deeming it a bad request: %v\n", versionId, e)
			util.WriteHttpResponse(w, http.StatusBadRequest, nil)
		} else if !strings.HasPrefix(data.Name, credentialsPath.ServiceInstanceId(serviceInstanceId)) {
			fmt.Printf("[API] %s trying to reinstate a version to service %s from another service.\n", username, serviceInstanceId)
			util.WriteHttpResponse(w, http.StatusBadRequest, nil)
		} else if currentData, e := credhub.GetCredhubData(credentialsPath.ServiceInstanceId(serviceInstanceId), 0); e != nil {
			fmt.Printf("Unable to get service %s credentials, deeming it a bad request.\n", serviceInstanceId)
			util.WriteHttpResponse(w, http.StatusBadRequest, nil)
		} else if len(currentData.Data) == 0 {
			fmt.Printf("[API] %s trying to reinstate a version for a non-existing credhub service %s\n", username, serviceInstanceId)
			util.WriteHttpResponse(w, http.StatusNotFound, "Not Found")
		} else if currentData.Data[0].ID == versionId {
			fmt.Printf("Credentials version %s being reinstated for service %s is already the current one", versionId, serviceInstanceId)
			util.WriteHttpResponse(w, http.StatusNotModified, "Unmodified")
		} else if e = credhub.SetCredhubJson(model.CredhubJsonRequest{Type: "json", Name: credentialsPath.ServiceInstanceId(serviceInstanceId), Value: data.Value.(map[string]interface{})}); e != nil {
			fmt.Printf("Failed to submit credentials update to credhub: %v\n", e)
			util.WriteHttpResponse(w, http.StatusInternalServerError, "Failed to update service")
		} else {
			fmt.Printf("[API] %s has updated credentials in service %s\n", username, serviceInstanceId)
			util.WriteHttpResponse(w, http.StatusAccepted, "Credentials Updated")
		}
	}
}
