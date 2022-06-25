package api

import (
	"net/http"

	"github.com/MarySmirnova/api_gateway/internal/config"
	"github.com/gorilla/mux"
)

type Gateway struct {
	httpServer *http.Server
}

func NewGateway(cfg config.Server) *Gateway {
	g := &Gateway{}

	handler := mux.NewRouter()
	handler.Name("get_news_list").Methods(http.MethodGet).Path("/news").HandlerFunc(g.NewsListHandler)
	handler.Name("news_filter").Methods(http.MethodGet).Path("/news/filter").HandlerFunc(g.NewsFilterHandler)
	handler.Name("get_full_news").Methods(http.MethodGet).Path("/news/{id}").HandlerFunc(g.FullNewsHandler)
	handler.Name("add_commet").Methods(http.MethodPost).Path("/news/{id}/comment").HandlerFunc(g.AddCommentHandler)

	g.httpServer = &http.Server{
		Addr:         cfg.Listen,
		Handler:      handler,
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
	}

	return g
}
