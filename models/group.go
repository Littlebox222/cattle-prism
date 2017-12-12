package models

import (
	"time"
)

type BsGroup struct {
	Id            int64     `json:id,omitempty`
	GroupName     string    `json:groupName,omitempty`
	StackCount    int       `json:stackCount,omitempty`
	InstanceCount int       `json:instanceCount,omitempty`
	Created       time.Time `json:created,omitempty`
	Updated       time.Time `json:updated,omitempty`
	CompanyId     int64     `json:companyId,omitempty`
}

type BsGroupRequestData struct {
	GroupId   int64   `json:groupId,omitempty`
	GroupName string  `json:groupName,omitempty`
	CompanyId int64   `json:companyId,omitempty`
	IdcIds    []int64 `json:idcIds,omitempty`
}

type BsGroupResponseData struct {
	Group  BsGroup `json:group,omitempty`
	IdcIds []int64 `json:idcIds,omitempty`
}

type BsGroupResponseDataDetail struct {
	Group             BsGroup               `json:group,omitempty`
	Idcs              []BsIdcResponseData   `json:idcs,omitempty`
	GroupUserResource []BsGroupUserResource `json:groupUserResource,omitempty`
}
