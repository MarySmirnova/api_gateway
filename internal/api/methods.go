package api

import (
	"encoding/json"
	"net/http"
)

func (g *Gateway) NewsListHandler(w http.ResponseWriter, r *http.Request) {
	newsList := []NewsShortDetailed{
		{
			ID:    1,
			Title: "Title",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(newsList)
}

func (g *Gateway) NewsFilterHandler(w http.ResponseWriter, r *http.Request) {
	newsList := []NewsShortDetailed{
		{
			ID:    1,
			Title: "Title",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(newsList)

}

func (g *Gateway) FullNewsHandler(w http.ResponseWriter, r *http.Request) {
	fullNews := NewsFullDetailed{
		ID:      1,
		Title:   "Title",
		Content: "Content",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(fullNews)
}

func (g *Gateway) AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
