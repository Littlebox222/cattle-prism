package models

import (
	"time"
)

type BsGroup struct {
	Id            int64     `json:id,omitempty`
	Name          string    `json:name,omitempty`
	StackCount    int       `json:stackCount,omitempty`
	InstanceCount int       `json:instanceCount,omitempty`
	Created       time.Time `json:created,omitempty`
	Updated       time.Time `json:updated,omitempty`
	CompanyId     int64     `json:companyId,omitempty`
}

type BsIdc struct {
	Id                int64     `json:id,omitempty`
	Name              string    `json:name,omitempty`
	CarrierOperatorId int64     `json:carrierOperatorId,omitempty`
	AreaId            int64     `json:areaId,omitempty`
	State             string    `json:state,omitempty`
	Created           time.Time `json:created,omitempty`
	Updated           time.Time `json:updated,omitempty`
}

type BsCarrierOperator struct {
	Id   int64  `json:id,omitempty`
	Name string `json:name,omitempty`
}

type BsArea struct {
	Id   int64  `json:id,omitempty`
	Name string `json:name,omitempty`
}

type BsGroupIdcMap struct {
	Id        int64     `json:id,omitempty`
	GroupId   int64     `json:groupId,omitempty`
	IdcId     int64     `json:idcId,omitempty`
	CompanyId int64     `json:companyId,omitempty`
	Created   time.Time `json:created,omitempty`
}

type BsIdcHostMap struct {
	Id        int64     `json:id,omitempty`
	IdcId     int64     `json:idcId,omitempty`
	HostId    int64     `json:hostId,omitempty`
	CompanyId int64     `json:companyId,omitempty`
	Created   time.Time `json:created,omitempty`
}

///////////

type BsGroupIdcMapReqData struct {
	GroupId   int64   `json:groupId,omitempty`
	Idcs      []int64 `json:idcs,omitempty`
	CompanyId int64   `json:companyId,omitempty`
}
