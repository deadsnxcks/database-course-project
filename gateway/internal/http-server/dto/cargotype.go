package dto

import (
	cargotypev1 "github.com/deadsnxcks/dbcp/protos/gen/go/cargotype"
)

type CargoType struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	ProcessCost float64 `json:"processCost"`
}

func CargoTypeFromProto(ct *cargotypev1.CargoType) CargoType {
	return CargoType{
		ID:          ct.GetId(),
		Title:       ct.GetTitle(),
		ProcessCost: ct.GetProcessCost(),
	}
}

func CargoTypesFromProto(cts []*cargotypev1.CargoType) []CargoType {
	out := make([]CargoType, 0, len(cts))
	for _, ct := range cts {
		out = append(out, CargoTypeFromProto(ct))
	}

	return out
}