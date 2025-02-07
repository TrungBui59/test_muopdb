package http

import (
	"github.com/TrungBui59/test_muopdb/interval/configs"
	"github.com/TrungBui59/test_muopdb/interval/muopdbclient"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type App struct {
	muopDBClient muopdbclient.MuopDbClient
	cfg          configs.Config
}

func (app App) routes() http.Handler {
	// set up router
	mux := chi.NewRouter()

	// specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))
	return mux
}
