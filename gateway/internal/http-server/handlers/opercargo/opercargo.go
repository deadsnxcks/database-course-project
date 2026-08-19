package opercargo

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/deadsnxcks/dbcp/gateway/internal/http-server/dto"
	"github.com/deadsnxcks/dbcp/gateway/internal/http-server/response"
	"github.com/deadsnxcks/dbcp/gateway/internal/lib/logger/sl"
	opercargov1 "github.com/deadsnxcks/dbcp/protos/gen/go/opercargo"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

const (
	opStart = "handlers.opercargo"
	timeout = 5 * time.Second
)

type Handler struct {
	log    *slog.Logger
	client opercargov1.OperationCargoServiceClient
}

func New(
	log *slog.Logger,
	client opercargov1.OperationCargoServiceClient,
) *Handler {
	return &Handler{log: log, client: client}
}

func (h *Handler) logger(r *http.Request, op string) *slog.Logger {
	return h.log.With(
		slog.String("op", opStart+"."+op),
		slog.String("req_id", middleware.GetReqID(r.Context())),
	)
}

type operCargoRequest struct {
	OperationID int64 `json:"operationId"`
	CargoID     int64 `json:"cargoId"`
}

func (h *Handler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := h.logger(r, "List")

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		resp, err := h.client.List(ctx, &opercargov1.ListRequest{})
		if err != nil {
			response.GRPCError(w, r, log, err)

			return
		}

		render.JSON(w, r, dto.OperationCargosFromProto(resp.GetOperationsCargos()))
	}
}

func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := h.logger(r, "Create")

		var req operCargoRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Debug("failed to decode body", sl.Err(err))
			response.BadRequest(w, r, "invalid request body")

			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		_, err := h.client.Create(ctx, &opercargov1.CreateRequest{
			OperationId: req.OperationID,
			CargoId:     req.CargoID,
		})
		if err != nil {
			response.GRPCError(w, r, log, err)

			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, dto.Ok())
	}
}

func (h *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := h.logger(r, "Delete")

		var req operCargoRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Debug("failed to decode body", sl.Err(err))
			response.BadRequest(w, r, "invalid request body")

			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		_, err := h.client.Delete(ctx, &opercargov1.DeleteRequest{
			OperationId: req.OperationID,
			CargoId:     req.CargoID,
		})
		if err != nil {
			response.GRPCError(w, r, log, err)

			return
		}

		render.JSON(w, r, dto.Ok())
	}
}