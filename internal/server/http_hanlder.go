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
	mux.HandleFunc("GET /object/{id}", h.DownloadFile)
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

func (h *HTTPHandler) DownloadFile(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	// TODO: validate id

	file, err := h.gateway.Download(req.Context(), id)
	// TODO: return not found
	if err != nil {
		h.logger.Error("Failed to download", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = io.Copy(w, file)
	if err != nil {
		h.logger.Error("Failed to copy file to writer", slog.String("error", err.Error()))
		return
	}
	h.logger.Debug("Successfully downloaded", slog.String("id", id))
}
