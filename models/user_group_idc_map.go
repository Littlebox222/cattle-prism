package models

import ()

type BsUserGroupIdcMap struct {
	Id        int64 `json:id,omitempty`
	CompanyId int64 `json:companyId,omitempty`
	GroupId   int64 `json:groupId,omitempty`
	IdcId     int64 `json:idcId,omitempty`
}
