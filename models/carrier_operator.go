package models

import ()

type BsCarrierOperator struct {
	Id                  int64  `json:id,omitempty`
	CarrierOperatorName string `json:carrierOperatorName,omitempty`
}

type BsCarrierOperatorResponseData struct {
	CarrierOperatorId   int64  `json:carrierOperatorId,omitempty`
	CarrierOperatorName string `json:carrierOperatorName,omitempty`
}
