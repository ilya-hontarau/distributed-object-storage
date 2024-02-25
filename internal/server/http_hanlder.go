package server

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"

	"github.com/ilya-hontarau/distributed-object-storage/internal/gateway"
)

type HTTPHandler struct {
	http.Handler
	gateway *gateway.Gateway
	logger  *slog.Logger
}

func NewHTTPHandler(gateway *gateway.Gateway, logger *slog.Logger) *HTTPHandler {
	h := &HTTPHandler{}
	h.gateway = gateway
	h.logger = logger
	mux := http.NewServeMux()
	mux.HandleFunc("PUT /object/{id}", h.UploadFile)
	h.Handler = mux
	return h
}

func (h *HTTPHandler) UploadFile(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	// TODO: validate id
	// TODO: check content length
	// TODO: set max value
	bodyBytes, err := io.ReadAll(http.MaxBytesReader(w, req.Body, 10000))
	if err != nil {
		h.logger.Debug("Request is too big")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.gateway.Upload(req.Context(), id, bytes.NewReader(bodyBytes), len(bodyBytes))
	if err != nil {
		h.logger.Error("Failed to upload", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.logger.Debug("Successfully uploaded", slog.String("id", id))
	w.WriteHeader(http.StatusOK)
}
