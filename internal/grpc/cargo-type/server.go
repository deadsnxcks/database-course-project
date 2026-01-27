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

func (s *serverAPI) List(
	ctx context.Context,
	req *cargotypev1.ListRequest,
) (*cargotypev1.ListResponse, error) {
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

	return &cargotypev1.ListResponse{CargoTypes: pbList}, nil
}

func (s *serverAPI) Get(
	ctx context.Context,
	req *cargotypev1.GetRequest,
) (*cargotypev1.GetResponse, error) {
	ct, err := s.cargoType.Get(ctx, req.GetId())
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrCargoTypeNotFound):
			return nil, status.Error(codes.NotFound, "cargo type not found")
		default:
			return nil, status.Error(codes.Internal, "failed to get cargo type")
		}
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

	ct := models.CargoType{
		Title:       req.GetTitle(),
		ProcessCost: req.GetProcessCost(),
	}

	id, err := s.cargoType.Create(ctx, ct)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrCargoTypeExists):
			return nil, status.Error(codes.AlreadyExists, "cargo type already exists")
		default:
			return nil, status.Error(codes.Internal, "failed to create cargo type")
		}
	}

	return &cargotypev1.CreateResponse{Id: id}, nil
}


func (s *serverAPI) Update(
	ctx context.Context,
	req *cargotypev1.UpdateRequest,
) (*cargotypev1.UpdateResponse, error) {
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
			return nil, status.Error(codes.NotFound, "cargo type not found")
		case errors.Is(err, storage.ErrCargoTypeExists):
			return nil, status.Error(codes.AlreadyExists, "cargo type already exists")
		default:
			return nil, status.Error(codes.Internal, "failed to update cargp type")
		}
	}

	return &cargotypev1.UpdateResponse{}, nil
}

func (s *serverAPI) Delete(
	ctx context.Context,
	req *cargotypev1.DeleteRequest,
) (*cargotypev1.DeleteResponse, error) {
	err := s.cargoType.Delete(ctx, req.GetId())
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrCargoTypeInUse):
			return nil, status.Error(codes.FailedPrecondition, "cargo type is used")
		case errors.Is(err, storage.ErrCargoTypeNotFound):
			return nil, status.Error(codes.NotFound, "cargo type not found")
		default:
			return nil, status.Error(codes.Internal, "failed to delete cargo type")
		}
	}

	return &cargotypev1.DeleteResponse{}, nil
}