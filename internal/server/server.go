package server

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PouryaMansouri/BucketDesk/internal/profiles"
	"github.com/PouryaMansouri/BucketDesk/internal/storage"
)

type Server struct {
	store  *profiles.Store
	logger *slog.Logger
}

func New(store *profiles.Store, logger *slog.Logger) *Server {
	return &Server{store: store, logger: logger}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", s.handleHealth)
	mux.HandleFunc("GET /api/profiles", s.handleListProfiles)
	mux.HandleFunc("POST /api/profiles", s.handleSaveProfile)
	mux.HandleFunc("DELETE /api/profiles/", s.handleDeleteProfile)
	mux.HandleFunc("POST /api/test", s.handleTest)
	mux.HandleFunc("GET /api/objects", s.handleListObjects)
	mux.HandleFunc("POST /api/upload", s.handleUpload)
	mux.HandleFunc("DELETE /api/objects", s.handleDeleteObjects)
	mux.HandleFunc("POST /api/shutdown", s.handleShutdown)
	mux.Handle("/", s.staticUI())

	return s.withMiddleware(mux)
}

func (s *Server) handleShutdown(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "shutting down"})
	go func() {
		time.Sleep(250 * time.Millisecond)
		os.Exit(0)
	}()
}

func (s *Server) withMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func (s *Server) handleListProfiles(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"profiles": s.store.List()})
}

func (s *Server) handleSaveProfile(w http.ResponseWriter, r *http.Request) {
	var profile profiles.Profile
	if err := decodeJSON(r, &profile); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	saved, err := s.store.Save(profile)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	writeJSON(w, http.StatusOK, saved)
}

func (s *Server) handleDeleteProfile(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/profiles/")
	if id == "" {
		writeError(w, http.StatusBadRequest, errors.New("profile id is required"))
		return
	}
	if err := s.store.Delete(id); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (s *Server) handleTest(w http.ResponseWriter, r *http.Request) {
	var profile profiles.Profile
	if err := decodeJSON(r, &profile); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	profile = profile.Normalized()
	if profile.SecretKey == "" && profile.ID != "" {
		if saved, ok := s.store.Get(profile.ID); ok {
			profile.SecretKey = saved.SecretKey
		}
	}
	if err := profile.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	result, err := storage.New(profile).Test(r.Context())
	if err != nil {
		s.logger.Warn("storage test failed", "error", err)
		writeError(w, http.StatusBadGateway, errors.New("connection failed or write permission is missing"))
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) handleListObjects(w http.ResponseWriter, r *http.Request) {
	client, ok := s.clientForRequest(w, r)
	if !ok {
		return
	}

	limit := int32(100)
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil {
			limit = int32(parsed)
		}
	}

	result, err := client.List(
		r.Context(),
		r.URL.Query().Get("prefix"),
		r.URL.Query().Get("token"),
		limit,
	)
	if err != nil {
		s.logger.Warn("list objects failed", "error", err)
		writeError(w, http.StatusBadGateway, errors.New("failed to browse bucket"))
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	client, ok := s.clientForRequest(w, r)
	if !ok {
		return
	}

	if err := r.ParseMultipartForm(512 << 20); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		writeError(w, http.StatusBadRequest, errors.New("no file selected"))
		return
	}

	uploaded := make([]storage.Object, 0, len(files))
	for _, header := range files {
		file, err := header.Open()
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		object, err := client.Upload(r.Context(), r.URL.Query().Get("prefix"), header.Filename, header.Header.Get("Content-Type"), file)
		closeErr := file.Close()
		if err != nil {
			writeError(w, http.StatusBadGateway, err)
			return
		}
		if closeErr != nil {
			writeError(w, http.StatusInternalServerError, closeErr)
			return
		}
		uploaded = append(uploaded, object)
	}

	writeJSON(w, http.StatusOK, map[string]any{"objects": uploaded})
}

func (s *Server) handleDeleteObjects(w http.ResponseWriter, r *http.Request) {
	client, ok := s.clientForRequest(w, r)
	if !ok {
		return
	}

	var payload struct {
		Keys []string `json:"keys"`
	}
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if len(payload.Keys) == 0 {
		writeError(w, http.StatusBadRequest, errors.New("no object keys selected"))
		return
	}

	if err := client.Delete(r.Context(), payload.Keys); err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"deleted": len(payload.Keys)})
}

func (s *Server) clientForRequest(w http.ResponseWriter, r *http.Request) (*storage.Client, bool) {
	id := r.URL.Query().Get("profile")
	if id == "" {
		writeError(w, http.StatusBadRequest, errors.New("profile is required"))
		return nil, false
	}

	profile, ok := s.store.Get(id)
	if !ok {
		writeError(w, http.StatusNotFound, errors.New("profile not found"))
		return nil, false
	}

	return storage.New(profile), true
}

func decodeJSON(r *http.Request, target any) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func (s *Server) staticUI() http.Handler {
	content, err := fs.Sub(webFiles, "web/dist")
	if err != nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`<html><body style="font-family:sans-serif;padding:32px"><h1>BucketDesk API</h1><p>Build the React UI with <code>npm run build:web</code>, then run the binary again.</p></body></html>`))
		})
	}

	fileServer := http.FileServer(http.FS(content))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fileServer.ServeHTTP(w, r)
	})
}
