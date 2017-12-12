package models

import ()

type Environment struct {
	Id        int64 `json:id,omitempty`
	CompanyId int64 `json:companyId,omitempty`
}
