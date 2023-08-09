package credhub

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cloudfoundry-community/go-uaa"
	"github.com/rabobank/credhub-service-broker/conf"
	"github.com/rabobank/credhub-service-broker/model"
	"github.com/rabobank/credhub-service-broker/util"
	"golang.org/x/oauth2"
)

var (
	UaaUrl   string
	uaaToken *oauth2.Token
	uaac     *uaa.API
)

func Initialize() {
	UaaUrl = util.CfClient.Endpoint.TokenEndpoint
	if client, e := uaa.New(UaaUrl, uaa.WithClientCredentials(conf.ClientId, conf.ClientSecret, uaa.JSONWebToken), uaa.WithSkipSSLValidation(conf.SkipSslValidation)); e == nil {
		uaac = client
	} else {
		fmt.Printf("Failed to authenticate with UAA: %v\n", e)
		os.Exit(8)
	}
}

func SetCredhubData(credhubData model.CredhubDataRequest) error {
	jsonBytes, err := json.Marshal(credhubData)
	if err != nil {
		return err
	}
	return setCredhub(jsonBytes)
}

func SetCredhubJson(credhubJson model.CredhubJsonRequest) error {
	jsonBytes, err := json.Marshal(credhubJson)
	if err != nil {
		return err
	}
	return setCredhub(jsonBytes)
}

func setCredhub(jsonBytes []byte) error {
	var resp *http.Response
	var credhubDataResponse model.CredhubDataResponse
	path := "/api/v1/data"
	client, req, err := getHttpClientAndRequest(http.MethodPut, path, jsonBytes)
	if err == nil {
		resp, err = client.Do(req)
		if err == nil && resp != nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			if err = json.Unmarshal(body, &credhubDataResponse); err != nil {
				return errors.New(fmt.Sprintf("cannot unmarshal JSON response from %s: %s\n", conf.CredhubURL+path, err))
			}
		} else {
			// response code != http.StatusOK, so handle that error as well
			if err == nil && resp != nil {
				respText, _ := util.LinesFromReader(resp.Body)
				return errors.New(fmt.Sprintf("cannot set credhub data, response: %s, body: %v", resp.Status, *respText))
			}
		}
	}
	return err
}

func GetCredhubData(credhubPath string) (model.CredhubEntry, error) {
	var err error
	var credhubEntry model.CredhubEntry
	var resp *http.Response
	path := fmt.Sprintf("/api/v1/data?name=%s&current=true", credhubPath)
	client, req, err := getHttpClientAndRequest(http.MethodGet, path, nil)
	if err == nil {
		resp, err = client.Do(req)
		if err == nil && resp != nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			if err = json.Unmarshal(body, &credhubEntry); err != nil {
				return credhubEntry, errors.New(fmt.Sprintf("cannot unmarshal JSON response from %s: %s\n", conf.CredhubURL+path, err))
			}
		}
	}
	return credhubEntry, err
}

func DeleteCredhubData(credhubPath string) error {
	var err error
	var resp *http.Response
	client, req, err := getHttpClientAndRequest(http.MethodDelete, fmt.Sprintf("/api/v1/data?name=%s", credhubPath), nil)
	if err == nil {
		resp, err = client.Do(req)
		if err == nil && resp != nil && resp.StatusCode == http.StatusNoContent {
			return nil
		}
	}
	return err
}

func CreateCredhubPermission(credhubPermission model.CredhubPermissionRequest) error {
	var err error
	var jsonBytes []byte
	var credhubPermissionResponse model.CredhubPermissionResponse
	jsonBytes, err = json.Marshal(credhubPermission)
	if err != nil {
		return err
	}
	var resp *http.Response
	path := "/api/v2/permissions"
	client, req, err := getHttpClientAndRequest(http.MethodPost, path, jsonBytes)
	if err == nil {
		resp, err = client.Do(req)
		if err == nil && resp != nil && resp.StatusCode == http.StatusCreated {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			if err = json.Unmarshal(body, &credhubPermissionResponse); err != nil {
				return errors.New(fmt.Sprintf("Can not unmarshal JSON response from %s: %s\n", conf.CredhubURL+path, err))
			}
		} else {
			// response code != http.StatusCreated, so handle that error as well
			if err == nil && resp != nil {
				respText, _ := util.LinesFromReader(resp.Body)
				for _, line := range *respText {
					if strings.Contains(line, "A permission entry for this actor and path already exists") {
						fmt.Printf("%s\n", *respText)
						// ignore the error
						return err
					}
				}
				return errors.New(fmt.Sprintf("response %s, body: %v", resp.Status, *respText))
			}
		}
	}
	return err
}

func GetCredhubPermission(credhubPath, actor string) (model.CredhubPermissionResponse, error) {
	var err error
	var resp *http.Response
	var credhubPermissionResponse model.CredhubPermissionResponse
	path := fmt.Sprintf("/api/v2/permissions?path=%s&actor=%s", credhubPath, actor)
	client, req, err := getHttpClientAndRequest(http.MethodGet, path, nil)
	if err == nil {
		resp, err = client.Do(req)
		if err == nil && resp != nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			if err = json.Unmarshal(body, &credhubPermissionResponse); err != nil {
				return credhubPermissionResponse, errors.New(fmt.Sprintf("Can not unmarshal JSON response from %s: %s\n", conf.CredhubURL+path, err))
			}
		}
	}
	return credhubPermissionResponse, err
}

func DeleteCredhubPermission(uuid string) error {
	var err error
	var resp *http.Response
	client, req, err := getHttpClientAndRequest(http.MethodDelete, fmt.Sprintf("/api/v2/permissions/%s", uuid), nil)
	if err == nil {
		resp, err = client.Do(req)
		if err == nil && resp != nil && resp.StatusCode == http.StatusOK {
			return nil
		}
	}
	return err
}

func getHttpClientAndRequest(method, path string, postData []byte) (*http.Client, *http.Request, error) {
	var err error
	var client http.Client
	var req *http.Request
	if uaaToken == nil || !uaaToken.Valid() {
		uaaToken, err = Token()
	}

	if err == nil {
		transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		client = http.Client{Timeout: time.Duration(conf.HttpTimeout) * time.Second, Transport: transport}
		req, _ = http.NewRequest(method, conf.CredhubURL+path, bytes.NewReader(postData))
		uaaToken.SetAuthHeader(req)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-type", "application/json")
	}
	return &client, req, err
}

func Token() (*oauth2.Token, error) {
	return uaac.Token(context.TODO())
}
