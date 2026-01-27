package opercargo

import (
	"context"
	"dbcp/internal/domain/models"
	"dbcp/internal/storage"
	opercargov1 "dbcp/protos/gen/go/opercargo"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OperationCargo interface {
	List(ctx context.Context) ([]models.OperationCargo, error)
	Create(ctx context.Context, operID, cargoID int64) (error)
	Delete(ctx context.Context, operID, cargoID int64) (error)
}

type serverAPI struct {
	opercargov1.UnimplementedOperationCargoServiceServer
	operCargo OperationCargo
}

func Register(gRPCServer *grpc.Server, operCargo OperationCargo) {
	opercargov1.RegisterOperationCargoServiceServer(
		gRPCServer,
		&serverAPI{
			operCargo: operCargo,
		},
	)
}

func (s *serverAPI) List(
	ctx context.Context,
	req *opercargov1.ListRequest,
) (*opercargov1.ListResponse, error) {
	operCargos, err := s.operCargo.List(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list operation cargos")
	}

	resp := make([]*opercargov1.OperationCargo, 0, len(operCargos))
	for _, oc := range operCargos {
		resp = append(resp, &opercargov1.OperationCargo{
			OperationId: oc.OperationID,
			CargoId:     oc.CargoID,
		})
	}

	return &opercargov1.ListResponse{
		OperationsCargos: resp,
	}, nil
}

func (s *serverAPI) Create(
	ctx context.Context,
	req *opercargov1.CreateRequest,
) (*opercargov1.CreateResponse, error) {
	if req.GetOperationId() <= 0 || req.GetCargoId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "operation_id and cargo_id are required")
	}

	err := s.operCargo.Create(ctx, req.GetOperationId(), req.GetCargoId())
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrOperCargoAlreadyExist):
			return nil, status.Error(codes.AlreadyExists, "operarion_cargo already exists")
		case errors.Is(err, storage.ErrRelatedEntityNotFound):
			return nil, status.Error(codes.FailedPrecondition, "one or more related entities not found")
		default:
			return nil, status.Error(codes.Internal, "failed to create operation_cargo")
		}
	}

	return &opercargov1.CreateResponse{}, nil
}

func (s *serverAPI) Delete(
	ctx context.Context,
	req *opercargov1.DeleteRequest,
) (*opercargov1.DeleteResponse, error) {
	if req.GetOperationId() <= 0 || req.GetCargoId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "operation_id and cargo_id are required")
	}

	err := s.operCargo.Delete(ctx, req.GetOperationId(), req.GetCargoId())
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrOperCargoNotFound):
			return nil, status.Error(codes.NotFound, "operation_cargo not found")
		default:
			return nil, status.Error(codes.Internal, "operation_cargo deletion failed")
		}
	}

	return &opercargov1.DeleteResponse{}, nil
}