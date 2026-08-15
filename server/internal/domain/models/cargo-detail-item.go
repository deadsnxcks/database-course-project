package models

import "time"

type CargoDetailItem struct {
	CargoName 		string;
	Weight			float64;
	CargoType		string;
	VesselName		string;
	UnloadingDate	time.Time;
}