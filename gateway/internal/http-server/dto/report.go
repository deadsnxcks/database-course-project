package dto

import (
	reportv1 "github.com/deadsnxcks/dbcp/protos/gen/go/report"
)

type CargoDetailItem struct {
	CargoName  string  `json:"cargoName"`
	Weight     float64 `json:"weight"`
	CargoType  string  `json:"cargoType"`
	VesselName string  `json:"vesselName"`
	UnloadDate string  `json:"unloadDate"`
}

func CargoDetailItemsFromProto(items []*reportv1.CargoDetailItem) []CargoDetailItem {
	out := make([]CargoDetailItem, 0, len(items))
	for _, i := range items {
		out = append(out, CargoDetailItem{
			CargoName:  i.GetCargoName(),
			Weight:     i.GetWeightTons(),
			CargoType:  i.GetCargoType(),
			VesselName: i.GetVesselName(),
			UnloadDate: i.GetUnloadingDate(),
		})
	}

	return out
}

type CargoTypeItem struct {
	CargoTypeName string  `json:"cargoTypeName"`
	Count         int64   `json:"count"`
	Weight        float64 `json:"weight"`
	Volume        float64 `json:"volume"`
	ProcessCost   float64 `json:"processCost"`
}

func CargoTypeItemsFromProto(items []*reportv1.CargoTypeItem) []CargoTypeItem {
	out := make([]CargoTypeItem, 0, len(items))
	for _, i := range items {
		out = append(out, CargoTypeItem{
			CargoTypeName: i.GetCargoTypeName(),
			Count:         int64(i.GetCargoCount()),
			Weight:        i.GetTotalWeightTons(),
			Volume:        i.GetTotalVolumeM3(),
			ProcessCost:   i.GetProcessCost(),
		})
	}

	return out
}