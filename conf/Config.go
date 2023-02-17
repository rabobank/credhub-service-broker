package conf

import (
	"fmt"
	"github.com/rabobank/credhub-service-broker/model"
	"os"
	"strconv"
	"strings"
)

var (
	//  PCSB  :  Panzer Credhub Service Broker
	debugStr           = os.Getenv("CSB_DEBUG")
	Debug              = false
	httpTimeoutStr     = os.Getenv("CSB_HTTP_TIMEOUT")
	HttpTimeout        int
	HttpTimeoutDefault = 10
	ClientId           = os.Getenv("CSB_CLIENT_ID")
	ClientSecret       string // will be resolved from config in credhub path CredhubCredsPath
	CredhubURL         = os.Getenv("CSB_CREDHUB_URL")

	Catalog    model.Catalog
	ListenPort int

	BrokerUser           = os.Getenv("CSB_BROKER_USER")
	BrokerPassword       string // will be resolved from config in credhub path CredhubCredsPath
	CatalogDir           = os.Getenv("CSB_CATALOG_DIR")
	ListenPortStr        = os.Getenv("CSB_LISTEN_PORT")
	CfApiURL             = os.Getenv("CSB_CFAPI_URL")
	SkipSslValidationStr = os.Getenv("CSB_SKIP_SSL_VALIDATION")
	SkipSslValidation    bool
	CredhubCredsPath     = os.Getenv("CREDHUB_CREDS_PATH") // something like /credhub-service-broker/config
)

const BasicAuthRealm = "PCSB Panzer Credhub Service Broker"

// EnvironmentComplete - Check for required environment variables and exit if not all are there.
func EnvironmentComplete() {
	envComplete := true
	if debugStr == "true" {
		Debug = true
	}
	if httpTimeoutStr == "" {
		HttpTimeout = HttpTimeoutDefault
	} else {
		var err error
		HttpTimeout, err = strconv.Atoi(httpTimeoutStr)
		if err != nil {
			fmt.Printf("failed reading envvar CSB_HTTP_TIMEOUT, err: %s\n", err)
			envComplete = false
		}
	}
	if CredhubURL == "" {
		CredhubURL = "https://credhub.service.cf.internal:8844"
	}
	if ClientId == "" {
		envComplete = false
		fmt.Println("missing envvar: CSB_CLIENT_ID")
	}
	if BrokerUser == "" {
		envComplete = false
		fmt.Println("missing envvar: CSB_BROKER_USER")
	}
	if CatalogDir == "" {
		CatalogDir = "./catalog"
	}
	if ListenPortStr == "" {
		ListenPort = 8080
	} else {
		var err error
		ListenPort, err = strconv.Atoi(ListenPortStr)
		if err != nil {
			fmt.Printf("failed reading envvar LISTEN_PORT, err: %s\n", err)
			envComplete = false
		}
	}
	if CredhubCredsPath == "" {
		envComplete = false
		fmt.Println("missing envvar: CREDHUB_CREDS_PATH")
	}
	if CfApiURL == "" {
		envComplete = false
		fmt.Println("missing envvar: CSB_CFAPI_URL")
	}

	if strings.EqualFold(SkipSslValidationStr, "true") {
		SkipSslValidation = true
	}

	if !envComplete {
		fmt.Println("one or more required environment variables missing, aborting...")
		os.Exit(8)
	}
}
