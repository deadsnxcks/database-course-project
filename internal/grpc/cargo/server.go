package cargo

import (
	"context"
	"dbcp/internal/domain/models"
	"dbcp/internal/storage"
	cargov1 "dbcp/protos/gen/go/cargo"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Cargo interface {
	List(ctx context.Context) ([]models.Cargo, error)
	ListByVesselID(ctx context.Context, vesselID int64) ([]models.Cargo, error)
	ListByTypeID(ctx context.Context, typeID int64) ([]models.Cargo, error)
	Get(ctx context.Context, id int64) (models.Cargo, error)
	Create(ctx context.Context, cargo models.Cargo) (int64, error)
	Delete(ctx context.Context, id int64) (error)
	Update(
		ctx context.Context, 
		id int64,
		title *string, 
		cargoTypeID *int64, 
		weight *float64,
		volume *float64,
		vesselID *int64,
	) (error)
}

type serverAPI struct {
	cargov1.UnimplementedCargoServiceServer
	cargo Cargo
}

func Register(gRPCServer *grpc.Server, cargo Cargo) {
	cargov1.RegisterCargoServiceServer(gRPCServer, &serverAPI{cargo: cargo})
}

func (s *serverAPI) ListCargos(
	ctx context.Context,
	_ *cargov1.ListCargosRequest,
) (*cargov1.ListCargosResponse, error) {

	cargos, err := s.cargo.List(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list cargos")
	}

	resp := make([]*cargov1.Cargo, 0, len(cargos))
	for _, c := range cargos {
		resp = append(resp, toProtoCargo(c))
	}

	return &cargov1.ListCargosResponse{Cargos: resp}, nil
}

func (s *serverAPI) ListCargosByVesselID(
	ctx context.Context,
	req *cargov1.GetCargosByVesselIDRequest,
) (*cargov1.GetCargosByVesselIDResponse, error) {

	if req.GetVesselId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "vessel_id is required")
	}

	cargos, err := s.cargo.ListByVesselID(ctx, req.GetVesselId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := make([]*cargov1.Cargo, 0, len(cargos))
	for _, c := range cargos {
		resp = append(resp, toProtoCargo(c))
	}

	return &cargov1.GetCargosByVesselIDResponse{Cargos: resp}, nil
}

func (s *serverAPI) GetCargo(
	ctx context.Context,
	req *cargov1.GetCargoRequest,
) (*cargov1.GetCargoResponse, error) {

	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	cargo, err := s.cargo.Get(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, storage.ErrCargoNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &cargov1.GetCargoResponse{Cargo: toProtoCargo(cargo)}, nil
}

func (s *serverAPI) CreateCargo(
	ctx context.Context,
	req *cargov1.CreateCargoRequest,
) (*cargov1.CreateCargoResponse, error) {

	if req.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}
	if req.GetTypeId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "cargo_type_id is required")
	}
	if req.GetWeight() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "weight must be positive")
	}

	cargo := models.Cargo{
		Title:      req.GetTitle(),
		TypeID:		req.GetTypeId(),
		Weight:     req.GetWeight(),
		Volume:		req.GetVolume(),
		VesselID:	req.GetVesselId(),
	}

	id, err := s.cargo.Create(ctx, cargo)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &cargov1.CreateCargoResponse{Id: id}, nil
}

func (s *serverAPI) UpdateCargo(
	ctx context.Context,
	req *cargov1.UpdateCargoRequest,
) (*cargov1.UpdateCargoResponse, error) {

	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	var (
		title       *string
		cargoTypeID *int64
		weight      *float64
		volume      *float64
		vesselID    *int64
	)

	if req.GetTitle() != "" {
		t := req.GetTitle()
		title = &t
	}
	if req.GetTypeId() > 0 {
		ct := req.GetTypeId()
		cargoTypeID = &ct
	}
	if req.GetWeight() > 0 {
		w := req.GetWeight()
		weight = &w
	}
	if req.GetVolume() > 0 {
		v := req.GetVolume()
		volume = &v
	}
	if req.GetVesselId() > 0 {
		vid := req.GetVesselId()
		vesselID = &vid
	}

	if err := s.cargo.Update(
		ctx,
		req.GetId(),
		title,
		cargoTypeID,
		weight,
		volume,
		vesselID,
	); err != nil {
		if errors.Is(err, storage.ErrCargoNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &cargov1.UpdateCargoResponse{}, nil
}

func (s *serverAPI) DeleteCargo(
	ctx context.Context,
	req *cargov1.DeleteCargoRequest,
) (*cargov1.DeleteCargoResponse, error) {

	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := s.cargo.Delete(ctx, req.GetId()); err != nil {
		if errors.Is(err, storage.ErrCargoNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &cargov1.DeleteCargoResponse{}, nil
}

func toProtoCargo(c models.Cargo) *cargov1.Cargo {
    return &cargov1.Cargo{
        Id:         c.ID,
        Title:      c.Title,
        TypeId: 	c.TypeID,
        Weight:    	c.Weight,
		Volume: 	c.Volume,
		VesselId: 	c.VesselID,	
    }
}