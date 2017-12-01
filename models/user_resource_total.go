package models

import ()

type BsUserResourceTotal struct {
	Id              int64 `json:id,omitempty`
	CompanyId       int64 `json:companyId,omitempty`
	ContainerTypeId int64 `json:containerTypeId,omitempty`
	IdcId           int64 `json:idcId,omitempty`
	Total           int   `json:total,omitempty`
	Used            int   `json:used,omitempty`
	Free            int   `json:free,omitempty`
}

type BsUserResourceTotalResponseData struct {
	Id            int64           `json:id,omitempty`
	CompanyId     int64           `json:companyId,omitempty`
	ContainerType BsContainerType `json:containerTypeId,omitempty`
	IdcId         int64           `json:idcId,omitempty`
	Total         int             `json:total,omitempty`
	Used          int             `json:used,omitempty`
	Free          int             `json:free,omitempty`
}

type BsUserResourceTotalQueryData struct {
	Id              int64  `json:id,omitempty`
	CompanyId       int64  `json:companyId,omitempty`
	IdcId           int64  `json:idcId,omitempty`
	Total           int    `json:total,omitempty`
	Used            int    `json:used,omitempty`
	Free            int    `json:free,omitempty`
	ContainerTypeId int64  `json:containerTypeId,omitempty`
	Name            string `json:name,omitempty`
	Cpu             int64  `json:cpu,omitempty`
	Memory          int64  `json:memory,omitempty`
	Storage         int64  `json:storage,omitempty`
}
