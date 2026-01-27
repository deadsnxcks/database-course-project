package vesselservice

import (
	"context"
	"dbcp/internal/domain/models"
	"dbcp/internal/lib/logger/sl"
	"fmt"
	"log/slog"
)

const (
	opStart = "services.vessel"
)

type VesselService struct {
	log *slog.Logger
	vProvider VesselProvider
}

type VesselProvider interface {
	Vessels(ctx context.Context) ([]models.Vessel, error)
	SaveVessel(ctx context.Context, vessel models.Vessel) (int64, error)
	DeleteVessel(ctx context.Context, id int64) error
	Vessel(ctx context.Context, id int64) (models.Vessel, error)
	UpdateVessel(
		ctx context.Context,
		id int64,
		title *string,
		vesselType *string,
		maxLoad *float64,
	) error
}

func New(
	log *slog.Logger,
	vProvider VesselProvider,
) *VesselService {
	return &VesselService{
		log: log,
		vProvider: vProvider,
	}
}

func (v *VesselService) List(ctx context.Context) ([]models.Vessel, error) {
	const op = opStart + ".List"

	log := v.log.With(
		slog.String("op", op),
	)

	log.Info("Listing vessels")

	vessels, err := v.vProvider.Vessels(ctx)
	if err != nil {
		log.Error("failed to list vessels", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return vessels, nil
}

func (v *VesselService) Get(ctx context.Context, id int64) (models.Vessel, error) {
	const op = opStart + ".Get"

	log := v.log.With(slog.String("op", op), slog.Int64("id", id))

	if id <= 0 {
		return models.Vessel{}, fmt.Errorf("%s: invalid id", op)
	}

	vessel, err := v.vProvider.Vessel(ctx, id)
	if err != nil {
		log.Error("failed to get vessel", sl.Err(err))
		return models.Vessel{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Vessel geted", slog.Int64("id", id))
	return vessel, nil
}

func (v *VesselService) Create(ctx context.Context, vessel models.Vessel) (int64, error) {
	const op = opStart + ".Create"

	log := v.log.With(slog.String("op", op), slog.String("title", vessel.Title))

	// Валидация
	if vessel.Title == "" {
		return 0, fmt.Errorf("%s: title is required", op)
	}
	if vessel.VesselType == "" {
		return 0, fmt.Errorf("%s: vesselType is required", op)
	}
	if vessel.MaxLoad <= 0 {
		return 0, fmt.Errorf("%s: maxLoad must be positive", op)
	}

	id, err := v.vProvider.SaveVessel(ctx, vessel)
	if err != nil {
		log.Error("failed to create vessel", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Vessel created", slog.Int64("id", id))
	return id, nil
}

func (v *VesselService) Delete(ctx context.Context, id int64) error {
	const op = opStart + ".Delete"

	log := v.log.With(slog.String("op", op), slog.Int64("id", id))

	if id <= 0 {
		return fmt.Errorf("%s: invalid id", op)
	}

	err := v.vProvider.DeleteVessel(ctx, id)
	if err != nil {
		log.Error("failed to delete vessel", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Vessel deleted")
	return nil
}

func (v *VesselService) Update(
	ctx context.Context,
	id int64,
	title *string,
	vesselType *string,
	maxLoad *float64,
) error {
	const op = opStart + ".Update"

	log := v.log.With(slog.String("op", op), slog.Int64("id", id))

	if id <= 0 {
		return fmt.Errorf("%s: invalid id", op)
	}

	err := v.vProvider.UpdateVessel(ctx, id, title, vesselType, maxLoad)
	if err != nil {
		log.Error("failed to update vessel", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Vessel updated")
	return nil
}