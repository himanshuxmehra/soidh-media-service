package api

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"soidh-media-service/internal/api/handlers"
	custommiddleware "soidh-media-service/internal/middleware"
)

func SetupRoutes(db *sql.DB, logger *zap.Logger) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(custommiddleware.ZapLogger(logger))

	h := handlers.New(db, logger)

	r.Get("/", h.Home)
	r.Post("/upload/{accountId}/{folderId}", h.UploadFile)

	return r
}
