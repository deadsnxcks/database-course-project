package vessel

import (
	"context"
	"github.com/deadsnxcks/dbcp/gateway/internal/lib/api/response"
	"github.com/deadsnxcks/dbcp/gateway/internal/lib/logger/sl"
	vesselv1 "github.com/deadsnxcks/dbcp/protos/gen/go/vessel"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/grpc/status"
    "google.golang.org/grpc/codes"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Handler struct {
	log    *slog.Logger
	client vesselv1.VesselServiceClient
}

func New(log *slog.Logger, client vesselv1.VesselServiceClient) *Handler {
	return &Handler{
		log:    log,
		client: client,
	}
}

func (h *Handler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.vessel.List"
		log := h.log.With(
			slog.String("op", op),
			slog.String("req_id", middleware.GetReqID(r.Context())),
		)

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		resp, err := h.client.List(ctx, &vesselv1.ListRequest{})
		if err != nil {
			log.Error("grpc call failed", sl.Err(err))
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		result := []map[string]interface{}{}
		for _, v := range resp.GetVessels() {
			result = append(result, map[string]interface{}{
				"id":          v.GetId(),
				"title":       v.GetTitle(),
				"vesselType": v.GetVesselType(),
				"maxLoad":    v.GetMaxLoad(),
			})
		}

		render.JSON(w, r, result)
	}
}

type createRequest struct {
	Title      string  `json:"title"`
	VesselType string  `json:"vesselType"`
	MaxLoad    float64 `json:"maxLoad"`
}

func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.vessel.Create"
		log := h.log.With(
			slog.String("op", op),
			slog.String("req_id", middleware.GetReqID(r.Context())),
		)

		var req createRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode body", sl.Err(err))
			render.JSON(w, r, response.Error("invalid request"))
			return
		}

		log.Debug("received request",
			slog.String("title", req.Title),
			slog.String("vesselType", req.VesselType),
			slog.Float64("maxLoad", req.MaxLoad))

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		grpcReq := &vesselv1.CreateRequest{
			Title:      req.Title,
			VesselType: req.VesselType,
			MaxLoad:    req.MaxLoad,
		}

		resp, err := h.client.Create(ctx, grpcReq)
		if err != nil {
			log.Error("grpc call failed", sl.Err(err))
			st, ok := status.FromError(err)
			if !ok {
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, response.Error("internal error"))
				return
			}

			httpStatus := http.StatusInternalServerError

			switch st.Code() {
			case codes.AlreadyExists:
				httpStatus = http.StatusConflict
			case codes.NotFound:
				httpStatus = http.StatusNotFound
			case codes.InvalidArgument:
				httpStatus = http.StatusBadRequest
			}

			render.Status(r, httpStatus)
			render.JSON(w, r, map[string]string{
				"code":    st.Code().String(),
				"message": st.Message(),
			})
			
			return
		}

		render.JSON(w, r, map[string]interface{}{
			"id": resp.GetId(),
		})
	}
}

func (h *Handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.vessel.Get"
		log := h.log.With(
			slog.String("op", op),
			slog.String("req_id", middleware.GetReqID(r.Context())),
		)

		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			render.JSON(w, r, response.Error("invalid id"))
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		resp, err := h.client.Get(ctx, &vesselv1.GetRequest{Id: id})
		if err != nil {
			log.Error("grpc call failed", sl.Err(err))
			st, ok := status.FromError(err)
			if !ok {
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, response.Error("internal error"))
				return
			}

			httpStatus := http.StatusInternalServerError

			switch st.Code() {
			case codes.AlreadyExists:
				httpStatus = http.StatusConflict
			case codes.NotFound:
				httpStatus = http.StatusNotFound
			case codes.InvalidArgument:
				httpStatus = http.StatusBadRequest
			}

			render.Status(r, httpStatus)
			render.JSON(w, r, map[string]string{
				"code":    st.Code().String(),
				"message": st.Message(),
			})
			
			return
		}

		v := resp.GetVessel()
		render.JSON(w, r, map[string]interface{}{
			"id":          v.GetId(),
			"title":       v.GetTitle(),
			"vessel_type": v.GetVesselType(),
			"max_load":    v.GetMaxLoad(),
		})
	}
}

type updateRequest struct {
	Title      *string  `json:"title,omitempty"`
	VesselType *string  `json:"vesselType,omitempty"`
	MaxLoad    *float64 `json:"maxLoad,omitempty"`
}

func (h *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.vessel.Update"
		log := h.log.With(
			slog.String("op", op),
			slog.String("req_id", middleware.GetReqID(r.Context())),
		)

		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			render.JSON(w, r, response.Error("invalid id"))
			return
		}

		var req updateRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode body", sl.Err(err))
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		grpcReq := &vesselv1.UpdateRequest{
			Id:         id,
			Title:      req.Title,
			VesselType: req.VesselType,
		}
		if req.MaxLoad != nil {
			grpcReq.MaxLoad = req.MaxLoad
		}

		_, err = h.client.Update(ctx, grpcReq)
		if err != nil {
			log.Error("grpc call failed", sl.Err(err))
			st, ok := status.FromError(err)
			if !ok {
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, response.Error("internal error"))
				return
			}

			httpStatus := http.StatusInternalServerError

			switch st.Code() {
			case codes.AlreadyExists:
				httpStatus = http.StatusConflict
			case codes.NotFound:
				httpStatus = http.StatusNotFound
			case codes.InvalidArgument:
				httpStatus = http.StatusBadRequest
			}

			render.Status(r, httpStatus)
			render.JSON(w, r, map[string]string{
				"code":    st.Code().String(),
				"message": st.Message(),
			})
			
			return
		}

		render.JSON(w, r, map[string]string{"status": "ok"})
	}
}

func (h *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.vessel.Delete"
		log := h.log.With(
			slog.String("op", op),
			slog.String("req_id", middleware.GetReqID(r.Context())),
		)

		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Error("invalid vessel id", sl.Err(err))
			render.JSON(w, r, response.Error("invalid id"))
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		_, err = h.client.Delete(ctx, &vesselv1.DeleteRequest{Id: id})
		if err != nil {
			log.Error("grpc call failed", sl.Err(err))
			st, ok := status.FromError(err)
			if !ok {
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, response.Error("internal error"))
				return
			}

			httpStatus := http.StatusInternalServerError

			switch st.Code() {
			case codes.FailedPrecondition:
				httpStatus = http.StatusConflict
			case codes.NotFound:
				httpStatus = http.StatusNotFound
			case codes.InvalidArgument:
				httpStatus = http.StatusBadRequest
			}

			render.Status(r, httpStatus)
			render.JSON(w, r, map[string]string{
				"code":    st.Code().String(),
				"message": st.Message(),
			})
			
			return
		}

		render.JSON(w, r, map[string]string{"status": "ok"})
	}
}