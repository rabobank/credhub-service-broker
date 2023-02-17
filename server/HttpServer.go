package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rabobank/credhub-service-broker/conf"
	"github.com/rabobank/credhub-service-broker/controllers"
	"net/http"
	"os"
)

func StartServer() {
	router := mux.NewRouter()

	router.Use(controllers.DebugMiddleware)
	router.Use(controllers.BasicAuthMiddleware)
	router.Use(controllers.AuditLogMiddleware)

	router.HandleFunc("/v2/catalog", controllers.Catalog).Methods("GET")
	router.HandleFunc("/v2/service_instances/{service_instance_guid}", controllers.CreateOrUpdateServiceInstance).Methods("PUT", "PATCH")
	router.HandleFunc("/v2/service_instances/{service_instance_guid}", controllers.DeleteServiceInstance).Methods("DELETE")
	router.HandleFunc("/v2/service_instances/{service_instance_guid}/service_bindings/{service_binding_guid}", controllers.CreateServiceBinding).Methods("PUT")
	router.HandleFunc("/v2/service_instances/{service_instance_guid}/service_bindings/{service_binding_guid}", controllers.DeleteServiceBinding).Methods("DELETE")

	http.Handle("/", router)

	router.Use(controllers.AddHeadersMiddleware)

	fmt.Printf("server started, listening on port %d...\n", conf.ListenPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", conf.ListenPort), nil)
	if err != nil {
		fmt.Printf("failed to start http server on port %d, err: %s\n", conf.ListenPort, err)
		os.Exit(8)
	}
}
