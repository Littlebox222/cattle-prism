package models

import ()

type BsUserResourceInstanceMap struct {
	Id             int64 `json:id,omitempty`
	CompanyId      int64 `json:companyId,omitempty`
	UserResourceId int64 `json:userResourceId,omitempty`
	InstanceId     int64 `json:instanceId,omitempty`
	ServiceId      int64 `json:serviceId,omitempty`
	StackId        int64 `json:stackId,omitempty`
}
