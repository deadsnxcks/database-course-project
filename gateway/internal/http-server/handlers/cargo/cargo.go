package cargo

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/deadsnxcks/dbcp/gateway/internal/http-server/dto"
	"github.com/deadsnxcks/dbcp/gateway/internal/http-server/response"
	"github.com/deadsnxcks/dbcp/gateway/internal/lib/logger/sl"
	cargov1 "github.com/deadsnxcks/dbcp/protos/gen/go/cargo"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

const (
	opStart = "handlers.cargo"
	timeout = 5 * time.Second
)

type Handler struct {
	log    *slog.Logger
	client cargov1.CargoServiceClient
}

func New(
	log *slog.Logger,
	client cargov1.CargoServiceClient,
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

		resp, err := h.client.List(ctx, &cargov1.ListRequest{})
		if err != nil {
			response.GRPCError(w, r, log, err)

			return
		}

		render.JSON(w, r, dto.CargosFromProto(resp.GetCargos()))
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

		resp, err := h.client.Get(ctx, &cargov1.GetRequest{Id: id})
		if err != nil {
			response.GRPCError(w, r, log, err)

			return
		}

		render.JSON(w, r, dto.CargoFromProto(resp.GetCargo()))
	}
}

type createRequest struct {
	Title    string  `json:"title"`
	TypeID   int64   `json:"typeId"`
	Weight   float64 `json:"weight"`
	Volume   float64 `json:"volume"`
	VesselID int64   `json:"vesselId"`
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

		resp, err := h.client.Create(ctx, &cargov1.CreateRequest{
			Title:    req.Title,
			TypeId:   req.TypeID,
			Weight:   req.Weight,
			Volume:   req.Volume,
			VesselId: req.VesselID,
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
	Title    *string  `json:"title,omitempty"`
	TypeID   *int64   `json:"typeId,omitempty"`
	Weight   *float64 `json:"weight,omitempty"`
	Volume   *float64 `json:"volume,omitempty"`
	VesselID *int64   `json:"vesselId,omitempty"` // была незакрытая кавычка в теге
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

		// Проверки на nil не нужны: присвоение nil в nil-поле ничего не меняет
		_, err := h.client.Update(ctx, &cargov1.UpdateRequest{
			Id:       id,
			Title:    req.Title,
			TypeId:   req.TypeID,
			Weight:   req.Weight,
			Volume:   req.Volume,
			VesselId: req.VesselID,
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

		if _, err := h.client.Delete(ctx, &cargov1.DeleteRequest{Id: id}); err != nil {
			response.GRPCError(w, r, log, err)

			return
		}

		render.JSON(w, r, dto.Ok())
	}
}
