package model

type ServiceInstance struct {
	ServiceId        string                 `json:"service_id"`
	PlanId           string                 `json:"plan_id"`
	OrganizationGuid string                 `json:"organization_guid"`
	SpaceGuid        string                 `json:"space_guid"`
	Context          *Context               `json:"context"`
	Parameters       map[string]interface{} `json:"parameters,omitempty"`
}

type CreateServiceInstanceResponse struct {
	ServiceId     string         `json:"service_id"`
	PlanId        string         `json:"plan_id"`
	DashboardUrl  string         `json:"dashboard_url"`
	LastOperation *LastOperation `json:"last_operation,omitempty"`
}

type DeleteServiceInstanceResponse struct {
	Result string `json:"result,omitempty"`
}
