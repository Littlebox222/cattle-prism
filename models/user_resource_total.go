package models

import (
	"time"
)

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
	Id            int64             `json:id,omitempty`
	CompanyId     int64             `json:companyId,omitempty`
	ContainerType BsContainerType   `json:containerTypeId,omitempty`
	Idc           BsIdcResponseData `json:idc,omitempty`
	Total         int               `json:total,omitempty`
	Used          int               `json:used,omitempty`
	Free          int               `json:free,omitempty`
}

type BsUserResourceTotalCollectedByContainerTypeResponseData struct {
	ContainerType BsContainerType `json:containerTypeId,omitempty`
	Total         int             `json:total,omitempty`
	Used          int             `json:used,omitempty`
	Free          int             `json:free,omitempty`
}

type BsUserResourceTotalCollectedByIdcResponseData struct {
	Idc            BsIdcResponseData `json:idc,omitempty`
	ContainerTypes []BsUserResourceTotalCollectedByContainerTypeResponseData
}

type BsUserResourceTotalQueryData struct {
	Id                  int64     `json:id,omitempty`
	CompanyId           int64     `json:companyId,omitempty`
	IdcId               int64     `json:idcId,omitempty`
	IdcName             string    `json:idcName,omitempty`
	CarrierOperatorId   int64     `json:carrierOperatorId,omitempty`
	CarrierOperatorName string    `json:carrierOperatorName,omitempty`
	AreaId              int64     `json:areaId,omitempty`
	AreaName            string    `json:areaName,omitempty`
	State               string    `json:state,omitempty`
	Created             time.Time `json:created,omitempty`
	Updated             time.Time `json:updated,omitempty`
	Total               int       `json:total,omitempty`
	Used                int       `json:used,omitempty`
	Free                int       `json:free,omitempty`
	ContainerTypeId     int64     `json:containerTypeId,omitempty`
	ContainerTypeName   string    `json:containerTypeName,omitempty`
	Cpu                 int64     `json:cpu,omitempty`
	Memory              int64     `json:memory,omitempty`
	Storage             int64     `json:storage,omitempty`
}
