package storage

import "errors"

var (
	ErrVesselExists = errors.New("vessel already exists")
	ErrVesselNotFound = errors.New("vessel not found")
)