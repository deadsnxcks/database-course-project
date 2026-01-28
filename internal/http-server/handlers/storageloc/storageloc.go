package storageloc

import (
	"context"
	"dbcp-api-gateway/internal/lib/api/response"
	"dbcp-api-gateway/internal/lib/logger/sl"
	storagelocv1 "dbcp-api-gateway/protos/gen/go/storageloc"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const opStart = "handlers.storageloc"

type Handler struct {
	log    *slog.Logger
	client storagelocv1.StorageLocationServiceClient
}

func New(
	log *slog.Logger,
	client storagelocv1.StorageLocationServiceClient,
) *Handler {
	return &Handler{
		log:    log,
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
		
		resp, err := h.client.List(ctx, &storagelocv1.ListRequest{})
		if err != nil {
			log.Error("grpc call failed", sl.Err(err))
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		var result []map[string]interface{}
		for _, v := range resp.GetStorageLocations() {
			result = append(result, map[string]interface{}{
				"id":    			v.GetId(),
				"cargoTypeId": 		v.GetCargoTypeId(),
				"maxWeight": 		v.GetMaxWeight(),
				"maxVolume": 		v.GetMaxVolume(),
				"cargoId": 			v.GetCargoId(),
				"dateOfPlacement": 	v.GetDateOfPlacement(),
			})
		}

		render.JSON(w, r, result)
	}
}

type createRequest struct {
	CargoTypeId	int64	`json:"cargoTypeId"`
	MaxWeight	float64	`json:"maxWeight"`
	MaxVolume	float64	`json:"maxVolume"`
}

func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = opStart + ".Create"

		log := h.log.With(
			slog.String("op", op),
			slog.String("req_id", middleware.GetReqID(r.Context())),
		)
		
		var req createRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode json", sl.Err(err))
			render.JSON(w, r, response.Error("invalid request body"))
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		_, err := h.client.Create(ctx, &storagelocv1.CreateRequest{
			CargoTypeId:	req.CargoTypeId,
			MaxWeight:		req.MaxWeight,
			MaxVolume:		req.MaxVolume,
		})
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
			log.Error("invalid id param", sl.Err(err))
			render.JSON(w, r, response.Error("invalid id param"))
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		_, err = h.client.Delete(ctx, &storagelocv1.DeleteRequest{Id: id})
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

func (h *Handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = opStart + ".Get"

		log := h.log.With(
			slog.String("op", op),
			slog.String("req_id", middleware.GetReqID(r.Context())),
		)
		
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Error("invalid id param", sl.Err(err))
			render.JSON(w, r, response.Error("invalid id param"))
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		resp, err := h.client.Get(ctx, &storagelocv1.GetRequest{Id: id})
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

		sl := resp.GetStorageLocation()
		render.JSON(w, r, map[string]interface{}{
			"id":    			sl.GetId(),
			"cargoTypeId": 		sl.GetCargoTypeId(),
			"maxWeight": 		sl.GetMaxWeight(),
			"maxVolume": 		sl.GetMaxVolume(),
			"cargoId": 			sl.GetCargoId(),
			"dateOfPlacement": 	sl.GetDateOfPlacement(),
		})	
	}
}

type updateRequest struct {
	CargoTypeId	*int64	`json:"cargoTypeId"`
	MaxWeight	*float64	`json:"maxWeight"`
	MaxVolume	*float64	`json:"maxVolume"`
}

func (h *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = opStart + ".Update"
		
		log := h.log.With(
			slog.String("op", op),
			slog.String("req_id", middleware.GetReqID(r.Context())),
		)

		var req updateRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode json", sl.Err(err))
			render.JSON(w, r, response.Error("invalid request body"))
			return
		}

		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Error("invalid id param", sl.Err(err))
			render.JSON(w, r, response.Error("invalid id param"))
			return
		}
		
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		_, err = h.client.Update(ctx, &storagelocv1.UpdateRequest{
			Id:				id,
			CargoTypeId:	req.CargoTypeId,
			MaxWeight:		req.MaxWeight,
			MaxVolume:		req.MaxVolume,
		})
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

type useRequest struct {
	CargoId			int64		`json:"cargoId"`
	DateOfPlacement	time.Time	`json:"dateOfPlacement"`
}

func (h *Handler) Use() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = opStart + ".Use"

		log := h.log.With(
			slog.String("op", op),
			slog.String("req_id", middleware.GetReqID(r.Context())),
		)

		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Error("invalid id param", sl.Err(err))
			render.JSON(w, r, response.Error("invalid id param"))
			return
		}

		var req useRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode json", sl.Err(err))
			render.JSON(w, r, response.Error("invalid request body"))
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		dateOfPlacement := timestamppb.New(req.DateOfPlacement)
		_, err = h.client.Use(ctx, &storagelocv1.UseRequest{
			StorageLocationId:	id,
			CargoId:			req.CargoId,
			DateOfPlacement:	dateOfPlacement,
		})
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

func (h *Handler) Reset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = opStart + ".Reset"

		log := h.log.With(
			slog.String("op", op),
			slog.String("req_id", middleware.GetReqID(r.Context())),
		)

		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Error("invalid id param", sl.Err(err))
			render.JSON(w, r, response.Error("invalid id param"))
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		_, err = h.client.Reset(ctx, &storagelocv1.ResetRequest{Id: id})
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