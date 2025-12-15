package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/merev/ds-player-api/internal/player"
)

func NewRouter(ph *player.Handler) http.Handler {
	r := chi.NewRouter()

	r.Route("/api", func(api chi.Router) {
		api.Get("/players", ph.ListPlayers)
		api.Post("/players", ph.CreatePlayer)
	})

	return r
}
