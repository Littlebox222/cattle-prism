package models

import ()

type BsIdcHostMap struct {
	Id     int64 `json:id,omitempty`
	IdcId  int64 `json:idcId,omitempty`
	HostId int64 `json:hostId,omitempty`
}
