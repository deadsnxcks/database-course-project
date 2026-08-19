package vessel

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/deadsnxcks/dbcp/gateway/internal/http-server/dto"
	"github.com/deadsnxcks/dbcp/gateway/internal/http-server/response"
	"github.com/deadsnxcks/dbcp/gateway/internal/lib/logger/sl"
	vesselv1 "github.com/deadsnxcks/dbcp/protos/gen/go/vessel"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

const timeout = 5 * time.Second

type Handler struct {
	log    *slog.Logger
	client vesselv1.VesselServiceClient
}

func New(log *slog.Logger, client vesselv1.VesselServiceClient) *Handler {
	return &Handler{log: log, client: client}
}

func (h *Handler) logger(r *http.Request, op string) *slog.Logger {
	return h.log.With(
		slog.String("op", op),
		slog.String("req_id", middleware.GetReqID(r.Context())),
	)
}

func (h *Handler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := h.logger(r, "handlers.vessel.List")

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		resp, err := h.client.List(ctx, &vesselv1.ListRequest{})
		if err != nil {
			response.GRPCError(w, r, log, err)

			return
		}

		render.JSON(w, r, dto.VesselsFromProto(resp.GetVessels()))
	}
}

func (h *Handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := h.logger(r, "handlers.vessel.Get")

		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			response.BadRequest(w, r, "invalid id")

			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		resp, err := h.client.Get(ctx, &vesselv1.GetRequest{Id: id})
		if err != nil {
			response.GRPCError(w, r, log, err)

			return
		}

		render.JSON(w, r, dto.VesselFromProto(resp.GetVessel()))
	}
}

type createRequest struct {
	Title      string  `json:"title"`
	VesselType string  `json:"vesselType"`
	MaxLoad    float64 `json:"maxLoad"`
}

func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := h.logger(r, "handlers.vessel.Create")

		var req createRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Debug("failed to decode body", sl.Err(err))
			response.BadRequest(w, r, "invalid request body")

			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		resp, err := h.client.Create(ctx, &vesselv1.CreateRequest{
			Title:      req.Title,
			VesselType: req.VesselType,
			MaxLoad:    req.MaxLoad,
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
	Title      *string  `json:"title,omitempty"`
	VesselType *string  `json:"vesselType,omitempty"`
	MaxLoad    *float64 `json:"maxLoad,omitempty"`
}

func (h *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := h.logger(r, "handlers.vessel.Update")

		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
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

		_, err = h.client.Update(ctx, &vesselv1.UpdateRequest{
			Id:         id,
			Title:      req.Title,
			VesselType: req.VesselType,
			MaxLoad:    req.MaxLoad,
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
		log := h.logger(r, "handlers.vessel.Delete")

		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			response.BadRequest(w, r, "invalid id")

			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		if _, err := h.client.Delete(ctx, &vesselv1.DeleteRequest{Id: id}); err != nil {
			response.GRPCError(w, r, log, err)

			return
		}

		render.JSON(w, r, dto.Ok())
	}
}
