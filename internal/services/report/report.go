package reportservice

import (
	"context"
	"dbcp/internal/domain/models"
	"dbcp/internal/lib/logger/sl"
	"fmt"
	"log/slog"
)

type ReportService struct {
	log *slog.Logger
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
		log: log,
		rProvider: rProvider,
	}
}

func (s *ReportService) CargoDetailReport(
	ctx context.Context,
) ([]models.CargoDetailItem, error) {
	const op = "services.report.CargoDetailReport"

	log := s.log.With(slog.String("op", op))
	log.Info("Generating report \"Cargo detail\"")

	cargoItems, err := s.rProvider.CargoDetailReport(ctx)
	if err != nil {
		log.Error("failed to generate report \"Cargo detail\"", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cargoItems, nil
}

func (s *ReportService) CargoTypeReport(
	ctx context.Context,
) ([]models.CargoTypeItem, error) {
	const op = "services.report.CargoDetailReport"

	log := s.log.With(slog.String("op", op))
	log.Info("Generating report \"Cargo detail\"")

	cargoTypeItems, err := s.rProvider.CargoTypeReport(ctx)
	if err != nil {
		log.Error("failed to generate report \"Cargo detail\"", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cargoTypeItems, nil
}