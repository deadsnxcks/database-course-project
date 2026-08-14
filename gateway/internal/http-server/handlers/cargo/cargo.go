package cargo

import (
	"context"
	"github.com/deadsnxcks/dbcp/gateway/internal/lib/api/response"
	"github.com/deadsnxcks/dbcp/gateway/internal/lib/logger/sl"
	cargov1 "github.com/deadsnxcks/dbcp/protos/gen/go/cargo"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const opStart = "handlers.cargo"

type Handler struct {
	log *slog.Logger
	client cargov1.CargoServiceClient
}

func New(
	log *slog.Logger, 
	client cargov1.CargoServiceClient,
) *Handler {
	return &Handler{
		log: log,
		client: client,
	}
}

func (h *Handler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = opStart + ".List"

		log := h.log.With(
			slog.String("op", op),
			slog.String("req_id", middleware.GetReqID(r.Context())),
		)

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		resp, err := h.client.List(ctx, &cargov1.ListRequest{}) 
		if err != nil {
			log.Error("grpc call failed", sl.Err(err))
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		result := []map[string]interface{}{}
		for _, v := range resp.GetCargos() {
			result = append(result, map[string]interface{}{
				"id": v.GetId(),
				"title": v.GetTitle(),
				"typeId": v.GetTypeId(),
				"weight": v.GetWeight(),
				"volume": v.GetVolume(),
				"vesselId": v.GetVesselId(), 
			})
		}

		render.JSON(w, r, result)
	}
}

type createRequest struct {
	Title 		string	`json:"title"`
	TypeID 		int64	`json:"typeId"`
	Weight 		float64	`json:"weight"`
	Volume 		float64	`json:"volume"`
	VesselID 	int64	`json:"vesselId"`
}

func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = opStart + ".Create"

		log := h.log.With(
			slog.String("op", op),
			slog.String("req_id", middleware.GetReqID(r.Context())),
		)

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var req createRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode body", sl.Err(err))
			render.JSON(w, r, response.Error("invalid request"))
			return
		}

		grpcReq := cargov1.CreateRequest{
			Title: req.Title,
			TypeId: req.TypeID,
			Weight: req.Weight,
			Volume: req.Volume,
			VesselId: req.VesselID,
		}
		resp, err := h.client.Create(ctx, &grpcReq)
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
		const op = opStart + ".Get"

		log := h.log.With(
			slog.String("op", op),
			slog.String("req_id", middleware.GetReqID(r.Context())),
		)

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			render.JSON(w, r, response.Error("invalid id"))
			return
		}

		resp, err := h.client.Get(ctx, &cargov1.GetRequest{Id: id})
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

		c := resp.GetCargo()
		render.JSON(w, r, map[string]interface{}{
			"id": c.GetId(),
			"title": c.GetTitle(),
			"typeId": c.GetTypeId(),
			"weight": c.GetWeight(),
			"volume": c.GetVolume(),
			"vesselId": c.GetVesselId(),
		})
	}
}

type updateRequest struct {
	Title 		*string		`json:"title"`
	TypeID 		*int64		`json:"typeId"`
	Weight 		*float64	`json:"weight"`
	Volume 		*float64	`json:"volume"`
	VesselID 	*int64		`json:"vesselId`
}

func (h *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = opStart + ".Update"

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
			render.JSON(w, r, response.Error("invalid request"))
			return
		}

		grpcReq := &cargov1.UpdateRequest{
			Id: id,
		}
		
		if req.Title != nil {
			grpcReq.Title = req.Title
		}
		if req.TypeID != nil {
			grpcReq.TypeId = req.TypeID
		}
		if req.Weight != nil {
			grpcReq.Weight = req.Weight
		}
		if req.Volume != nil {
			grpcReq.Volume = req.Volume
		}
		if req.VesselID != nil {
			grpcReq.VesselId = req.VesselID
		}

		_, err = h.client.Update(r.Context(), grpcReq)
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
		const op = opStart + ".Delete"

		log := h.log.With(
			slog.String("op", op),
			slog.String("req_id", middleware.GetReqID(r.Context())),
		)

		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Error("invalid cargo id", sl.Err(err))
			render.JSON(w, r, response.Error("invalid id"))
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		_, err = h.client.Delete(ctx, &cargov1.DeleteRequest{Id: id})
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