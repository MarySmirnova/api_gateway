package api

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MarySmirnova/api_gateway/internal/config"
	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

type ContextKey string

const ContextReqIDKey ContextKey = "request_id"

type Gateway struct {
	httpServer      *http.Server
	newsAddress     string
	commentsAddress string
	moderateAddress string
}

func NewGateway(cfg config.Server) *Gateway {
	g := &Gateway{
		newsAddress:     cfg.NewsAddress,
		commentsAddress: cfg.CommentsAddress,
		moderateAddress: cfg.ModerateAddress,
	}

	handler := mux.NewRouter()
	handler.Use(g.reqIDMiddleware, g.logMiddleware)
	handler.Name("get_news_list").Methods(http.MethodGet).Path("/news").HandlerFunc(g.NewsListHandler)
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

func (g *Gateway) GetHTTPServer() *http.Server {
	return g.httpServer
}

func (g *Gateway) reqIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqID int
		reqIDString := r.FormValue("request_id")

		if reqIDString == "" {
			reqID = g.generateReqID()
		}

		if reqIDString != "" {
			id, err := strconv.Atoi(reqIDString)
			if err != nil {
				g.writeResponseError(w, err, http.StatusBadRequest)
				return
			}
			reqID = id
		}

		ctx := context.WithValue(r.Context(), ContextReqIDKey, reqID)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (g *Gateway) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			log.WithFields(log.Fields{
				"request_time": time.Now().Format("2006-01-02 15:04:05.000000"),
				"request_ip":   strings.TrimPrefix(strings.Split(r.RemoteAddr, ":")[1], "["),
				"code":         w.Header().Get("Code"),
				"request_id":   r.Context().Value(ContextReqIDKey),
			}).Info("news reader response")
		}()

		next.ServeHTTP(w, r)
	})
}

func (g *Gateway) generateReqID() int {
	max := 999999999999
	min := 100000

	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func (g *Gateway) getPageAndFilterParams(w http.ResponseWriter, r *http.Request) (int, string, error) {
	var page int
	filter := r.FormValue("filter")
	pageString := r.FormValue("page")
	if pageString == "" {
		page = 1
	}
	if pageString != "" {
		p, err := strconv.Atoi(pageString)
		if err != nil {
			return 0, "", err
		}
		page = p
	}

	return page, filter, nil
}

func (g *Gateway) writeResponseError(w http.ResponseWriter, err error, code int) {
	w.Header().Add("Code", strconv.Itoa(code))
	log.WithError(err).Error("api error")
	w.WriteHeader(code)
	_, _ = w.Write([]byte(err.Error()))
}

func (g *Gateway) setReqParameters(r *http.Request, filter string, page int) {
	q := r.URL.Query()
	q.Add("filter", filter)
	q.Add("page", strconv.Itoa(page))
	r.URL.RawQuery = q.Encode()
}

func (g *Gateway) setReqID(r *http.Request) error {
	var reqID string

	switch id := r.Context().Value(ContextReqIDKey).(type) {
	case string:
		reqID = string(id)
	case int:
		reqID = strconv.Itoa(id)
	default:
		fmt.Printf("id: %v, type: %T", id, id)
		return fmt.Errorf("неизветсный тип данных параметра request_id")

	}
	q := r.URL.Query()
	q.Add("request_id", reqID)
	r.URL.RawQuery = q.Encode()

	return nil
}
