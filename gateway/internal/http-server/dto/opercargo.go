package dto

import (
	opercargov1 "github.com/deadsnxcks/dbcp/protos/gen/go/opercargo"
)

type OperationCargo struct {
	OperationID int64 `json:"operationId"`
	CargoID     int64 `json:"cargoId"`
}

func OperationCargoFromProto(oc *opercargov1.OperationCargo) OperationCargo {
	return OperationCargo{
		OperationID: oc.GetOperationId(),
		CargoID:     oc.GetCargoId(),
	}
}

func OperationCargosFromProto(ocs []*opercargov1.OperationCargo) []OperationCargo {
	out := make([]OperationCargo, 0, len(ocs))
	for _, oc := range ocs {
		out = append(out, OperationCargoFromProto(oc))
	}

	return out
}