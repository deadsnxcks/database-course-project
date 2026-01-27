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
	Delete(ctx context.Context, id int64) error
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
	) error
	Use(
		ctx context.Context,
		id int64,
		cargoId int64,
		date time.Time,
	) error
	Reset(ctx context.Context, id int64) error
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

func (s *serverAPI) List(
	ctx context.Context,
	_ *storagelocv1.ListRequest,
) (*storagelocv1.ListResponse, error) {

	locs, err := s.storageLocation.List(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list storage locations")
	}

	resp := make([]*storagelocv1.StorageLocation, 0, len(locs))
	for _, l := range locs {
		resp = append(resp, toProtoStorageLoc(l))
	}

	return &storagelocv1.ListResponse{
		StorageLocations: resp,
	}, nil
}

func (s *serverAPI) Get(
	ctx context.Context,
	req *storagelocv1.GetRequest,
) (*storagelocv1.GetResponse, error) {

	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	loc, err := s.storageLocation.Get(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, storage.ErrStorageLocNotFound) {
			return nil, status.Error(codes.NotFound, "storage location not found")
		}
		return nil, status.Error(codes.Internal, "failed to get storage location")
	}

	return &storagelocv1.GetResponse{
		StorageLocation: toProtoStorageLoc(loc),
	}, nil
}

func (s *serverAPI) Create(
	ctx context.Context,
	req *storagelocv1.CreateRequest,
) (*storagelocv1.CreateResponse, error) {

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
		switch {
		case errors.Is(err, storage.ErrRelatedEntityNotFound):
			return nil, status.Error(codes.NotFound, "related entity not found")
		default:
			return nil, status.Error(codes.Internal, "failed to create storage location")
		}
	}

	return &storagelocv1.CreateResponse{Id: id}, nil
}

func (s *serverAPI) Update(
	ctx context.Context,
	req *storagelocv1.UpdateRequest,
) (*storagelocv1.UpdateResponse, error) {

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
		switch {
		case errors.Is(err, storage.ErrStorageLocNotFound):
			return nil, status.Error(codes.NotFound, "storage location not found")
		case errors.Is(err, storage.ErrRelatedEntityNotFound):
			return nil, status.Error(codes.NotFound, "related entity not found")
		default:
			return nil, status.Error(codes.Internal, "failed to update storage location")
		}
	}

	return &storagelocv1.UpdateResponse{}, nil
}

func (s *serverAPI) Delete(
	ctx context.Context,
	req *storagelocv1.DeleteRequest,
) (*storagelocv1.DeleteResponse, error) {

	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := s.storageLocation.Delete(ctx, req.GetId())
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrStorageLocInUse):
			return nil, status.Error(codes.FailedPrecondition, "storage location is used")
		case errors.Is(err, storage.ErrStorageLocNotFound):
			return nil, status.Error(codes.NotFound, "storage location not found")
		default:
			return nil, status.Error(codes.Internal, "failed to delete storage location")
		}
	}

	return &storagelocv1.DeleteResponse{}, nil
}

func (s *serverAPI) Use(
	ctx context.Context,
	req *storagelocv1.UseRequest,
) (*storagelocv1.UseResponse, error) {

	if req.GetStorageLocationId() <= 0 {
		return nil, status.Error(codes.InvalidArgument,
			"storage_location_id must be positive")
	}

	if req.GetCargoId() <= 0 {
		return nil, status.Error(codes.InvalidArgument,
			"cargo_id must be positive")
	}

	date := time.Now()
	if req.GetDateOfPlacement() != nil {
		date = req.GetDateOfPlacement().AsTime()
		
		_, offset := time.Now().Zone()
		adjustedNow := time.Now().UTC().Add(time.Duration(offset) * time.Second)	
				
		if date.After(adjustedNow) {
			return nil, status.Error(codes.InvalidArgument,
				"date_of_placement cannot be in the future")
		}
	}

	err := s.storageLocation.Use(ctx, req.GetStorageLocationId(), req.GetCargoId(), date)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrStorageLocNotFound):
			return nil, status.Error(codes.NotFound, "storage location not found")
		case errors.Is(err, storage.ErrCargoNotFound):
			return nil, status.Error(codes.NotFound, "cargo not found")
		case errors.Is(err, storage.ErrStorageLocInUse):
			return nil, status.Error(codes.FailedPrecondition, "storage location is already in use")
		case errors.Is(err, storage.ErrStorageLocNotSuitable):
			return nil, status.Error(codes.FailedPrecondition, "storage location not suitable for this cargo")
		case errors.Is(err, storage.ErrCargoAlreadyPlaced):
			return nil, status.Error(codes.FailedPrecondition, "cargo is already placed in a storage location")
		default:
			return nil, status.Error(codes.Internal, "failed to use storage location")
		}
	}

	return &storagelocv1.UseResponse{}, nil
}

func (s *serverAPI) Reset(
	ctx context.Context,
	req *storagelocv1.ResetRequest,
) (*storagelocv1.ResetResponse, error) {

	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id is must be positive")
	}

	err := s.storageLocation.Reset(ctx, req.GetId())
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrStorageLocAlreadyEmpty):
			return nil, status.Error(codes.FailedPrecondition, "stiorage location is already empty")
		case errors.Is(err, storage.ErrStorageLocNotFound):
			return nil, status.Error(codes.NotFound, "storage location not found")
		default:
			return nil, status.Error(codes.Internal, "failed to reset storage location")
		}
	}

	return &storagelocv1.ResetResponse{}, nil
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
		Id:              sl.ID,
		CargoTypeId:     sl.CargoTypeID,
		MaxWeight:       sl.MaxWeight,
		MaxVolume:       sl.MaxVolume,
		CargoId:         cargoID,
		DateOfPlacement: date,
	}
}
