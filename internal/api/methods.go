package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const (
	urlGetNewsList   = "/news"
	urlGetFullNews   = "/news/full/"
	urlGetComments   = ""
	urlCheckModerate = ""
)

func (g *Gateway) NewsListHandler(w http.ResponseWriter, r *http.Request) {
	page, filter, err := g.getPageAndFilterParams(w, r)
	if err != nil {
		g.writeResponseError(w, err, http.StatusBadRequest)
		return
	}

	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, "http://"+g.newsAddress+urlGetNewsList, nil)
	if err != nil {
		g.writeResponseError(w, err, http.StatusInternalServerError)
		return
	}

	if err := g.setReqParameters(req, filter, page); err != nil {
		g.writeResponseError(w, err, http.StatusBadRequest)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		g.writeResponseError(w, err, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		g.writeResponseError(w, fmt.Errorf("error"), resp.StatusCode)
		return
	}

	var newsList NewsList
	if err := json.NewDecoder(resp.Body).Decode(&newsList); err != nil {
		g.writeResponseError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Add("Code", strconv.Itoa(http.StatusOK))
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
