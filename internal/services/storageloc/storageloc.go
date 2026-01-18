package storagelocservice

import (
	"context"
	"dbcp/internal/domain/models"
	"dbcp/internal/lib/logger/sl"
	"fmt"
	"log/slog"
	"time"
)

const (
	opStart = "services.storageloc"
)

type StorageLocService struct {
	log *slog.Logger
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
		log: log,
		slProvider: slProvider,
	}
}

func (s *StorageLocService) List(
	ctx context.Context,
) ([]models.StorageLocation, error) {
	const op = opStart + ".List"

	log := s.log.With(slog.String("op", op))
	log.Info("listing storage locations")

	locs, err := s.slProvider.StorageLocations(ctx)
	if err != nil {
		log.Error("failed to list storage locations", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return locs, nil
}


func (s *StorageLocService) Get(
	ctx context.Context, 
	id int64,
) (models.StorageLocation, error) {
	const op = opStart + ".Get"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("id", id),
	)

	if id <= 0 {
		return models.StorageLocation{}, fmt.Errorf("%s: invalid id", op)
	}

	loc, err := s.slProvider.StorageLocation(ctx, id)
	if err != nil {
		log.Error("failed to get storage location", sl.Err(err))
		return models.StorageLocation{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("storage location fetched")
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

	if cargoTypeID <= 0 {
		return 0, fmt.Errorf("%s: cargoTypeID is required", op)
	}
	if maxWeight <= 0 {
		return 0, fmt.Errorf("%s: maxWeight must be positive", op)
	}
	if maxVolume <= 0 {
		return 0, fmt.Errorf("%s: maxVolume must be positive", op)
	}

	id, err := s.slProvider.SaveStorageLoc(ctx, cargoTypeID, maxWeight, maxVolume)
	if err != nil {
		log.Error("failed to create storage location", sl.Err(err))
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

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("id", id),
	)

	if id <= 0 {
		return fmt.Errorf("%s: invalid id", op)
	}

	if err := s.slProvider.DeleteStorageLoc(ctx, id); err != nil {
		log.Error("failed to delete storage location", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("storage location deleted")
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

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("id", id),
	)

	if id <= 0 {
		return fmt.Errorf("%s: invalid id", op)
	}

	if err := s.slProvider.UpdateStorageLoc(
		ctx,
		id,
		cargoTypeID,
		maxWeight,
		maxVolume,
	); err != nil {
		log.Error("failed to update storage location", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("storage location updated")
	return nil
}

func (s *StorageLocService) Use(
	ctx context.Context,
	id int64,
	cargoID int64,
	date time.Time,
) error {
	const op = opStart + ".Use"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("id", id),
		slog.Int64("cargoID", cargoID),
	)

	if id <= 0 {
		return fmt.Errorf("%s: invalid storage location id", op)
	}
	if cargoID <= 0 {
		return fmt.Errorf("%s: invalid cargo id", op)
	}
	if date.IsZero() {
		date = time.Now()
	}

	if err := s.slProvider.UseStorageLoc(ctx, id, cargoID, date); err != nil {
		log.Error("failed to use storage location", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("storage location used")
	return nil
}

func (s *StorageLocService) Reset(
	ctx context.Context, 
	id int64,
) error {
	const op = opStart + ".Reset"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("id", id),
	)

	if id <= 0 {
		return fmt.Errorf("%s: invalid id", op)
	}

	if err := s.slProvider.ResetStorageLoc(ctx, id); err != nil {
		log.Error("failed to reset storage location", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("storage location reset")
	return nil
}