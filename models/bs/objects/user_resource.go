package models

import ()

type BsUserResource struct {
	Id              int64 `json:id,omitempty`
	CompanyId       int64 `json:companyId,omitempty`
	ContainerTypeId int64 `json:containerTypeId,omitempty`
	IdcId           int64 `json:idcId,omitempty`
	CanUse          bool  `json:canUse,omitempty`
}
