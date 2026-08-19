package vesselservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/deadsnxcks/dbcp/server/internal/domain"
	"github.com/deadsnxcks/dbcp/server/internal/domain/models"
	"github.com/deadsnxcks/dbcp/server/internal/storage"
)

const (
	opStart = "services.vessel"
)

type VesselService struct {
	log       *slog.Logger
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
		log:       log,
		vProvider: vProvider,
	}
}

func (v *VesselService) List(ctx context.Context) ([]models.Vessel, error) {
	const op = opStart + ".List"

	vessels, err := v.vProvider.Vessels(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return vessels, nil
}

func (v *VesselService) Get(ctx context.Context, id int64) (models.Vessel, error) {
	const op = opStart + ".Get"

	vessel, err := v.vProvider.Vessel(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrVesselNotFound) {
			return models.Vessel{}, fmt.Errorf("%s: %w", op, domain.ErrVesselNotFound)
		}

		return models.Vessel{}, fmt.Errorf("%s: %w", op, err)
	}

	return vessel, nil
}

func (v *VesselService) Create(ctx context.Context, vessel models.Vessel) (int64, error) {
	const op = opStart + ".Create"

	id, err := v.vProvider.SaveVessel(ctx, vessel)
	if err != nil {
		if errors.Is(err, storage.ErrVesselExists) {
			return 0, fmt.Errorf("%s: %w", op, domain.ErrVesselExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	v.log.Info("vessel created", slog.Int64("id", id), slog.String("title", vessel.Title))

	return id, nil
}

func (v *VesselService) Delete(ctx context.Context, id int64) error {
	const op = opStart + ".Delete"

	if err := v.vProvider.DeleteVessel(ctx, id); err != nil {
		switch {
		case errors.Is(err, storage.ErrVesselInUse):
			return fmt.Errorf("%s: %w", op, domain.ErrVesselInUse)
		case errors.Is(err, storage.ErrVesselNotFound):
			return fmt.Errorf("%s: %w", op, domain.ErrVesselNotFound)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

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

	if err := v.vProvider.UpdateVessel(ctx, id, title, vesselType, maxLoad); err != nil {
		switch {
		case errors.Is(err, storage.ErrVesselNotFound):
			return fmt.Errorf("%s: %w", op, domain.ErrVesselNotFound)
		case errors.Is(err, storage.ErrVesselExists):
			return fmt.Errorf("%s: %w", op, domain.ErrVesselExists)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
