package model

import "time"

type CredhubEntry struct {
	Data []CredhubDataResponse `json:"data"`
}

type CredhubDataResponse struct {
	Type             string      `json:"type"`
	VersionCreatedAt time.Time   `json:"version_created_at"`
	ID               string      `json:"id"`
	Name             string      `json:"name"`
	Value            interface{} `json:"value"`
}

type CredhubDataRequest struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CredhubJsonRequest struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value JSON   `json:"value"`
}

type JSON map[string]interface{}

type CredhubInfoResponse struct {
	AuthServer struct {
		URL string `json:"url"`
	} `json:"auth-server"`
	App struct {
		Name string `json:"name"`
	} `json:"app"`
}

type CredhubPermissionRequest struct {
	Path       string   `json:"path"`
	Actor      string   `json:"actor"`
	Operations []string `json:"operations"`
}

type CredhubPermissionResponse struct {
	Path       string   `json:"path"`
	Operations []string `json:"operations"`
	Actor      string   `json:"actor"`
	UUID       string   `json:"uuid"`
}
type UAATokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	Jti         string `json:"jti"`
}

// CSBCredentials - The structure that is returned when querying credhub for the csb credentials (containing the broker password and the csb cf client secret)
type CSBCredentials struct {
	Data []struct {
		Type             string      `json:"type"`
		VersionCreatedAt time.Time   `json:"version_created_at"`
		ID               string      `json:"id"`
		Name             string      `json:"name"`
		Metadata         interface{} `json:"metadata"`
		Value            struct {
			CsbBrokerPassword string `json:"CSB_BROKER_PASSWORD"`
			CsbClientSecret   string `json:"CSB_CLIENT_SECRET"`
		} `json:"value"`
	} `json:"data"`
}
