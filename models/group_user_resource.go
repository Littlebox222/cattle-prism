package models

import ()

type BsGroupUserResource struct {
	IdcId         int64           `json:idcId,omitempty`
	ContainerType BsContainerType `json:containerType,omitempty`
	Total         int             `json:total,omitempty`
	Used          int             `json:used,omitempty`
	Free          int             `json:free,omitempty`
}
