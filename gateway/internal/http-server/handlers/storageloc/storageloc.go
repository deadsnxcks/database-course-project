package storageloc

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/deadsnxcks/dbcp/gateway/internal/http-server/dto"
	"github.com/deadsnxcks/dbcp/gateway/internal/http-server/response"
	"github.com/deadsnxcks/dbcp/gateway/internal/lib/logger/sl"
	storagelocv1 "github.com/deadsnxcks/dbcp/protos/gen/go/storageloc"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	opStart = "handlers.storageloc"
	timeout = 5 * time.Second
)

type Handler struct {
	log    *slog.Logger
	client storagelocv1.StorageLocationServiceClient
}

func New(
	log *slog.Logger,
	client storagelocv1.StorageLocationServiceClient,
) *Handler {
	return &Handler{log: log, client: client}
}

func (h *Handler) logger(r *http.Request, op string) *slog.Logger {
	return h.log.With(
		slog.String("op", opStart+"."+op),
		slog.String("req_id", middleware.GetReqID(r.Context())),
	)
}

func pathID(r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}

	return id, true
}

func (h *Handler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := h.logger(r, "List")

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		resp, err := h.client.List(ctx, &storagelocv1.ListRequest{})
		if err != nil {
			response.GRPCError(w, r, log, err)

			return
		}

		render.JSON(w, r, dto.StorageLocationsFromProto(resp.GetStorageLocations()))
	}
}

func (h *Handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := h.logger(r, "Get")

		id, ok := pathID(r)
		if !ok {
			response.BadRequest(w, r, "invalid id")

			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		resp, err := h.client.Get(ctx, &storagelocv1.GetRequest{Id: id})
		if err != nil {
			response.GRPCError(w, r, log, err)

			return
		}

		render.JSON(w, r, dto.StorageLocationFromProto(resp.GetStorageLocation()))
	}
}

type createRequest struct {
	CargoTypeID int64   `json:"cargoTypeId"`
	MaxWeight   float64 `json:"maxWeight"`
	MaxVolume   float64 `json:"maxVolume"`
}

func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := h.logger(r, "Create")

		var req createRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Debug("failed to decode body", sl.Err(err))
			response.BadRequest(w, r, "invalid request body")

			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		resp, err := h.client.Create(ctx, &storagelocv1.CreateRequest{
			CargoTypeId: req.CargoTypeID,
			MaxWeight:   req.MaxWeight,
			MaxVolume:   req.MaxVolume,
		})
		if err != nil {
			response.GRPCError(w, r, log, err)

			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, dto.Created{ID: resp.GetId()})
	}
}

type updateRequest struct {
	CargoTypeID *int64   `json:"cargoTypeId,omitempty"`
	MaxWeight   *float64 `json:"maxWeight,omitempty"`
	MaxVolume   *float64 `json:"maxVolume,omitempty"`
}

func (h *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := h.logger(r, "Update")

		id, ok := pathID(r)
		if !ok {
			response.BadRequest(w, r, "invalid id")

			return
		}

		var req updateRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Debug("failed to decode body", sl.Err(err))
			response.BadRequest(w, r, "invalid request body")

			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		_, err := h.client.Update(ctx, &storagelocv1.UpdateRequest{
			Id:          id,
			CargoTypeId: req.CargoTypeID,
			MaxWeight:   req.MaxWeight,
			MaxVolume:   req.MaxVolume,
		})
		if err != nil {
			response.GRPCError(w, r, log, err)

			return
		}

		render.JSON(w, r, dto.Ok())
	}
}

func (h *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := h.logger(r, "Delete")

		id, ok := pathID(r)
		if !ok {
			response.BadRequest(w, r, "invalid id")

			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		if _, err := h.client.Delete(ctx, &storagelocv1.DeleteRequest{Id: id}); err != nil {
			response.GRPCError(w, r, log, err)

			return
		}

		render.JSON(w, r, dto.Ok())
	}
}

type useRequest struct {
	CargoID         int64      `json:"cargoId"`
	DateOfPlacement *time.Time `json:"dateOfPlacement,omitempty"`
}

func (h *Handler) Use() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := h.logger(r, "Use")

		id, ok := pathID(r)
		if !ok {
			response.BadRequest(w, r, "invalid id")

			return
		}

		var req useRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Debug("failed to decode body", sl.Err(err))
			response.BadRequest(w, r, "invalid request body")

			return
		}

		grpcReq := &storagelocv1.UseRequest{
			StorageLocationId: id,
			CargoId:           req.CargoID,
		}
		if req.DateOfPlacement != nil {
			grpcReq.DateOfPlacement = timestamppb.New(*req.DateOfPlacement)
		}

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		if _, err := h.client.Use(ctx, grpcReq); err != nil {
			response.GRPCError(w, r, log, err)

			return
		}

		render.JSON(w, r, dto.Ok())
	}
}

func (h *Handler) Reset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := h.logger(r, "Reset")

		id, ok := pathID(r)
		if !ok {
			response.BadRequest(w, r, "invalid id")

			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		if _, err := h.client.Reset(ctx, &storagelocv1.ResetRequest{Id: id}); err != nil {
			response.GRPCError(w, r, log, err)

			return
		}

		render.JSON(w, r, dto.Ok())
	}
}
