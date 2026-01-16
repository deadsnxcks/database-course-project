package models

import "time"

type StorageLocation struct {
	ID 				int64
	CargoTypeID 	int64
	MaxWeight		float64
	MaxVolume		float64
	CargoID			int64
	DateOfPlacement	time.Time
}