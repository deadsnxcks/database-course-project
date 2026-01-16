package storage

import "errors"

var (
	ErrVesselExists = errors.New("vessel already exists")
	ErrVesselNotFound = errors.New("vessel not found")

	ErrCargoTypeExists = errors.New("cargo type already exists")
	ErrCargoTypeNotFound = errors.New("cargo type not found")
	ErrCargoTypeInUse = errors.New("cargo type used")

	ErrOperationExists = errors.New("operation already exists")
	ErrOperationNotFound = errors.New("operation not found")
	ErrOperationInUse = errors.New("operation used")

	ErrCargoNotFound = errors.New("cargo not found")
	ErrCargoInUse = errors.New("cargo used")

	ErrStorageLocNotFound = errors.New("storage location not found")
	ErrStorageLocNotSuitable = errors.New("storage location not suitable")
	ErrStorageLocAlreadyEmpty = errors.New("storage location already empty")

	ErrOperCargoAlreadyExist = errors.New("an operation with such cargo already exists")
	ErrOperCargoNotFound = errors.New("an opearation with such cargo not found")

	ErrRelatedEntityNotFound = errors.New("related entity not found")
	ErrForeignKeyViolation = errors.New("foreign key violation")
)