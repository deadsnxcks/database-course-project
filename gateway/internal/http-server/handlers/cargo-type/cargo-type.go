package cargotype

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/deadsnxcks/dbcp/gateway/internal/http-server/dto"
	"github.com/deadsnxcks/dbcp/gateway/internal/http-server/response"
	"github.com/deadsnxcks/dbcp/gateway/internal/lib/logger/sl"
	cargotypev1 "github.com/deadsnxcks/dbcp/protos/gen/go/cargotype"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

const (
	opStart = "handlers.cargo-type"
	timeout = 5 * time.Second
)

type Handler struct {
	log    *slog.Logger
	client cargotypev1.CargoTypeServiceClient
}

func New(
	log *slog.Logger,
	client cargotypev1.CargoTypeServiceClient,
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

		resp, err := h.client.List(ctx, &cargotypev1.ListRequest{})
		if err != nil {
			response.GRPCError(w, r, log, err)

			return
		}

		render.JSON(w, r, dto.CargoTypesFromProto(resp.GetCargoTypes()))
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

		resp, err := h.client.Get(ctx, &cargotypev1.GetRequest{Id: id})
		if err != nil {
			response.GRPCError(w, r, log, err)

			return
		}

		render.JSON(w, r, dto.CargoTypeFromProto(resp.GetCargoType()))
	}
}

type createRequest struct {
	Title       string  `json:"title"`
	ProcessCost float64 `json:"processCost"`
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

		resp, err := h.client.Create(ctx, &cargotypev1.CreateRequest{
			Title:       req.Title,
			ProcessCost: req.ProcessCost,
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
	Title       *string  `json:"title,omitempty"`
	ProcessCost *float64 `json:"processCost,omitempty"`
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

		_, err := h.client.Update(ctx, &cargotypev1.UpdateRequest{
			Id:          id,
			Title:       req.Title,
			ProcessCost: req.ProcessCost,
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

		if _, err := h.client.Delete(ctx, &cargotypev1.DeleteRequest{Id: id}); err != nil {
			response.GRPCError(w, r, log, err)

			return
		}

		render.JSON(w, r, dto.Ok())
	}
}