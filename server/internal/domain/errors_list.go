package domain

var (
	ErrVesselNotFound = NotFound("vessel not found")
	ErrVesselExists   = Conflict("vessel already exists")
	ErrVesselInUse    = InUse("vessel is in use")

	ErrCargoTypeNotFound = NotFound("cargo type not found")
	ErrCargoTypeExists   = Conflict("cargo type already exists")
	ErrCargoTypeInUse    = InUse("cargo type is in use")

	ErrCargoNotFound      = NotFound("cargo not found")
	ErrCargoExists        = Conflict("cargo already exists")
	ErrCargoInUse         = InUse("cargo is in use")
	ErrCargoAlreadyPlaced = Unprocessable("cargo is already placed in a storage location")

	ErrOperationNotFound = NotFound("operation not found")
	ErrOperationExists   = Conflict("operation already exists")
	ErrOperationInUse    = InUse("operation is in use")

	ErrStorageLocNotFound        = NotFound("storage location not found")
	ErrStorageLocInUse           = InUse("storage location is occupied")
	ErrStorageLocNotSuitable     = Unprocessable("storage location is not suitable")
	ErrStorageLocTypeNotSuitable = Unprocessable("cargo type is not allowed in this location")
	ErrStorageLocAlreadyEmpty    = Unprocessable("storage location is already empty")

	ErrOperCargoAlreadyExist = Conflict("operation with such cargo already exists")
	ErrOperCargoNotFound     = NotFound("operation with such cargo not found")

	ErrRelatedEntityNotFound = Unprocessable("one or more related entities not found")
)
