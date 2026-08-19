package cargo

import (
	"context"
	cargov1 "github.com/deadsnxcks/dbcp/protos/gen/go/cargo"
	"github.com/deadsnxcks/dbcp/server/internal/domain/models"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Cargo interface {
	List(ctx context.Context) ([]models.Cargo, error)
	Get(ctx context.Context, id int64) (models.Cargo, error)
	Create(ctx context.Context, cargo models.Cargo) (int64, error)
	Delete(ctx context.Context, id int64) error
	Update(
		ctx context.Context,
		id int64,
		title *string,
		cargoTypeID *int64,
		weight *float64,
		volume *float64,
		vesselID *int64,
	) error
}

type serverAPI struct {
	cargov1.UnimplementedCargoServiceServer
	cargo Cargo
}

func Register(gRPCServer *grpc.Server, cargo Cargo) {
	cargov1.RegisterCargoServiceServer(gRPCServer, &serverAPI{cargo: cargo})
}

func (s *serverAPI) List(
	ctx context.Context,
	_ *cargov1.ListRequest,
) (*cargov1.ListResponse, error) {

	cargos, err := s.cargo.List(ctx)
	if err != nil {
		return nil, err
	}

	resp := make([]*cargov1.Cargo, 0, len(cargos))
	for _, c := range cargos {
		resp = append(resp, toProtoCargo(c))
	}

	return &cargov1.ListResponse{Cargos: resp}, nil
}

func (s *serverAPI) Get(
	ctx context.Context,
	req *cargov1.GetRequest,
) (*cargov1.GetResponse, error) {
	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	cargo, err := s.cargo.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &cargov1.GetResponse{Cargo: toProtoCargo(cargo)}, nil
}

func (s *serverAPI) Create(
	ctx context.Context,
	req *cargov1.CreateRequest,
) (*cargov1.CreateResponse, error) {

	if req.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}
	if req.GetTypeId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "cargo_type_id is required")
	}
	if req.GetWeight() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "weight must be positive")
	}
	if req.GetVolume() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "volume must be positive")
	}
	if req.GetVesselId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "vessel_id is required")
	}

	cargo := models.Cargo{
		Title:    req.GetTitle(),
		TypeID:   req.GetTypeId(),
		Weight:   req.GetWeight(),
		Volume:   req.GetVolume(),
		VesselID: req.GetVesselId(),
	}

	id, err := s.cargo.Create(ctx, cargo)
	if err != nil {
		return nil, err
	}

	return &cargov1.CreateResponse{Id: id}, nil
}

func (s *serverAPI) Update(
	ctx context.Context,
	req *cargov1.UpdateRequest,
) (*cargov1.UpdateResponse, error) {
	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if req.Weight != nil && *req.Weight <= 0 {
		return nil, status.Error(codes.InvalidArgument, "weight must be positive")
	}

	if err := s.cargo.Update(
		ctx,
		req.GetId(),
		req.Title,
		req.TypeId,
		req.Weight,
		req.Volume,
		req.VesselId,
	); err != nil {
		return nil, err
	}

	return &cargov1.UpdateResponse{}, nil
}

func (s *serverAPI) Delete(
	ctx context.Context,
	req *cargov1.DeleteRequest,
) (*cargov1.DeleteResponse, error) {

	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := s.cargo.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}

	return &cargov1.DeleteResponse{}, nil
}

func toProtoCargo(c models.Cargo) *cargov1.Cargo {
	return &cargov1.Cargo{
		Id:       c.ID,
		Title:    c.Title,
		TypeId:   c.TypeID,
		Weight:   c.Weight,
		Volume:   c.Volume,
		VesselId: c.VesselID,
	}
}
