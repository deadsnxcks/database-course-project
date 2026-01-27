package report

import (
	"context"
	"dbcp/internal/domain/models"
	reportv1 "dbcp/protos/gen/go/report"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Report interface {
	CargoDetailReport(ctx context.Context) ([]models.CargoDetailItem, error)
	CargoTypeReport(ctx context.Context) ([]models.CargoTypeItem, error)
}

type serverAPI struct {
	reportv1.UnimplementedReportServiceServer
	r Report
}

func Register(gRPCServer *grpc.Server, report Report) {
	reportv1.RegisterReportServiceServer(gRPCServer, &serverAPI{r: report})
}

func (s *serverAPI) GenerateUnloadedCargoReport(
	ctx context.Context,
	_ *reportv1.UnloadedCargoReportRequest,	
) (*reportv1.CargoDetailReport, error) {

	report, err := s.r.CargoDetailReport(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generrate report")
	}

	resp := make([]*reportv1.CargoDetailItem, 0, len(report))
	for _, item := range report {
		resp = append(resp,
			&reportv1.CargoDetailItem{
				CargoName: item.CargoName,
				WeightTons: item.Weight,
				CargoType: item.CargoType,
				VesselName: item.VesselName,
				UnloadingDate: item.UnloadingDate.Format(time.DateTime),
			},
		)
	}

	return &reportv1.CargoDetailReport{
		Items: resp,
	}, nil
}

func (s *serverAPI) GenerateCargoTypeSummaryReport(
	ctx context.Context,
	_ *reportv1.CargoTypeReportRequest,
) (*reportv1.CargoTypeReport, error) {
	report, err := s.r.CargoTypeReport(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generrate report")
	}

	resp := make([]*reportv1.CargoTypeItem, 0, len(report))
	for _, item := range report {
		resp = append(resp,
			&reportv1.CargoTypeItem{
				CargoTypeName: item.CargoTypeName,
				CargoCount: item.CargoCount,
				TotalWeightTons: item.TotalWeight,
				TotalVolumeM3: item.TotalVolume,
				ProcessCost: item.TotalProcessCost,
			},
		)
	}

	return &reportv1.CargoTypeReport{
		Items: resp,
	}, nil
}