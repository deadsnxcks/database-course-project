package dto

import (
	"time"

	storagelocv1 "github.com/deadsnxcks/dbcp/protos/gen/go/storageloc"
)

type StorageLocation struct {
	ID          int64   `json:"id"`
	CargoTypeID int64   `json:"cargoTypeId"`
	MaxWeight   float64 `json:"maxWeight"`
	MaxVolume   float64 `json:"maxVolume"`
	CargoID         *int64     `json:"cargoId"`
	DateOfPlacement *time.Time `json:"dateOfPlacement"`
}

func StorageLocationFromProto(l *storagelocv1.StorageLocation) StorageLocation {
	out := StorageLocation{
		ID:          l.GetId(),
		CargoTypeID: l.GetCargoTypeId(),
		MaxWeight:   l.GetMaxWeight(),
		MaxVolume:   l.GetMaxVolume(),
	}

	if l.CargoId != nil {
		v := *l.CargoId
		out.CargoID = &v
	}

	if ts := l.GetDateOfPlacement(); ts != nil {
		t := ts.AsTime()
		out.DateOfPlacement = &t
	}

	return out
}

func StorageLocationsFromProto(ls []*storagelocv1.StorageLocation) []StorageLocation {
	out := make([]StorageLocation, 0, len(ls))
	for _, l := range ls {
		out = append(out, StorageLocationFromProto(l))
	}

	return out
}