package models

import (
	"time"
)

type BsIdc struct {
	Id                int64     `json:id,omitempty`
	IdcName           string    `json:idcName,omitempty`
	CarrierOperatorId int64     `json:carrierOperatorId,omitempty`
	AreaId            int64     `json:areaId,omitempty`
	State             string    `json:state,omitempty`
	Created           time.Time `json:created,omitempty`
	Updated           time.Time `json:updated,omitempty`
}

type BsIdcResponseData struct {
	IdcId           int64                         `json:idcId,omitempty`
	IdcName         string                        `json:idcName,omitempty`
	CarrierOperator BsCarrierOperatorResponseData `json:carrierOperator,omitempty`
	Area            BsAreaResponseData            `json:area,omitempty`
	State           string                        `json:state,omitempty`
	Created         time.Time                     `json:created,omitempty`
	Updated         time.Time                     `json:updated,omitempty`
}
