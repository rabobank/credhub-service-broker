package util

import (
	"bufio"
	"bytes"
	"crypto/subtle"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/rabobank/credhub-service-broker/conf"
	"github.com/rabobank/credhub-service-broker/model"
)

const (
	cfCertPathEnv        = "CF_INSTANCE_CERT"
	cfKeyPathEnv         = "CF_INSTANCE_KEY"
	TokenRefreshInterval = 90
)

var CfClient cfclient.Client

func Initialize() {
	c := &cfclient.Config{ApiAddress: conf.CfApiURL, ClientID: conf.ClientId, ClientSecret: conf.ClientSecret, SkipSslValidation: conf.SkipSslValidation}
	client, err := cfclient.NewClient(c)
	if err == nil {
		CfClient = *client
		fmt.Printf("cf client initialized with api address %s\n", CfClient.Config.ApiAddress)
	} else {
		fmt.Printf("Unable to authenticate with CF: %s\n", err)
		os.Exit(8)
	}

	// refresh the client every hour to get a new refresh token
	go func() {
		channel := time.Tick(time.Duration(TokenRefreshInterval) * time.Minute)
		for range channel {
			client, err = cfclient.NewClient(c)
			if err != nil {
				fmt.Printf("failed to refresh cfclient, error is %s\n", err)
			}
			fmt.Println("refreshed cf client, got new token")
			CfClient = *client
		}
	}()

}

func WriteHttpResponse(w http.ResponseWriter, code int, object interface{}) {
	data, err := json.Marshal(object)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, err.Error())
		return
	}

	w.WriteHeader(code)
	_, _ = fmt.Fprintf(w, string(data))
	fmt.Printf("response: code:%d, body: %s\n", code, string(data))
}

// BasicAuth validate if user/pass in the http request match the configured service broker user/pass
func BasicAuth(w http.ResponseWriter, r *http.Request, username, password string) bool {
	user, pass, ok := r.BasicAuth()
	if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
		w.Header().Set("WWW-Authenticate", `Basic realm="`+conf.BasicAuthRealm+`"`)
		w.WriteHeader(401)
		_, _ = w.Write([]byte("Unauthorised.\n"))
		return false
	}
	return true
}

func DumpRequest(r *http.Request) {
	if conf.Debug {
		fmt.Printf("dumping %s request for URL: %s\n", r.Method, r.URL)
		fmt.Println("dumping request headers...")
		// Loop over header names
		for name, values := range r.Header {
			if name == "Authorization" {
				fmt.Printf(" %s: %s\n", name, "<redacted>")
			} else {
				// Loop over all values for the name.
				for _, value := range values {
					fmt.Printf(" %s: %s\n", name, value)
				}
			}
		}

		// dump the request body
		fmt.Println("dumping request body...")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Error reading body: %v\n", err)
		} else {
			fmt.Println(string(body))
		}
		// Restore the io.ReadCloser to it's original state
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}
}

func ProvisionObjectFromRequest(r *http.Request, object interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("failed to read json object from request, error: %s\n", err)
		return err
	}
	if conf.Debug {
		fmt.Printf("received body:%v\n", string(body))
	}
	err = json.Unmarshal(body, object)
	if err != nil {
		fmt.Printf("failed to parse json object from request, error: %s\n", err)
		return err
	}
	return nil
}

// ResolveCredhubCredentials - Resolve the credentials by querying credhub for the given paths.
//
//	We implicitly use the app-containers key/cert and use mTLS to get access to the credhub path.
func ResolveCredhubCredentials() {
	// Read the key pair to create certificate
	cert, err := tls.LoadX509KeyPair(os.Getenv(cfCertPathEnv), os.Getenv(cfKeyPathEnv))
	if err != nil {
		log.Fatal("failed to parse the keypair from the app-container", err)
	}
	// Create a HTTPS client and supply the (created CA pool and) certificate
	// client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: caCertPool, Certificates: []tls.Certificate{cert}}}}
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: conf.SkipSslValidation}}}

	// Do the actual mTLS http request
	path := fmt.Sprintf("/api/v1/data?name=%s&current=true", conf.CredhubCredsPath)
	fmt.Printf("trying to get credentials from %s ...\n", conf.CredhubURL+path)
	resp, err := client.Get(conf.CredhubURL + path)
	if err != nil {
		fmt.Printf("Failed to read the credentials from path %s in credhub: %s\n", conf.CredhubCredsPath, err)
		os.Exit(8)
	}
	if resp != nil && resp.StatusCode != http.StatusOK {
		respText, _ := LinesFromReader(resp.Body)
		fmt.Printf("failed to to get credentials from credhub, response code %d, response: %s", resp.StatusCode, *respText)
		os.Exit(8)
	}
	fmt.Println("successfully got the credentials from credhub")

	// Read the response body
	defer func() { _ = resp.Body.Close() }()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("reading response from credhub failed: %s\n", err)
		os.Exit(8)
	}

	// Print the response body to stdout
	// fmt.Printf("response from credhub (DEBUG, REMOVE): \n%s\n", body)

	// parse the response into the model we expect:
	var credhubResponse model.CSBCredentials

	if err = json.Unmarshal(body, &credhubResponse); err != nil {
		fmt.Printf("cannot unmarshal JSON response from %s: %s\n", conf.CredhubURL+path, err)
		os.Exit(8)
	}
	conf.BrokerPassword = credhubResponse.Data[0].Value.CsbBrokerPassword
	conf.ClientSecret = credhubResponse.Data[0].Value.CsbClientSecret

}

func LinesFromReader(r io.Reader) (*[]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &lines, nil
}
