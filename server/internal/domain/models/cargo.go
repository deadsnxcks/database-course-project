package models

type Cargo struct {
	ID 			int64
	Title 		string
	TypeID 		int64
	Weight 		float64
	Volume		float64
	VesselID	int64
}