package reportservice

import (
	"context"
	"fmt"
	"github.com/deadsnxcks/dbcp/server/internal/domain/models"
	"log/slog"
)

type ReportService struct {
	log       *slog.Logger
	rProvider ReportProvider
}

type ReportProvider interface {
	CargoDetailReport(ctx context.Context) ([]models.CargoDetailItem, error)
	CargoTypeReport(ctx context.Context) ([]models.CargoTypeItem, error)
}

func New(
	log *slog.Logger,
	rProvider ReportProvider,
) *ReportService {
	return &ReportService{
		log:       log,
		rProvider: rProvider,
	}
}

func (s *ReportService) CargoDetailReport(
	ctx context.Context,
) ([]models.CargoDetailItem, error) {
	const op = "services.report.CargoDetailReport"

	cargoItems, err := s.rProvider.CargoDetailReport(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cargoItems, nil
}

func (s *ReportService) CargoTypeReport(
	ctx context.Context,
) ([]models.CargoTypeItem, error) {
	const op = "services.report.CargoTypeReport"

	cargoTypeItems, err := s.rProvider.CargoTypeReport(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cargoTypeItems, nil
}
