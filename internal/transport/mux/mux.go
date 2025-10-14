package mux

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"sensor-api/internal/config"
	"sensor-api/internal/transport/websocket"
)

func GetMux(pathToSpec string, port string, cfg *config.Server, ws *websocket.Ws) *chi.Mux {
	r := chi.NewMux()

	if cfg.Swagger {
		r.Get("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/yaml")
			http.ServeFile(w, r, fmt.Sprintf("%s/oapi-spec.yml", pathToSpec))
		})

		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("/openapi.yaml")))
	}

	if cfg.Websocket {
		r.HandleFunc("/ws/send", ws.SendData)
		r.HandleFunc("/ws/read", ws.ReadData)
	}

	return r
}
