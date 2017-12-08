package models

import ()

type Service struct {
	Id        int64 `json:id,omitempty`
	CompanyId int64 `json:companyId,omitempty`
}
