package server

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ilya-hontarau/distributed-object-storage/internal/gateway"
)

const maxSize = 1024 * 1024 * 50

type HTTPHandler struct {
	http.Handler
	gateway *gateway.Gateway
	logger  *slog.Logger
}

func NewHTTPHandler(gateway *gateway.Gateway, logger *slog.Logger) *HTTPHandler {
	h := &HTTPHandler{
		Handler: nil,
		gateway: gateway,
		logger:  logger,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("PUT /object/{id}", h.UploadFile)
	mux.HandleFunc("GET /object/{id}", h.DownloadFile)
	h.Handler = mux
	return h
}

func isValidID(id string) bool {
	return len(id) > 1 && len(id) <= 32 && !strings.Contains(id, " ")
}

func (h *HTTPHandler) UploadFile(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	if !isValidID(id) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bodyBytes, err := io.ReadAll(http.MaxBytesReader(w, req.Body, maxSize))
	if err != nil {
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
	if !isValidID(id) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file, err := h.gateway.Download(req.Context(), id)
	if err != nil {
		if errors.Is(err, gateway.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
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
