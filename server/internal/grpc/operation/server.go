package operation

import (
	"context"
	operationv1 "github.com/deadsnxcks/dbcp/protos/gen/go/operation"
	"github.com/deadsnxcks/dbcp/server/internal/domain/models"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Operation interface {
	List(ctx context.Context) ([]models.Operation, error)
	Get(ctx context.Context, id int64) (models.Operation, error)
	Create(ctx context.Context, title string) (int64, error)
	Delete(ctx context.Context, id int64) error
	Update(
		ctx context.Context,
		id int64,
		title *string,
	) error
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

func (s *serverAPI) List(
	ctx context.Context,
	_ *operationv1.ListRequest,
) (*operationv1.ListResponse, error) {
	ops, err := s.operation.List(ctx)
	if err != nil {
		return nil, err
	}

	resp := make([]*operationv1.Operation, 0, len(ops))
	for _, op := range ops {
		resp = append(resp, toProtoOperation(op))
	}

	return &operationv1.ListResponse{
		Operations: resp,
	}, nil
}

func (s *serverAPI) Get(
	ctx context.Context,
	req *operationv1.GetRequest,
) (*operationv1.GetResponse, error) {

	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	op, err := s.operation.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &operationv1.GetResponse{
		Operation: toProtoOperation(op),
	}, nil
}

func (s *serverAPI) Create(
	ctx context.Context,
	req *operationv1.CreateRequest,
) (*operationv1.CreateResponse, error) {

	if req.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}

	id, err := s.operation.Create(ctx, req.GetTitle())
	if err != nil {
		return nil, err
	}

	return &operationv1.CreateResponse{
		Id: id,
	}, nil
}

func (s *serverAPI) Update(
	ctx context.Context,
	req *operationv1.UpdateRequest,
) (*operationv1.UpdateResponse, error) {
	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := s.operation.Update(ctx, req.GetId(), req.Title); err != nil {
		return nil, err
	}

	return &operationv1.UpdateResponse{}, nil
}

func (s *serverAPI) Delete(
	ctx context.Context,
	req *operationv1.DeleteRequest,
) (*operationv1.DeleteResponse, error) {

	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := s.operation.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}

	return &operationv1.DeleteResponse{}, nil
}

func toProtoOperation(
	op models.Operation,
) *operationv1.Operation {

	return &operationv1.Operation{
		Id:        op.ID,
		Title:     op.Title,
		CreatedAt: timestamppb.New(op.CreatedAt),
	}
}
