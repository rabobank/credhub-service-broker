package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rabobank/credhub-service-broker/conf"
	"github.com/rabobank/credhub-service-broker/credhub"
	"github.com/rabobank/credhub-service-broker/security"
	"github.com/rabobank/credhub-service-broker/server"
	"github.com/rabobank/credhub-service-broker/util"
)

func main() {
	fmt.Printf("credhub-service-broker starting, version:%s, commit:%s\n", conf.VERSION, conf.COMMIT)

	conf.EnvironmentComplete()

	util.ResolveCredhubCredentials()

	security.Initialize()

	initialize()

	server.StartServer()
}

// initialize credhub-service-broker, reading the catalog json file, initializing a cf client, and check for the uaa client.
func initialize() {
	catalogFile := fmt.Sprintf("%s/catalog.json", conf.CatalogDir)
	file, err := os.ReadFile(catalogFile)
	if err != nil {
		fmt.Printf("failed reading catalog file %s: %s\n", catalogFile, err)
		os.Exit(8)
	}
	err = json.Unmarshal(file, &conf.Catalog)
	if err != nil {
		fmt.Printf("failed unmarshalling catalog file %s, error: %s\n", catalogFile, err)
		os.Exit(8)
	}

	util.InitCFClient()

	credhub.Initialize()

}
