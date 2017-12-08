package models

import ()

type BsArea struct {
	Id       int64  `json:id,omitempty`
	AreaName string `json:areaName,omitempty`
}

type BsAreaResponseData struct {
	AreaId   int64  `json:areaId,omitempty`
	AreaName string `json:areaName,omitempty`
}
