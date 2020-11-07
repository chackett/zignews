package mobileapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/chackett/zignews/pkg/storage"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

const apiPrefixWithVersion = "/api/v1"

// Handler implements HTTP API
type Handler struct {
	service Service
	server  *http.Server
}

// Service defines the functionality required by the mobile-api
type Service interface {
	GetArticles(ctx context.Context, offset, count int, category, provider []string) ([]storage.Article, error)
	SaveProvider(ctx context.Context, provider storage.Provider) (string, error)
}

// ErrorResponse is returned to requests that result in some error state such as bad request or internal server error
type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}

// NewHandler constructs a new instance of Handler
func NewHandler(service Service, addr string) (Handler, error) {
	if service == nil {
		return Handler{}, errors.New("service is nil")
	}

	h := Handler{
		service: service,
	}

	r := mux.NewRouter().PathPrefix(apiPrefixWithVersion).Subrouter()
	// r.HandleFunc(fmt.Sprintf("%s/article", apiPrefixWithVersion), h.HandleGetArticles()).Methods(http.MethodGet)
	// r.HandleFunc(fmt.Sprintf("%s/provider", apiPrefixWithVersion), h.HandlePostProvider()).Methods(http.MethodPost)
	// r.HandleFunc(fmt.Sprintf("%s/ping", apiPrefixWithVersion), h.HandlePing()).Methods(http.MethodGet)
	r.HandleFunc("/article", h.HandleGetArticles()).Methods(http.MethodGet)
	r.HandleFunc("/provider", h.HandlePostProvider()).Methods(http.MethodPost)
	r.HandleFunc("/ping", h.HandlePing()).Methods(http.MethodGet)
	r.StrictSlash(true) // Saves being a nuisance

	h.server = &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return h, nil
}

// Start starts the underlying http server
func (h *Handler) Start() error {
	err := h.server.ListenAndServe()
	if err != nil {
		fmt.Println(err.Error())
	}
	return nil
}

// Stop stops the underlying http server
func (h *Handler) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*120))
	defer cancel()
	err := h.server.Shutdown(ctx)
	if err != nil {
		return errors.Wrap(err, "server Shutdown()")
	}
	return nil
}

func (h *Handler) returnError(err error, statusCode int, w http.ResponseWriter) {
	response := ErrorResponse{
		Error: err.Error(),
	}
	w.WriteHeader(statusCode)
	encodeErr := json.NewEncoder(w).Encode(response)
	if encodeErr != nil {
		log.Printf("ERROR: Unable to return error message to client. %s", err.Error())

		// Gracefully fail to respond with structured error
		w.Write([]byte(fmt.Sprintf("Error: %s", encodeErr.Error())))
	}
}

// HandleGetArticles returns articles matching specified criteria
func (h *Handler) HandleGetArticles() http.HandlerFunc {
	type ArticleResponse struct {
		Articles []storage.Article `json:"articles,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Offset and count used to satisfy the "scrollable list" via pagination.
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			log.Print("use default offset value")
			offset = defaultOffset
		}

		count, err := strconv.Atoi(r.URL.Query().Get("count"))
		if err != nil {
			log.Print("use default count value")
			count = defaultPageSize
		}

		// Filters - supports querystring format ?category=politics&category=technology&provider=msn&provider=bbc
		// These could be further sanitised
		categories := r.URL.Query()["category"]
		providers := r.URL.Query()["provider"]

		articles, err := h.service.GetArticles(r.Context(), offset, count, categories, providers)
		if err != nil {
			h.returnError(errors.Wrap(err, "get articles"), http.StatusInternalServerError, w)
			return
		}

		response := ArticleResponse{
			Articles: articles,
		}
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("ERROR: encoding response to client: %s", err.Error())
		}
	}
}

// HandlePostProvider saves provider information for use by the aggregator component
func (h *Handler) HandlePostProvider() http.HandlerFunc {
	type Request struct {
		Type                 string `json:"type,omitempty"`
		Label                string `json:"label,omitempty"`
		FeedURL              string `json:"feedURL,omitempty"`
		PollFrequencySeconds int    `json:"pollFrequencySeconds,omitempty"`
	}

	type Response struct {
		InsertedProviderID string `json:"insertedProviderID,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var requestObj Request
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&requestObj)
		if err != nil {
			h.returnError(errors.Wrap(err, "parse request"), http.StatusBadRequest, w)
			return
		}

		provider := storage.Provider{
			Label:                requestObj.Label,
			FeedURL:              requestObj.FeedURL,
			PollFrequencySeconds: requestObj.PollFrequencySeconds,
			Type:                 requestObj.Type,
		}

		providerID, err := h.service.SaveProvider(r.Context(), provider)
		if err != nil {
			h.returnError(errors.Wrap(err, "save provider"), http.StatusInternalServerError, w)
			return
		}

		response := Response{
			InsertedProviderID: providerID,
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("ERROR: encoding response to client: %s", err.Error())
		}
	}
}

// HandlePing is a simple handle to enable clients to test connectivity
func (h *Handler) HandlePing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("It works!"))
	}
}
