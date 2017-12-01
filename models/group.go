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
