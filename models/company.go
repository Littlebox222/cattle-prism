package models

import ()

type Company struct {
	Id   int64  `json:id,omitempty`
	Name string `json:name,omitempty`
	Uuid string `json:uuid,omitempty"`
}
