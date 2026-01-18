package cargotype

import (
	"context"
	"dbcp/internal/domain/models"
	"dbcp/internal/storage"
	cargotypev1 "dbcp/protos/gen/go/cargotype"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CargoType interface {
	List(ctx context.Context) ([]models.CargoType, error)
	Get(ctx context.Context, id int64) (models.CargoType, error)
	Create(ctx context.Context, vessel models.CargoType) (int64, error)
	Delete(ctx context.Context, id int64) (error)
	Update(
		ctx context.Context, 
		id int64,
		title *string, 
		processCost *float64,
	) (error)
}

type serverAPI struct {
	cargotypev1.UnimplementedCargoTypeServiceServer
	cargoType CargoType
}

func Register(gRPCServer *grpc.Server, cargoType CargoType) {
	cargotypev1.RegisterCargoTypeServiceServer(gRPCServer, &serverAPI{cargoType: cargoType})
}

func (s *serverAPI) ListCargoTypes(
	ctx context.Context,
	req *cargotypev1.ListCargoTypesRequest,
) (*cargotypev1.ListCargoTypesResponse, error) {
	ctList, err := s.cargoType.List(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list cargo types")
	}

	var pbList []*cargotypev1.CargoType
	for _, ct := range ctList {
		pbList = append(pbList, &cargotypev1.CargoType{
			Id:          ct.ID,
			Title:       ct.Title,
			ProcessCost: ct.ProcessCost,
		})
	}

	return &cargotypev1.ListCargoTypesResponse{CargoTypes: pbList}, nil
}

func (s *serverAPI) GetCargoType(
	ctx context.Context,
	req *cargotypev1.GetCargoTypeRequest,
) (*cargotypev1.GetCargoTypeResponse, error) {
	ct, err := s.cargoType.Get(ctx, req.GetId())
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrCargoTypeNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &cargotypev1.GetCargoTypeResponse{
		CargoType: &cargotypev1.CargoType{
			Id:          ct.ID,
			Title:       ct.Title,
			ProcessCost: ct.ProcessCost,
		},
	}, nil
}

func (s *serverAPI) CreateCargoType(
	ctx context.Context,
	req *cargotypev1.CreateCargoTypeRequest,
) (*cargotypev1.CreateCargoTypeResponse, error) {
	if req.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}

	ct := models.CargoType{
		Title:       req.GetTitle(),
		ProcessCost: req.GetProcessCost(),
	}

	id, err := s.cargoType.Create(ctx, ct)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrCargoTypeExists):
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &cargotypev1.CreateCargoTypeResponse{Id: id}, nil
}


func (s *serverAPI) UpdateCargoType(
	ctx context.Context,
	req *cargotypev1.UpdateCargoTypeRequest,
) (*cargotypev1.UpdateCargoTypeResponse, error) {
	var title *string
	if req.GetTitle() != "" {
		t := req.GetTitle()
		title = &t
	}

	var processCost *float64
	if req.GetProcessCost() != 0 {
		pc := req.GetProcessCost()
		processCost = &pc
	}

	err := s.cargoType.Update(ctx, req.GetId(), title, processCost)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrCargoTypeNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, storage.ErrCargoTypeExists):
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &cargotypev1.UpdateCargoTypeResponse{}, nil
}

func (s *serverAPI) DeleteCargoType(
	ctx context.Context,
	req *cargotypev1.DeleteCargoTypeRequest,
) (*cargotypev1.DeleteCargoTypeResponse, error) {
	err := s.cargoType.Delete(ctx, req.GetId())
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrCargoTypeNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &cargotypev1.DeleteCargoTypeResponse{}, nil
}