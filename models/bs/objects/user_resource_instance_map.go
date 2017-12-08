package models

import ()

type BsUserResourceInstanceMap struct {
	Id             int64  `json:id,omitempty`
	CompanyId      int64  `json:companyId,omitempty`
	UserResourceId int64  `json:userResourceId,omitempty`
	InstanceId     int64  `json:instanceId,omitempty`
	InstanceState  string `json:instanceState,omitempty`
	ServiceId      int64  `json:serviceId,omitempty`
	ServiceState   string `json:serviceState,omitempty`
	StackId        int64  `json:stackId,omitempty`
	StackState     string `json:stackState,omitempty`
}
