package cargotype

import (
	"context"
	cargotypev1 "github.com/deadsnxcks/dbcp/protos/gen/go/cargotype"
	"github.com/deadsnxcks/dbcp/server/internal/domain/models"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CargoType interface {
	List(ctx context.Context) ([]models.CargoType, error)
	Get(ctx context.Context, id int64) (models.CargoType, error)
	Create(ctx context.Context, vessel models.CargoType) (int64, error)
	Delete(ctx context.Context, id int64) error
	Update(
		ctx context.Context,
		id int64,
		title *string,
		processCost *float64,
	) error
}

type serverAPI struct {
	cargotypev1.UnimplementedCargoTypeServiceServer
	cargoType CargoType
}

func Register(gRPCServer *grpc.Server, cargoType CargoType) {
	cargotypev1.RegisterCargoTypeServiceServer(gRPCServer, &serverAPI{cargoType: cargoType})
}

func (s *serverAPI) List(
	ctx context.Context,
	_ *cargotypev1.ListRequest,
) (*cargotypev1.ListResponse, error) {
	ctList, err := s.cargoType.List(ctx)
	if err != nil {
		return nil, err
	}

	var pbList []*cargotypev1.CargoType
	for _, ct := range ctList {
		pbList = append(pbList, &cargotypev1.CargoType{
			Id:          ct.ID,
			Title:       ct.Title,
			ProcessCost: ct.ProcessCost,
		})
	}

	return &cargotypev1.ListResponse{CargoTypes: pbList}, nil
}

func (s *serverAPI) Get(
	ctx context.Context,
	req *cargotypev1.GetRequest,
) (*cargotypev1.GetResponse, error) {
	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	ct, err := s.cargoType.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &cargotypev1.GetResponse{
		CargoType: &cargotypev1.CargoType{
			Id:          ct.ID,
			Title:       ct.Title,
			ProcessCost: ct.ProcessCost,
		},
	}, nil
}

func (s *serverAPI) Create(
	ctx context.Context,
	req *cargotypev1.CreateRequest,
) (*cargotypev1.CreateResponse, error) {
	if req.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}
	if req.GetProcessCost() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "process_cost must be greater than 0")
	}

	id, err := s.cargoType.Create(ctx, models.CargoType{
		Title:       req.GetTitle(),
		ProcessCost: req.GetProcessCost(),
	})
	if err != nil {
		return nil, err
	}

	return &cargotypev1.CreateResponse{Id: id}, nil
}

func (s *serverAPI) Update(
	ctx context.Context,
	req *cargotypev1.UpdateRequest,
) (*cargotypev1.UpdateResponse, error) {
	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if req.Title != nil && *req.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "title must not be empty")
	}
	if req.ProcessCost != nil && *req.ProcessCost <= 0 {
		return nil, status.Error(codes.InvalidArgument, "process_cost must be greater than 0")
	}

	err := s.cargoType.Update(ctx, req.GetId(), req.Title, req.ProcessCost)
	if err != nil {
		return nil, err
	}

	return &cargotypev1.UpdateResponse{}, nil
}

func (s *serverAPI) Delete(
	ctx context.Context,
	req *cargotypev1.DeleteRequest,
) (*cargotypev1.DeleteResponse, error) {
	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := s.cargoType.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &cargotypev1.DeleteResponse{}, nil
}
