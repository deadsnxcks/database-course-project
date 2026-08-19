package dto

import (
	cargov1 "github.com/deadsnxcks/dbcp/protos/gen/go/cargo"
)

type Cargo struct {
	ID       int64   `json:"id"`
	Title    string  `json:"title"`
	TypeID   int64   `json:"typeId"`
	Weight   float64 `json:"weight"`
	Volume   float64 `json:"volume"`
	VesselID int64   `json:"vesselId"`
}

func CargoFromProto(c *cargov1.Cargo) Cargo {
	return Cargo{
		ID:       c.GetId(),
		Title:    c.GetTitle(),
		TypeID:   c.GetTypeId(),
		Weight:   c.GetWeight(),
		Volume:   c.GetVolume(),
		VesselID: c.GetVesselId(),
	}
}

func CargosFromProto(cs []*cargov1.Cargo) []Cargo {
	out := make([]Cargo, 0, len(cs))
	for _, c := range cs {
		out = append(out, CargoFromProto(c))
	}

	return out
}