package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rabobank/credhub-service-broker/conf"
	"github.com/rabobank/credhub-service-broker/credhub"
	"github.com/rabobank/credhub-service-broker/model"
	"github.com/rabobank/credhub-service-broker/util"
)

func Catalog(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("get service broker catalog from %s...\n", r.RemoteAddr)
	util.WriteHttpResponse(w, http.StatusOK, conf.Catalog)
}

func CreateOrUpdateServiceInstance(w http.ResponseWriter, r *http.Request) {
	var err error
	serviceInstanceId := mux.Vars(r)["service_instance_guid"]
	fmt.Printf("create/update service instance for %s...\n", serviceInstanceId)
	var serviceInstance model.ServiceInstance
	err = util.ProvisionObjectFromRequest(r, &serviceInstance)
	if err != nil {
		util.WriteHttpResponse(w, http.StatusBadRequest, model.BrokerError{Error: "FAILED", Description: err.Error(), InstanceUsable: false, UpdateRepeatable: false})
		return
	}
	// check max len of creds (64 KB //TODO)

	var jsonValues = model.JSON{}
	for key, value := range serviceInstance.Parameters {
		jsonValues[key] = value
	}
	credhubPath := fmt.Sprintf("/pcsb/%s/credentials", serviceInstanceId)
	err = credhub.SetCredhubJson(model.CredhubJsonRequest{Type: "json", Name: credhubPath, Value: jsonValues})
	if err != nil {
		util.WriteHttpResponse(w, http.StatusBadRequest, model.BrokerError{Error: "FAILED", Description: fmt.Sprintf("failed to set JSON value for path %s: %s", credhubPath, err.Error()), InstanceUsable: false, UpdateRepeatable: false})
		return
	}
	fmt.Printf("created credhub path %s\n", credhubPath)
	util.WriteHttpResponse(w, http.StatusOK, model.CreateServiceInstanceResponse{})
	return
}

func DeleteServiceInstance(w http.ResponseWriter, r *http.Request) {
	serviceInstanceId := mux.Vars(r)["service_instance_guid"]
	fmt.Printf("delete service instance %s...\n", serviceInstanceId)
	_, err := credhub.GetCredhubData(fmt.Sprintf("/pcsb/%s/credentials", serviceInstanceId), 0)
	if err != nil {
		util.WriteHttpResponse(w, http.StatusGone, err)
		return
	}

	credhubPath := fmt.Sprintf("pcsb/%s/credentials", serviceInstanceId)
	err = credhub.DeleteCredhubData(credhubPath)
	if err != nil {
		response := model.DeleteServiceInstanceResponse{Result: fmt.Sprintf("failed to delete path %s, error: %s", credhubPath, err)}
		util.WriteHttpResponse(w, http.StatusBadRequest, response)
		return
	}
	fmt.Printf("deleted credhub path %s\n", credhubPath)
	response := model.DeleteServiceInstanceResponse{Result: fmt.Sprintf("Service instance %s deleted", serviceInstanceId)}
	util.WriteHttpResponse(w, http.StatusOK, response)
}
