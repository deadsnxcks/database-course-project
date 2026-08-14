package opercargo

import (
	"context"
	"github.com/deadsnxcks/dbcp/gateway/internal/lib/api/response"
	"github.com/deadsnxcks/dbcp/gateway/internal/lib/logger/sl"
	opercargov1 "github.com/deadsnxcks/dbcp/protos/gen/go/opercargo"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const opStart = "handlers.opercargo"

type Handler struct {
	log    *slog.Logger
	client opercargov1.OperationCargoServiceClient
}

func New(
	log *slog.Logger,
	client opercargov1.OperationCargoServiceClient,
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

		resp, err := h.client.List(ctx, &opercargov1.ListRequest{}) 
		if err != nil {
			log.Error("grpc call failed", sl.Err(err))
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		result := []map[string]interface{}{}
		for _, v := range resp.GetOperationsCargos() {
			result = append(result, map[string]interface{}{
				"operationId":    	v.GetOperationId(),
				"cargoId": 			v.GetCargoId(),
			})
		}

		render.JSON(w, r, result)
	}
}

type createRequest struct {
	OperationId	int64	`json:"operationId"`
	CargoId		int64	`json:"cargoId"`
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
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, response.Error("invalid request body"))
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		grpcReq := &opercargov1.CreateRequest{
			OperationId: req.OperationId,
			CargoId:     req.CargoId,
		}
		_, err := h.client.Create(ctx, grpcReq)
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

type DeleteRequest struct {
	OperationId	int64	`json:"operationId"`
	CargoId		int64	`json:"cargoId"`
}

func (h *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = opStart + ".Delete"
		
		log := h.log.With(
			slog.String("op", op),
			slog.String("req_id", middleware.GetReqID(r.Context())),
		)

		var req DeleteRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, response.Error("invalid request body"))
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		grpcReq := &opercargov1.DeleteRequest{
			OperationId: req.OperationId,
			CargoId:     req.CargoId,
		}
		_, err := h.client.Delete(ctx, grpcReq)
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

		render.JSON(w, r, map[string]string{"status": "ok"})
	}
}