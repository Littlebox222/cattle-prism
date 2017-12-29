package models

import ()

type Service struct {
	Id        int64 `json:id,omitempty`
	CompanyId int64 `json:companyId,omitempty`
}

type ServiceRequestBody struct {
	Scale        int                            `json:scale,omitempty`
	LaunchConfig ServiceRequestBodyLaunchConfig `json:launchConfig,omitempty`
	IdcIds       []int64                        `json:idcIds,omitempty`
	Metadata     map[string]string              `json:metadata,omitempty`
}

type ServiceRequestBodyLaunchConfig struct {
	Labels map[string]string `json:labels,omitempty`
}

type ServiceData struct {
	Fields ServiceDataField `json:fields,omitempty`
}

type ServiceDataField struct {
	Metadata map[string]string `json:metadata,omitempty`
}
