package report

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/deadsnxcks/dbcp/gateway/internal/http-server/dto"
	"github.com/deadsnxcks/dbcp/gateway/internal/http-server/response"
	reportv1 "github.com/deadsnxcks/dbcp/protos/gen/go/report"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

const (
	opStart = "handlers.report"
	timeout = 15 * time.Second
)

type Handler struct {
	log    *slog.Logger
	client reportv1.ReportServiceClient
}

func New(
	log *slog.Logger,
	client reportv1.ReportServiceClient,
) *Handler {
	return &Handler{log: log, client: client}
}

func (h *Handler) logger(r *http.Request, op string) *slog.Logger {
	return h.log.With(
		slog.String("op", opStart+"."+op),
		slog.String("req_id", middleware.GetReqID(r.Context())),
	)
}

func (h *Handler) CargoDetailReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := h.logger(r, "CargoDetailReport")

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		resp, err := h.client.GenerateUnloadedCargoReport(
			ctx,
			&reportv1.UnloadedCargoReportRequest{},
		)
		if err != nil {
			response.GRPCError(w, r, log, err)

			return
		}

		render.JSON(w, r, dto.CargoDetailItemsFromProto(resp.GetItems()))
	}
}

func (h *Handler) CargoTypeReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := h.logger(r, "CargoTypeReport")

		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		resp, err := h.client.GenerateCargoTypeSummaryReport(
			ctx,
			&reportv1.CargoTypeReportRequest{},
		)
		if err != nil {
			response.GRPCError(w, r, log, err)

			return
		}

		render.JSON(w, r, dto.CargoTypeItemsFromProto(resp.GetItems()))
	}
}