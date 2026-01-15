package storage

import "errors"

var (
	ErrVesselExists = errors.New("vessel already exists")
	ErrVesselNotFound = errors.New("vessel not found")

	ErrCargoTypeExists = errors.New("cargo type already exists")
	ErrCargoNotFound = errors.New("cargo type not found")
)