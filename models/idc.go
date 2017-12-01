package models

import (
	"time"
)

type BsIdc struct {
	Id                int64     `json:id,omitempty`
	Name              string    `json:name,omitempty`
	CarrierOperatorId int64     `json:carrierOperatorId,omitempty`
	AreaId            int64     `json:areaId,omitempty`
	State             string    `json:state,omitempty`
	Created           time.Time `json:created,omitempty`
	Updated           time.Time `json:updated,omitempty`
}
