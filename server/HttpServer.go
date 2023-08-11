package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rabobank/credhub-service-broker/conf"
	"github.com/rabobank/credhub-service-broker/controllers"
	"github.com/rabobank/credhub-service-broker/security"
	"github.com/rabobank/credhub-service-broker/util"
)

var HealthStatus = struct{ Status string }{"UP"}

func Health(w http.ResponseWriter, _ *http.Request) {
	util.WriteHttpResponse(w, http.StatusOK, HealthStatus)
}

func StartServer() {
	router := mux.NewRouter()

	router.Use(controllers.DebugMiddleware)
	// oauth2 interceptor. It will only handle /api endpoints
	router.Use(security.MatchPrefix("/health").AuthenticateWith(security.Anonymous).
		MatchPrefix("/api").AuthenticateWith(security.UAA).
		Default(security.BasicAuth).Build())
	router.Use(controllers.AuditLogMiddleware)

	// service broker endpoints
	router.HandleFunc("/v2/catalog", controllers.Catalog).Methods("GET")
	router.HandleFunc("/v2/service_instances/{service_instance_guid}", controllers.CreateOrUpdateServiceInstance).Methods("PUT", "PATCH")
	router.HandleFunc("/v2/service_instances/{service_instance_guid}", controllers.DeleteServiceInstance).Methods("DELETE")
	router.HandleFunc("/v2/service_instances/{service_instance_guid}/service_bindings/{service_binding_guid}", controllers.CreateServiceBinding).Methods("PUT")
	router.HandleFunc("/v2/service_instances/{service_instance_guid}/service_bindings/{service_binding_guid}", controllers.DeleteServiceBinding).Methods("DELETE")

	// health endpoint
	router.HandleFunc("/health", Health).Methods("GET")

	// key management api endpoints
	router.HandleFunc("/api/{service_instance_guid}/keys", controllers.ListServiceKeys).Methods("GET")
	router.HandleFunc("/api/{service_instance_guid}/keys", controllers.UpdateServiceKeys).Methods("PUT")
	router.HandleFunc("/api/{service_instance_guid}/keys", controllers.DeleteServiceKeys).Methods("DELETE")
	router.HandleFunc("/api/{service_instance_guid}/versions", controllers.ListServiceVersions).Methods("GET")

	http.Handle("/", router)

	router.Use(controllers.AddHeadersMiddleware)

	fmt.Printf("server started, listening on port %d...\n", conf.ListenPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", conf.ListenPort), nil)
	if err != nil {
		fmt.Printf("failed to start http server on port %d, err: %s\n", conf.ListenPort, err)
		os.Exit(8)
	}
}
