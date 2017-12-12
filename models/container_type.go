package models

import ()

type BsContainerType struct {
	Id                int64  `json:id,omitempty`
	ContainerTypeName string `json:containerTypeName,omitempty`
	Cpu               int64  `json:cpu,omitempty`
	Memory            int64  `json:memory,omitempty`
	Storage           int64  `json:storage,omitempty`
}
