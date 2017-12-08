package models

import ()

type BsIdcResource struct {
	Id          int64 `json:id,omitempty`
	IdcId       int64 `json:idcId,omitempty`
	Cpu         int64 `json:cpu,omitempty`
	CpuUsed     int64 `json:cpuUsed,omitempty`
	CpuFree     int64 `json:cpuFree,omitempty`
	Memory      int64 `json:memory,omitempty`
	MemoryUsed  int64 `json:memoryUsed,omitempty`
	MemoryFree  int64 `json:meomoryFree,omitempty`
	Storage     int64 `json:storage,omitempty`
	StorageUsed int64 `json:storageUsed,omitempty`
	StorageFree int64 `json:storageFree,omitempty`
}
