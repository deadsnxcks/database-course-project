package vessel

import (
	"context"
	"dbcp/internal/domain/models"
	"dbcp/internal/storage"
	dbcpv1 "dbcp/protos/gen/go/dbcp"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Vessel interface {
	List(ctx context.Context) ([]models.Vessel, error)
	Get(ctx context.Context, id int64) (models.Vessel, error)
	Create(ctx context.Context, vessel models.Vessel) (int64, error)
	Delete(ctx context.Context, id int64) (error)
	Update(
		ctx context.Context, 
		id int64,
		title *string, 
		vesselType *string, 
		maxLoad *float64,
	) (error)
}

type serverAPI struct {
	dbcpv1.UnimplementedDbcpServiceServer
	vessel Vessel
}

func Register(gRPCServer *grpc.Server, vessel Vessel) {
	dbcpv1.RegisterDbcpServiceServer(gRPCServer, &serverAPI{vessel: vessel})
}

func (s *serverAPI) ListVessels(
	ctx context.Context,
	lv *dbcpv1.ListVesselsRequest,
	) (*dbcpv1.ListVesselsResponse, error) {

	vessels, err := s.vessel.List(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list vessels")
	}

	resp := make([]*dbcpv1.Vessel, 0, len(vessels))
	for _, v := range vessels {
		resp = append(resp, &dbcpv1.Vessel{
			Id:         v.ID,
			Title:      v.Title,
			VesselType: v.VesselType,
			MaxLoad:    v.MaxLoad,
		})
	}

	return &dbcpv1.ListVesselsResponse{Vessels: resp}, nil
}

func (s *serverAPI) GetVessel(
	ctx context.Context,
	gv *dbcpv1.GetVesselRequest,
) (*dbcpv1.GetVesselResponse, error) {
	if gv.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	vessel, err := s.vessel.Get(ctx, gv.GetId())
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrVesselExists):
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &dbcpv1.GetVesselResponse{Vessel: toProtoVessel(vessel)}, nil
}

func (s *serverAPI) CreateVessel(
	ctx context.Context,
	cv *dbcpv1.CreateVesselRequest,
) (*dbcpv1.CreateVesselResponse, error) {
	if cv.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required") 
	}
	if cv.GetMaxLoad() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "max load must be greater than 0")
	}
	if cv.GetVesselType() == "" {
		return nil, status.Error(codes.InvalidArgument, "vessel type is required")
	}

	vessel := models.Vessel{
		Title: cv.GetTitle(),
		MaxLoad: cv.GetMaxLoad(),
		VesselType: cv.GetVesselType(),
	}
	id, err := s.vessel.Create(ctx, vessel)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create vessel")
	}

	return &dbcpv1.CreateVesselResponse{Id: id}, nil
}

func (s *serverAPI) UpdateVessel(
	ctx context.Context,
	uv *dbcpv1.UpdateVesselRequest,
) (*dbcpv1.UpdateVesselResponse, error) {
	if uv.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}


	var title *string
	if uv.GetTitle() != "" {
		t := uv.GetTitle()
		title = &t
	}

	var vesselType *string
	if uv.GetVesselType() != "" {
		vt := uv.GetVesselType()
		vesselType = &vt
	}

	var maxLoad *float64
	if uv.GetMaxLoad() > 0 {
		ml := uv.GetMaxLoad()
		maxLoad = &ml
	}
	
	err := s.vessel.Update(
		ctx,
		uv.GetId(),
		title,
		vesselType,
		maxLoad,
	)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrVesselNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &dbcpv1.UpdateVesselResponse{}, nil
}

func (s *serverAPI) DeleteVessel(
	ctx context.Context,
	dv *dbcpv1.DeleteVesselRequest,
) (*dbcpv1.DeleteVesselResponse, error) {
	if dv.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := s.vessel.Delete(ctx, dv.GetId()); err != nil {
		if errors.Is(err, storage.ErrVesselNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &dbcpv1.DeleteVesselResponse{}, nil
}

func toProtoVessel(v models.Vessel) *dbcpv1.Vessel {
    return &dbcpv1.Vessel{
        Id:         v.ID,
        Title:      v.Title,
        VesselType: v.VesselType,
        MaxLoad:    v.MaxLoad,
    }
}