package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rabobank/credhub-service-broker/credhub"
	"github.com/rabobank/credhub-service-broker/model"
	"github.com/rabobank/credhub-service-broker/util"
	"net/http"
)

func CreateServiceBinding(w http.ResponseWriter, r *http.Request) {
	var err error
	serviceInstanceId := mux.Vars(r)["service_instance_guid"]
	serviceBindingId := mux.Vars(r)["service_binding_guid"]
	var serviceBinding model.ServiceBinding
	err = util.ProvisionObjectFromRequest(r, &serviceBinding)
	if err != nil {
		util.WriteHttpResponse(w, http.StatusBadRequest, model.BrokerError{Error: "FAILED", Description: err.Error(), InstanceUsable: false, UpdateRepeatable: false})
		return
	}
	credhubPath := fmt.Sprintf("/pcsb/%s/credentials", serviceInstanceId)
	fmt.Printf("create service binding id %s for service instance id %s, creating path %s...\n", serviceBindingId, serviceInstanceId, credhubPath)
	err = credhub.CreateCredhubPermission(model.CredhubPermissionRequest{Path: credhubPath, Actor: fmt.Sprintf("mtls-app:%s", serviceBinding.AppGuid), Operations: []string{"read"}})
	if err != nil {
		util.WriteHttpResponse(w, http.StatusBadRequest, model.BrokerError{Error: "FAILED", Description: err.Error(), InstanceUsable: false, UpdateRepeatable: false})
		return
	}
	// now store the app_guid in credhub as well, in a entry keyed with the binding_guid, since we need the app_guid when we do an unbind for the service
	err = credhub.SetCredhubData(model.CredhubDataRequest{Type: "value", Name: fmt.Sprintf("/pcsb/%s/%s", serviceInstanceId, serviceBindingId), Value: serviceBinding.AppGuid})
	if err != nil {
		util.WriteHttpResponse(w, http.StatusBadRequest, model.BrokerError{Error: "FAILED", Description: err.Error(), InstanceUsable: false, UpdateRepeatable: false})
		return
	}
	response := model.CreateServiceBindingResponse{Credentials: &model.Credentials{CredhubRef: credhubPath}}
	util.WriteHttpResponse(w, http.StatusCreated, response)
}

func DeleteServiceBinding(w http.ResponseWriter, r *http.Request) {
	serviceInstanceId := mux.Vars(r)["service_instance_guid"]
	serviceBindingId := mux.Vars(r)["service_binding_guid"]
	fmt.Printf("delete service binding id %s for service instance id %s...\n", serviceBindingId, serviceInstanceId)
	credhubPath := fmt.Sprintf("/pcsb/%s/%s", serviceInstanceId, serviceBindingId)
	credhubEntry, err := credhub.GetCredhubData(credhubPath)
	if err != nil {
		util.WriteHttpResponse(w, http.StatusOK, model.BrokerError{Error: "FAILED", Description: fmt.Sprintf("failed to read credhub at path %s: %s", credhubPath, err.Error()), InstanceUsable: false, UpdateRepeatable: false})
		return
	}
	if len(credhubEntry.Data) == 0 {
		util.WriteHttpResponse(w, http.StatusOK, model.BrokerError{Error: "FAILED", Description: fmt.Sprintf("credhub entry %s was not found or was not readable", credhubPath), InstanceUsable: false, UpdateRepeatable: false})
		return
	}
	path := fmt.Sprintf("/pcsb/%s/credentials", serviceInstanceId)
	actor := fmt.Sprintf("mtls-app:%s", credhubEntry.Data[0].Value)
	credhubPerm, err := credhub.GetCredhubPermission(path, actor)
	err = credhub.DeleteCredhubPermission(credhubPerm.UUID)
	if err != nil {
		util.WriteHttpResponse(w, http.StatusBadRequest, model.BrokerError{Error: "FAILED", Description: err.Error(), InstanceUsable: false, UpdateRepeatable: false})
		return
	}
	err = credhub.DeleteCredhubData(fmt.Sprintf("/pcsb/%s/%s", serviceInstanceId, serviceBindingId))
	if err != nil {
		util.WriteHttpResponse(w, http.StatusBadRequest, model.BrokerError{Error: "FAILED", Description: err.Error(), InstanceUsable: false, UpdateRepeatable: false})
		return
	}
	util.WriteHttpResponse(w, http.StatusOK, model.DeleteServiceBindingResponse{Result: "unbind completed"})
}
