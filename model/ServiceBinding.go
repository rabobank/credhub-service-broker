package model

type ServiceBinding struct {
	ServiceId         string        `json:"service_id"`
	PlanId            string        `json:"plan_id"`
	AppGuid           string        `json:"app_guid"`
	ServiceInstanceId string        `json:"service_instance_id"`
	BindResource      *BindResource `json:"bind_resource"`
	Context           *Context      `json:"context"`
}

type BindResource struct {
	AppGuid   string `json:"app_guid"`
	SpaceGuid string `json:"space_guid"`
}

type CreateServiceBindingResponse struct {
	// SyslogDrainUrl string      `json:"syslog_drain_url, omitempty"`
	Credentials *Credentials `json:"credentials"`
}

type DeleteServiceBindingResponse struct {
	// SyslogDrainUrl string      `json:"syslog_drain_url, omitempty"`
	Result string `json:"result"`
}

type Credentials struct {
	CredhubRef string `json:"credhub-ref"`
}
