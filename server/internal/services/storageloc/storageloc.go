package storagelocservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/deadsnxcks/dbcp/server/internal/domain"
	"github.com/deadsnxcks/dbcp/server/internal/domain/models"
	"github.com/deadsnxcks/dbcp/server/internal/storage"
)

const (
	opStart = "services.storageloc"
)

type StorageLocService struct {
	log        *slog.Logger
	slProvider StorageLocProvider
}

type StorageLocProvider interface {
	StorageLocations(ctx context.Context) ([]models.StorageLocation, error)
	SaveStorageLoc(
		ctx context.Context,
		cargoTypeID int64,
		maxWeight float64,
		maxVolume float64,
	) (int64, error)
	DeleteStorageLoc(ctx context.Context, id int64) error
	StorageLocation(ctx context.Context, id int64) (models.StorageLocation, error)
	UpdateStorageLoc(
		ctx context.Context,
		id int64,
		cargoTypeID *int64,
		maxWeight *float64,
		maxVolume *float64,
	) error
	UseStorageLoc(
		ctx context.Context,
		storageLocID int64,
		cargoID int64,
		date time.Time,
	) error
	ResetStorageLoc(ctx context.Context, id int64) error
}

func New(
	log *slog.Logger,
	slProvider StorageLocProvider,
) *StorageLocService {
	return &StorageLocService{
		log:        log,
		slProvider: slProvider,
	}
}

func (s *StorageLocService) List(
	ctx context.Context,
) ([]models.StorageLocation, error) {
	const op = opStart + ".List"

	locs, err := s.slProvider.StorageLocations(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return locs, nil
}

func (s *StorageLocService) Get(
	ctx context.Context,
	id int64,
) (models.StorageLocation, error) {
	const op = opStart + ".Get"

	loc, err := s.slProvider.StorageLocation(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrStorageLocNotFound) {
			return models.StorageLocation{}, fmt.Errorf("%s: %w", op, domain.ErrStorageLocNotFound)
		}
		return models.StorageLocation{}, fmt.Errorf("%s: %w", op, err)
	}

	return loc, nil
}

func (s *StorageLocService) Create(
	ctx context.Context,
	cargoTypeID int64,
	maxWeight float64,
	maxVolume float64,
) (int64, error) {
	const op = opStart + ".Create"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("cargoTypeID", cargoTypeID),
	)

	id, err := s.slProvider.SaveStorageLoc(ctx, cargoTypeID, maxWeight, maxVolume)
	if err != nil {
		if errors.Is(err, storage.ErrRelatedEntityNotFound) {
			return 0, fmt.Errorf("%s: %w", op, domain.ErrRelatedEntityNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("storage location created", slog.Int64("id", id))
	return id, nil
}

func (s *StorageLocService) Delete(
	ctx context.Context,
	id int64,
) error {
	const op = opStart + ".Delete"

	if err := s.slProvider.DeleteStorageLoc(ctx, id); err != nil {
		switch {
		case errors.Is(err, storage.ErrStorageLocInUse):
			return fmt.Errorf("%s: %w", op, domain.ErrStorageLocInUse)
		case errors.Is(err, storage.ErrStorageLocNotFound):
			return fmt.Errorf("%s: %w", op, domain.ErrStorageLocNotFound)
		default:
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (s *StorageLocService) Update(
	ctx context.Context,
	id int64,
	cargoTypeID *int64,
	maxWeight *float64,
	maxVolume *float64,
) error {
	const op = opStart + ".Update"

	if err := s.slProvider.UpdateStorageLoc(
		ctx,
		id,
		cargoTypeID,
		maxWeight,
		maxVolume,
	); err != nil {
		switch {
		case errors.Is(err, storage.ErrRelatedEntityNotFound):
			return fmt.Errorf("%s: %w", op, domain.ErrRelatedEntityNotFound)
		case errors.Is(err, storage.ErrStorageLocNotFound):
			return fmt.Errorf("%s: %w", op, domain.ErrStorageLocNotFound)
		default:
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (s *StorageLocService) Use(
	ctx context.Context,
	id int64,
	cargoID int64,
	date time.Time,
) error {
	const op = opStart + ".Use"

	if err := s.slProvider.UseStorageLoc(ctx, id, cargoID, date); err != nil {
		switch {
		case errors.Is(err, storage.ErrStorageLocNotFound):
			return fmt.Errorf("%s: %w", op, domain.ErrStorageLocNotFound)
		case errors.Is(err, storage.ErrStorageLocInUse):
			return fmt.Errorf("%s: %w", op, domain.ErrStorageLocInUse)
		case errors.Is(err, storage.ErrCargoNotFound):
			return fmt.Errorf("%s: %w", op, domain.ErrCargoNotFound)
		case errors.Is(err, storage.ErrCargoAlreadyPlaced):
			return fmt.Errorf("%s: %w", op, domain.ErrCargoAlreadyPlaced)
		case errors.Is(err, storage.ErrStorageLocNotSuitable):
			return fmt.Errorf("%s: %w", op, domain.ErrStorageLocNotSuitable)
		default:
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (s *StorageLocService) Reset(
	ctx context.Context,
	id int64,
) error {
	const op = opStart + ".Reset"

	if err := s.slProvider.ResetStorageLoc(ctx, id); err != nil {
		switch {
		case errors.Is(err, storage.ErrStorageLocNotFound):
			return fmt.Errorf("%s: %w", op, domain.ErrStorageLocNotFound)
		case errors.Is(err, storage.ErrStorageLocAlreadyEmpty):
			return fmt.Errorf("%s: %w", op, domain.ErrStorageLocAlreadyEmpty)
		default:
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}
