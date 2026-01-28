package report

import (
	"context"
	"dbcp-api-gateway/internal/lib/api/response"
	"dbcp-api-gateway/internal/lib/logger/sl"
	reportv1 "dbcp-api-gateway/protos/gen/go/report"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Handler struct {
	log *slog.Logger
	client reportv1.ReportServiceClient
}

func New(
	log *slog.Logger,
	client reportv1.ReportServiceClient,
) *Handler {
	return &Handler{
		log: log,
		client: client,
	}
}

func (h *Handler) CargoDetailReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.report.CargoDetailReport"

		log := h.log.With(
			slog.String("op", op),
			slog.String("req_id", middleware.GetReqID(r.Context())),
		)

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		
		resp, err := h.client.GenerateUnloadedCargoReport(ctx, &reportv1.UnloadedCargoReportRequest{})
		if err != nil {
			log.Error("grpc call failed", sl.Err(err))
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		var result []map[string]interface{}
		for _, item := range resp.GetItems() {
			result = append(result, map[string]interface{}{
				"cargoName": item.GetCargoName(),
				"weight": item.GetWeightTons(),
				"cargoType": item.GetCargoType(),
				"vesselName": item.GetVesselName(),
				"unloadDate": item.GetUnloadingDate(), 
			})
		}

		render.JSON(w, r, result)
	}
}

func (h *Handler) CargoTypeReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.report.CargoTypeReport"

		log := h.log.With(
			slog.String("op", op),
			slog.String("req_id", middleware.GetReqID(r.Context())),
		)

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		
		resp, err := h.client.GenerateCargoTypeSummaryReport(ctx, &reportv1.CargoTypeReportRequest{})
		if err != nil {
			log.Error("grpc call failed", sl.Err(err))
			render.JSON(w, r, response.Error(err.Error()))
			return
		}

		var result []map[string]interface{}
		for _, item := range resp.GetItems() {
			result = append(result, map[string]interface{}{
				"cargoTypeName": item.GetCargoTypeName(),
				"count": item.GetCargoCount(),
				"weight": item.GetTotalWeightTons(),
				"volume": item.GetTotalVolumeM3(),
				"processCost": item.GetProcessCost(),
			})
		}

		render.JSON(w, r, result)
	}
}