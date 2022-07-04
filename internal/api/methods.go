package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

const (
	urlGetNewsList   = "/news"
	urlGetFullNews   = "/news/full/"
	urlComments      = "/comment/"
	urlCheckModerate = "/moderate"
)

type ReqError struct {
	err        error
	statusCode int
}

//NewsListHandler возвращает страницу со списком новостей.
//Поддерживает фильтрацию по названию новости (параметр filter)
//и запрашивемый номер страницы (параметр page).
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

	g.setReqParameters(req, filter, page)

	resp, reqErr := g.executeRequest(req)
	if reqErr != nil {
		g.writeResponseError(w, reqErr.err, reqErr.statusCode)
		return
	}

	var newsList NewsList
	if err := json.NewDecoder(resp.Body).Decode(&newsList); err != nil {
		g.writeResponseError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Add("Code", strconv.Itoa(http.StatusOK))
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(newsList)
}

//AddCommentHandler добавляет комментарий к новости но id новости.
//Проверяет комментарий на цензуру, при непрохождении проверки вернет ошибку.
func (g *Gateway) AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	newsID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		g.writeResponseError(w, fmt.Errorf("invalid parameter passed: %s", err), http.StatusBadRequest)
		return
	}

	reqModerate, err := http.NewRequestWithContext(r.Context(), http.MethodPost, "http://"+g.moderateAddress+urlCheckModerate, r.Body)
	if err != nil {
		g.writeResponseError(w, err, http.StatusInternalServerError)
		return
	}

	_, reqErr := g.executeRequest(reqModerate)
	if reqErr != nil {
		g.writeResponseError(w, reqErr.err, reqErr.statusCode)
		return
	}

	reqComment, err := http.NewRequestWithContext(r.Context(), http.MethodPost, "http://"+g.commentsAddress+urlComments+strconv.Itoa(newsID), r.Body)
	if err != nil {
		g.writeResponseError(w, err, http.StatusInternalServerError)
		return
	}

	_, reqErr = g.executeRequest(reqComment)
	if reqErr != nil {
		g.writeResponseError(w, reqErr.err, reqErr.statusCode)
		return
	}

	w.Header().Add("Code", strconv.Itoa(http.StatusNoContent))
	w.WriteHeader(http.StatusNoContent)
}

//FullNewsHandler возвращает объект полной новости по ее id и комментарии к ней.
func (g *Gateway) FullNewsHandler(w http.ResponseWriter, r *http.Request) {
	newsID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		g.writeResponseError(w, fmt.Errorf("invalid parameter passed: %s", err), http.StatusBadRequest)
		return
	}

	chanRespones := make(chan interface{}, 2)

	wg := sync.WaitGroup{}

	wg.Add(2)
	go g.fullNewsRequest(r, newsID, chanRespones, &wg)
	go g.commentRequest(r, newsID, chanRespones, &wg)

	wg.Wait()

	close(chanRespones)

	var fullNews NewsFullDetailed

	for content := range chanRespones {
		switch t := content.(type) {
		case *ReqError:
			g.writeResponseError(w, t.err, t.statusCode)
			return

		case NewsFullDetailed:
			fullNews.ID = t.ID
			fullNews.Title = t.Title
			fullNews.PubTime = t.PubTime
			fullNews.Link = t.Link
			fullNews.Content = t.Content

		case []Comment:
			fullNews.Comments = t
		}
	}

	w.Header().Add("Code", strconv.Itoa(http.StatusOK))
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(fullNews)
}

func (g *Gateway) fullNewsRequest(r *http.Request, id int, chanRespones chan<- interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, "http://"+g.newsAddress+urlGetFullNews+strconv.Itoa(id), nil)
	if err != nil {
		chanRespones <- &ReqError{
			err:        err,
			statusCode: http.StatusInternalServerError,
		}
		return
	}

	resp, reqErr := g.executeRequest(req)
	if reqErr != nil {
		chanRespones <- reqErr
		return
	}

	var news NewsFullDetailed
	if err := json.NewDecoder(resp.Body).Decode(&news); err != nil {
		chanRespones <- &ReqError{
			err:        err,
			statusCode: http.StatusInternalServerError,
		}
		return
	}

	chanRespones <- news
}

func (g *Gateway) commentRequest(r *http.Request, id int, chanRespones chan<- interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, "http://"+g.commentsAddress+urlComments+strconv.Itoa(id), nil)
	if err != nil {
		chanRespones <- &ReqError{
			err:        err,
			statusCode: http.StatusInternalServerError,
		}
		return
	}

	resp, reqErr := g.executeRequest(req)
	if reqErr != nil {
		chanRespones <- reqErr
		return
	}

	var comments []Comment
	if err := json.NewDecoder(resp.Body).Decode(&comments); err != nil {
		chanRespones <- &ReqError{
			err:        err,
			statusCode: http.StatusInternalServerError,
		}
		return
	}

	chanRespones <- comments
}

func (g *Gateway) executeRequest(req *http.Request) (*http.Response, *ReqError) {
	if err := g.setReqID(req); err != nil {
		return nil, &ReqError{
			err:        err,
			statusCode: http.StatusBadRequest,
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, &ReqError{
			err:        err,
			statusCode: http.StatusInternalServerError,
		}
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, &ReqError{
				err:        err,
				statusCode: http.StatusInternalServerError,
			}
		}

		resp.Body.Close()
		return nil, &ReqError{
			err:        fmt.Errorf(string(body)),
			statusCode: resp.StatusCode,
		}
	}

	return resp, nil
}
