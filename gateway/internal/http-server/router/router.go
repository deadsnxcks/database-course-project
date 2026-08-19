package router

import (
	"log/slog"

	"github.com/deadsnxcks/dbcp/gateway/internal/http-server/handlers/cargo"
	cargotype "github.com/deadsnxcks/dbcp/gateway/internal/http-server/handlers/cargo-type"
	"github.com/deadsnxcks/dbcp/gateway/internal/http-server/handlers/operations"
	"github.com/deadsnxcks/dbcp/gateway/internal/http-server/handlers/opercargo"
	"github.com/deadsnxcks/dbcp/gateway/internal/http-server/handlers/report"
	"github.com/deadsnxcks/dbcp/gateway/internal/http-server/handlers/storageloc"
	"github.com/deadsnxcks/dbcp/gateway/internal/http-server/handlers/vessel"
	mw "github.com/deadsnxcks/dbcp/gateway/internal/http-server/middleware"
	cargov1 "github.com/deadsnxcks/dbcp/protos/gen/go/cargo"
	cargotypev1 "github.com/deadsnxcks/dbcp/protos/gen/go/cargotype"
	operationv1 "github.com/deadsnxcks/dbcp/protos/gen/go/operation"
	opercargov1 "github.com/deadsnxcks/dbcp/protos/gen/go/opercargo"
	reportv1 "github.com/deadsnxcks/dbcp/protos/gen/go/report"
	storagelocv1 "github.com/deadsnxcks/dbcp/protos/gen/go/storageloc"
	vesselv1 "github.com/deadsnxcks/dbcp/protos/gen/go/vessel"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func New(
	log *slog.Logger,
	vesselClient vesselv1.VesselServiceClient,
	cargoClient cargov1.CargoServiceClient,
	cargoTypeClient cargotypev1.CargoTypeServiceClient,
	operationsClient operationv1.OperationServiceClient,
	operCargoClient opercargov1.OperationCargoServiceClient,
	storageLocClient storagelocv1.StorageLocationServiceClient,
	reportClient reportv1.ReportServiceClient,
) chi.Router {
	router := chi.NewRouter()

	router.Use(chimw.RequestID)
	router.Use(chimw.Recoverer)
	router.Use(mw.Logging(log))

	vesselHandler := vessel.New(log, vesselClient)
	cargoTypeHandler := cargotype.New(log, cargoTypeClient)
	cargoHandler := cargo.New(log, cargoClient)
	operationsHandler := operations.New(log, operationsClient)
	operCargoHandler := opercargo.New(log, operCargoClient)
	storageLocHandler := storageloc.New(log, storageLocClient)
	reportHandler := report.New(log, reportClient)

	router.Route("/vessel", func(r chi.Router) {
		r.Get("/", vesselHandler.List())
		r.Post("/", vesselHandler.Create())
		r.Get("/{id}", vesselHandler.Get())
		r.Put("/{id}", vesselHandler.Update())
		r.Delete("/{id}", vesselHandler.Delete())
	})

	router.Route("/cargotype", func(r chi.Router) {
		r.Get("/", cargoTypeHandler.List())
		r.Post("/", cargoTypeHandler.Create())
		r.Get("/{id}", cargoTypeHandler.Get())
		r.Put("/{id}", cargoTypeHandler.Update())
		r.Delete("/{id}", cargoTypeHandler.Delete())
	})

	router.Route("/cargo", func(r chi.Router) {
		r.Get("/", cargoHandler.List())
		r.Post("/", cargoHandler.Create())
		r.Get("/{id}", cargoHandler.Get())
		r.Put("/{id}", cargoHandler.Update())
		r.Delete("/{id}", cargoHandler.Delete())
	})

	router.Route("/operation", func(r chi.Router) {
		r.Get("/", operationsHandler.List())
		r.Post("/", operationsHandler.Create())
		r.Get("/{id}", operationsHandler.Get())
		r.Put("/{id}", operationsHandler.Update())
		r.Delete("/{id}", operationsHandler.Delete())
	})

	router.Route("/opercargo", func(r chi.Router) {
		r.Get("/", operCargoHandler.List())
		r.Post("/", operCargoHandler.Create())
		r.Delete("/", operCargoHandler.Delete())
	})

	router.Route("/storageloc", func(r chi.Router) {
		r.Get("/", storageLocHandler.List())
		r.Post("/", storageLocHandler.Create())
		r.Get("/{id}", storageLocHandler.Get())
		r.Put("/{id}", storageLocHandler.Update())
		r.Delete("/{id}", storageLocHandler.Delete())
		r.Post("/{id}/use", storageLocHandler.Use())
		r.Post("/{id}/reset", storageLocHandler.Reset())
	})

	router.Method("Get", "/report-cargo-detail", reportHandler.CargoDetailReport())
	router.Method("Get", "/report-cargo-type", reportHandler.CargoTypeReport())

	return router
}
