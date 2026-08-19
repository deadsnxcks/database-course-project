package dto

import (
	vesselv1 "github.com/deadsnxcks/dbcp/protos/gen/go/vessel"
)

type Vessel struct {
	ID         int64   `json:"id"`
	Title      string  `json:"title"`
	VesselType string  `json:"vesselType"`
	MaxLoad    float64 `json:"maxLoad"`
}

func VesselFromProto(v *vesselv1.Vessel) Vessel {
	return Vessel{
		ID:         v.GetId(),
		Title:      v.GetTitle(),
		VesselType: v.GetVesselType(),
		MaxLoad:    v.GetMaxLoad(),
	}
}

func VesselsFromProto(vs []*vesselv1.Vessel) []Vessel {
	out := make([]Vessel, 0, len(vs))
	for _, v := range vs {
		out = append(out, VesselFromProto(v))
	}

	return out
}