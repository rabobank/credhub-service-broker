package model

type CfApiEndpoints struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		CloudControllerV2 struct {
			Href string `json:"href"`
			Meta struct {
				Version string `json:"version"`
			} `json:"meta"`
		} `json:"cloud_controller_v2"`
		CloudControllerV3 struct {
			Href string `json:"href"`
			Meta struct {
				Version string `json:"version"`
			} `json:"meta"`
		} `json:"cloud_controller_v3"`
		NetworkPolicyV0 struct {
			Href string `json:"href"`
		} `json:"network_policy_v0"`
		NetworkPolicyV1 struct {
			Href string `json:"href"`
		} `json:"network_policy_v1"`
		Login struct {
			Href string `json:"href"`
		} `json:"login"`
		Uaa struct {
			Href string `json:"href"`
		} `json:"uaa"`
		Credhub struct {
			Href string `json:"href"`
		} `json:"credhub"`
		Routing struct {
			Href string `json:"href"`
		} `json:"routing"`
		Logging struct {
			Href string `json:"href"`
		} `json:"logging"`
		LogCache struct {
			Href string `json:"href"`
		} `json:"log_cache"`
		LogStream struct {
			Href string `json:"href"`
		} `json:"log_stream"`
		AppSsh struct {
			Href string `json:"href"`
			Meta struct {
				HostKeyFingerprint string `json:"host_key_fingerprint"`
				OauthClient        string `json:"oauth_client"`
			} `json:"meta"`
		} `json:"app_ssh"`
	} `json:"links"`
}

type CfServiceInstancePermissions struct {
	Read   bool `json:"read"`
	Manage bool `json:"manage"`
}
