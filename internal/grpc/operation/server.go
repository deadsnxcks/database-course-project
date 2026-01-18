package operation

import (
	"context"
	"dbcp/internal/domain/models"
	"dbcp/internal/storage"
	operationv1 "dbcp/protos/gen/go/operation"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Operation interface {
	List(ctx context.Context) ([]models.Operation, error)
	Get(ctx context.Context, id int64) (models.Operation, error)
	Create(ctx context.Context, title string) (int64, error)
	Delete(ctx context.Context, id int64) (error)
	Update(
		ctx context.Context, 
		id int64,
		title *string, 
	) (error)
}

type serverAPI struct {
	operationv1.UnimplementedOperationServiceServer
	operation Operation
}

func Register(gRPCServer *grpc.Server, operation Operation) {
	operationv1.RegisterOperationServiceServer(
		gRPCServer, 
		&serverAPI{
			operation: operation,
		},
	)
}

func (s *serverAPI) ListOperations(
	ctx context.Context,
	_ *operationv1.ListOperationsRequest,
) (*operationv1.ListOperationsResponse, error) {

	ops, err := s.operation.List(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list operations")
	}

	resp := make([]*operationv1.Operation, 0, len(ops))
	for _, op := range ops {
		resp = append(resp, toProtoOperation(op))
	}

	return &operationv1.ListOperationsResponse{
		Operations: resp,
	}, nil
}

func (s *serverAPI) GetOperation(
	ctx context.Context,
	req *operationv1.GetOperationRequest,
) (*operationv1.GetOperationResponse, error) {

	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	op, err := s.operation.Get(ctx, req.GetId())
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrOperationNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &operationv1.GetOperationResponse{
		Operation: toProtoOperation(op),
	}, nil
}

func (s *serverAPI) CreateOperation(
	ctx context.Context,
	req *operationv1.CreateOperationRequest,
) (*operationv1.CreateOperationResponse, error) {

	if req.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}

	id, err := s.operation.Create(ctx, req.GetTitle())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &operationv1.CreateOperationResponse{
		Id: id,
	}, nil
}

func (s *serverAPI) UpdateOperation(
	ctx context.Context,
	req *operationv1.UpdateOperationRequest,
) (*operationv1.UpdateOperationResponse, error) {

	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	var title *string
	if req.GetTitle() != "" {
		t := req.GetTitle()
		title = &t
	}

	if err := s.operation.Update(
		ctx,
		req.GetId(),
		title,
	); err != nil {
		switch {
		case errors.Is(err, storage.ErrOperationNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &operationv1.UpdateOperationResponse{}, nil
}

func (s *serverAPI) DeleteOperation(
	ctx context.Context,
	req *operationv1.DeleteOperationRequest,
) (*operationv1.DeleteOperationResponse, error) {

	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := s.operation.Delete(ctx, req.GetId()); err != nil {
		switch {
		case errors.Is(err, storage.ErrOperationNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &operationv1.DeleteOperationResponse{}, nil
}

func toProtoOperation(
	op models.Operation,
) *operationv1.Operation {

	return &operationv1.Operation{
		Id:    op.ID,
		Title: op.Title,
		CreatedAt: timestamppb.New(op.CreatedAt),
	}
}