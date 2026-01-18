package storageloc

import (
	"context"
	"dbcp/internal/domain/models"
	"dbcp/internal/storage"
	storagelocv1 "dbcp/protos/gen/go/storageloc"
	"errors"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type StorageLoc interface {
	List(ctx context.Context) ([]models.StorageLocation, error)
	Get(ctx context.Context, id int64) (models.StorageLocation, error)
	Delete(ctx context.Context, id int64) (error)
	Create(
		ctx context.Context,
		cargoTypeID int64, 
		maxWeight float64,
		maxVolume float64,
	) (int64, error)
	Update(
		ctx context.Context, 
		id int64,
		cargoTypeID *int64, 
		maxWeight *float64, 
		maxVolume *float64,
	) (error)
	Use(
		ctx context.Context,
		id int64,
		cargoId int64,
		date time.Time,
	) (error)
	Reset(ctx context.Context, id int64) (error)
}

type serverAPI struct {
	storagelocv1.UnimplementedStorageLocationServiceServer
	storageLocation StorageLoc
}

func Register(gRPCServer *grpc.Server, storageLocation StorageLoc) {
	storagelocv1.RegisterStorageLocationServiceServer(
		gRPCServer, 
		&serverAPI{
			storageLocation: storageLocation,
		},
	)
}

func (s *serverAPI) ListStorageLocations(
	ctx context.Context,
	_ *storagelocv1.ListStorageLocationsRequest,
) (*storagelocv1.ListStorageLocationsResponse, error) {

	locs, err := s.storageLocation.List(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list storage locations")
	}

	resp := make([]*storagelocv1.StorageLocation, 0, len(locs))
	for _, l := range locs {
		resp = append(resp, toProtoStorageLoc(l))
	}

	return &storagelocv1.ListStorageLocationsResponse{
		StorageLocations: resp,
	}, nil
}

func (s *serverAPI) GetStorageLocation(
	ctx context.Context,
	req *storagelocv1.GetStorageLocationRequest,
) (*storagelocv1.GetStorageLocationResponse, error) {

	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	loc, err := s.storageLocation.Get(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, storage.ErrStorageLocNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &storagelocv1.GetStorageLocationResponse{
		StorageLocation: toProtoStorageLoc(loc),
	}, nil
}

func (s *serverAPI) CreateStorageLocation(
	ctx context.Context,
	req *storagelocv1.CreateStorageLocationRequest,
) (*storagelocv1.CreateStorageLocationResponse, error) {

	if req.GetCargoTypeId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "cargo_type_id is required")
	}
	if req.GetMaxWeight() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "max_weight must be > 0")
	}
	if req.GetMaxVolume() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "max_volume must be > 0")
	}

	id, err := s.storageLocation.Create(
		ctx,
		req.GetCargoTypeId(),
		req.GetMaxWeight(),
		req.GetMaxVolume(),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &storagelocv1.CreateStorageLocationResponse{Id: id}, nil
}

func (s *serverAPI) UpdateStorageLocation(
	ctx context.Context,
	req *storagelocv1.UpdateStorageLocationRequest,
) (*storagelocv1.UpdateStorageLocationResponse, error) {

	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	var cargoTypeID *int64
	if req.GetCargoTypeId() > 0 {
		id := req.GetCargoTypeId()
		cargoTypeID = &id
	}

	var maxWeight *float64
	if req.GetMaxWeight() > 0 {
		w := req.GetMaxWeight()
		maxWeight = &w
	}

	var maxVolume *float64
	if req.GetMaxVolume() > 0 {
		v := req.GetMaxVolume()
		maxVolume = &v
	}

	err := s.storageLocation.Update(
		ctx,
		req.GetId(),
		cargoTypeID,
		maxWeight,
		maxVolume,
	)
	if err != nil {
		if errors.Is(err, storage.ErrStorageLocNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &storagelocv1.UpdateStorageLocationResponse{}, nil
}

func (s *serverAPI) UseStorageLocation(
	ctx context.Context,
	req *storagelocv1.UseStorageLocationRequest,
) (*storagelocv1.UseStorageLocationResponse, error) {

	if req.GetStorageLocationId() <= 0 || req.GetCargoId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id and cargo_id are required")
	}

	date := time.Now()
	if req.GetDateOfPlacement() != nil {
		date = req.GetDateOfPlacement().AsTime()
	}

	err := s.storageLocation.Use(ctx, req.GetStorageLocationId(), req.GetCargoId(), date)
	if err != nil {
		if errors.Is(err, storage.ErrStorageLocNotSuitable) {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &storagelocv1.UseStorageLocationResponse{}, nil
}

func (s *serverAPI) ResetStorageLocation(
	ctx context.Context,
	req *storagelocv1.ResetStorageLocationRequest,
) (*storagelocv1.ResetStorageLocationResponse, error) {

	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := s.storageLocation.Reset(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, storage.ErrStorageLocNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &storagelocv1.ResetStorageLocationResponse{}, nil
}

func toProtoStorageLoc(
	sl models.StorageLocation,
) *storagelocv1.StorageLocation {
	var cargoID *int64
	if sl.CargoID != nil {
		cargoID = sl.CargoID
	}

	var date *timestamppb.Timestamp
	if sl.DateOfPlacement != nil {
		date = timestamppb.New(*sl.DateOfPlacement)
	}

	return &storagelocv1.StorageLocation{
		Id: 				sl.ID,
		CargoTypeId: 		sl.CargoTypeID,
		MaxWeight: 			sl.MaxWeight,
		MaxVolume: 			sl.MaxVolume,
		CargoId: 			cargoID,
		DateOfPlacement: 	date,
	}
}