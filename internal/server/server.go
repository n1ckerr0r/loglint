package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/n1ckerr0r/loglint/internal/platform/storage"
	"github.com/n1ckerr0r/loglint/internal/service"
)

type Server struct {
	repo  service.RunRepository
	queue service.Queue
	mux   *http.ServeMux
}

func New(repo service.RunRepository, queue service.Queue) *Server {
	server := &Server{
		repo:  repo,
		queue: queue,
		mux:   http.NewServeMux(),
	}
	server.routes()
	return server
}

func (server *Server) Handler() http.Handler {
	return server.mux
}

func (server *Server) routes() {
	server.mux.HandleFunc("GET /healthz", server.health)
	server.mux.HandleFunc("POST /api/v1/runs", server.createRun)
	server.mux.HandleFunc("GET /api/v1/runs/", server.getRun)
}

func (server *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (server *Server) createRun(w http.ResponseWriter, r *http.Request) {
	var request service.RunRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	request.Project = strings.TrimSpace(request.Project)
	if request.Project == "" {
		writeError(w, http.StatusBadRequest, "project is required")
		return
	}

	now := time.Now().UTC()
	run := service.Run{
		ID:        newID(),
		Project:   request.Project,
		Ref:       strings.TrimSpace(request.Ref),
		Commit:    strings.TrimSpace(request.Commit),
		Status:    service.RunStatusQueued,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := server.repo.CreateRun(r.Context(), run); err != nil {
		writeError(w, http.StatusInternalServerError, "create run failed")
		return
	}
	if err := server.queue.Publish(r.Context(), service.Job{
		RunID:   run.ID,
		Project: run.Project,
		Ref:     run.Ref,
		Commit:  run.Commit,
	}); err != nil {
		_ = server.repo.UpdateRun(r.Context(), run.ID, service.RunStatusFailed, err.Error())
		writeError(w, http.StatusBadGateway, "queue publish failed")
		return
	}

	writeJSON(w, http.StatusAccepted, run)
}

func (server *Server) getRun(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/runs/")
	if id == "" || strings.Contains(id, "/") {
		writeError(w, http.StatusNotFound, "run not found")
		return
	}
	run, err := server.repo.GetRun(r.Context(), id)
	if errors.Is(err, storage.ErrRunNotFound) {
		writeError(w, http.StatusNotFound, "run not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "get run failed")
		return
	}
	writeJSON(w, http.StatusOK, run)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func newID() string {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return hex.EncodeToString([]byte(time.Now().UTC().Format(time.RFC3339Nano)))
	}
	return hex.EncodeToString(bytes[:])
}

func ListenAndServe(ctx context.Context, addr string, handler http.Handler, shutdownTimeout time.Duration) error {
	httpServer := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
	errs := make(chan error, 1)
	go func() {
		errs <- httpServer.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		return httpServer.Shutdown(shutdownCtx)
	case err := <-errs:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}
