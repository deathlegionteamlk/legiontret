package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/deathlegionteam/legiontret/internal/config"
	"github.com/deathlegionteam/legiontret/internal/llama"
	"github.com/deathlegionteam/legiontret/internal/model"
	"github.com/deathlegionteam/legiontret/internal/registry"
)

// Server is the REST API server
type Server struct {
	cfg      *config.Config
	manager  *model.Manager
	registry *registry.Registry
	llama    *llama.Server
	server   *http.Server
	model    string // currently loaded model
}

// NewServer creates a new API server
func NewServer(cfg *config.Config, mgr *model.Manager, reg *registry.Registry, llamaServer *llama.Server) *Server {
	return &Server{
		cfg:      cfg,
		manager:  mgr,
		registry: reg,
		llama:    llamaServer,
	}
}

// Start starts the API server
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Core API endpoints (Ollama-compatible)
	mux.HandleFunc("/api/generate", s.handleGenerate)
	mux.HandleFunc("/api/chat", s.handleChat)
	mux.HandleFunc("/api/embeddings", s.handleEmbeddings)
	mux.HandleFunc("/api/pull", s.handlePull)
	mux.HandleFunc("/api/push", s.handlePush)
	mux.HandleFunc("/api/delete", s.handleDelete)
	mux.HandleFunc("/api/show", s.handleShow)
	mux.HandleFunc("/api/copy", s.handleCopy)
	mux.HandleFunc("/api/tags", s.handleTags)
	mux.HandleFunc("/api/ps", s.handlePS)
	mux.HandleFunc("/api/version", s.handleVersion)

	// Health and info
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/", s.handleRoot)

	// LegionTret-specific enhanced endpoints
	mux.HandleFunc("/api/v1/models", s.handleV1Models)      // OpenAI-compatible
	mux.HandleFunc("/api/v1/chat/completions", s.handleV1Chat) // OpenAI-compatible
	mux.HandleFunc("/api/v1/completions", s.handleV1Completions) // OpenAI-compatible
	mux.HandleFunc("/api/search", s.handleSearch)
	mux.HandleFunc("/api/info", s.handleInfo)
	mux.HandleFunc("/api/system", s.handleSystem)

	s.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port),
		Handler: corsMiddleware(mux),
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		s.Stop()
	}()

	fmt.Printf("LegionTret API server listening on http://%s:%d\n", s.cfg.Host, s.cfg.Port)
	return s.server.ListenAndServe()
}

// Stop stops the API server
func (s *Server) Stop() {
	if s.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.server.Shutdown(ctx)
	}
}

// corsMiddleware adds CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// handleRoot handles the root endpoint
func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"name":    "legiontret",
		"version": config.Version,
		"org":     config.OrgName,
		"message": "LegionTret is running! Visit /api/tags to see available models.",
	})
}

// handleHealth handles health checks
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// handleVersion returns version info
func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"version": config.Version})
}

// GenerateRequest represents a generation request
type GenerateRequest struct {
	Model     string         `json:"model"`
	Prompt    string         `json:"prompt"`
	System    string         `json:"system"`
	Template  string         `json:"template"`
	Context   []int          `json:"context"`
	Stream    *bool          `json:"stream"`
	Raw       bool           `json:"raw"`
	Format    string         `json:"format"`
	Options   GenerateOptions `json:"options"`
}

// GenerateOptions holds generation parameters
type GenerateOptions struct {
	NumPredict   int     `json:"num_predict"`
	Temperature  float64 `json:"temperature"`
	TopP         float64 `json:"top_p"`
	TopK         int     `json:"top_k"`
	RepeatPenalty float64 `json:"repeat_penalty"`
	Seed         int     `json:"seed"`
	Stop         []string `json:"stop"`
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   *bool         `json:"stream"`
	Format   string        `json:"format"`
	Options  GenerateOptions `json:"options"`
}

// ChatMessage represents a single chat message
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Images  []string `json:"images,omitempty"`
}

// handleGenerate handles text generation requests
func (s *Server) handleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if req.Model == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "model is required"})
		return
	}

	if req.Prompt == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "prompt is required"})
		return
	}

	stream := true
	if req.Stream != nil {
		stream = *req.Stream
	}

	// Forward to llama.cpp server
	llamaURL := fmt.Sprintf("http://%s:%d/completion", s.cfg.Host, s.cfg.Port)

	payload := map[string]interface{}{
		"prompt":      req.Prompt,
		"n_predict":   req.Options.NumPredict,
		"temperature": req.Options.Temperature,
		"top_p":       req.Options.TopP,
		"top_k":       req.Options.TopK,
		"repeat_penalty": req.Options.RepeatPenalty,
		"stream":      stream,
	}

	if req.System != "" {
		payload["system_prompt"] = req.System
	}

	if len(req.Options.Stop) > 0 {
		payload["stop"] = req.Options.Stop
	}

	if stream {
		s.streamProxy(w, r, llamaURL, payload)
	} else {
		s.proxyRequest(w, r, llamaURL, payload)
	}
}

// handleChat handles chat completion requests
func (s *Server) handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Convert chat messages to a single prompt
	var promptBuilder strings.Builder
	for _, msg := range req.Messages {
		switch msg.Role {
		case "system":
			promptBuilder.WriteString(fmt.Sprintf("System: %s\n", msg.Content))
		case "user":
			promptBuilder.WriteString(fmt.Sprintf("User: %s\n", msg.Content))
		case "assistant":
			promptBuilder.WriteString(fmt.Sprintf("Assistant: %s\n", msg.Content))
		}
	}
	promptBuilder.WriteString("Assistant: ")

	stream := true
	if req.Stream != nil {
		stream = *req.Stream
	}

	llamaURL := fmt.Sprintf("http://%s:%d/completion", s.cfg.Host, s.cfg.Port)

	payload := map[string]interface{}{
		"prompt":    promptBuilder.String(),
		"n_predict": req.Options.NumPredict,
		"temperature": req.Options.Temperature,
		"top_p":     req.Options.TopP,
		"stream":    stream,
	}

	if stream {
		s.streamProxy(w, r, llamaURL, payload)
	} else {
		s.proxyRequest(w, r, llamaURL, payload)
	}
}

// handleEmbeddings handles embedding requests
func (s *Server) handleEmbeddings(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	llamaURL := fmt.Sprintf("http://%s:%d/embedding", s.cfg.Host, s.cfg.Port)
	payload := map[string]interface{}{
		"content": req.Prompt,
	}
	s.proxyRequest(w, r, llamaURL, payload)
}

// handlePull handles model pull/download requests
func (s *Server) handlePull(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name     string `json:"name"`
		Insecure bool   `json:"insecure"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status": fmt.Sprintf("pulling %s - use CLI for download progress", req.Name),
	})
}

// handlePush handles model push requests
func (s *Server) handlePush(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "push not yet supported"})
}

// handleDelete handles model deletion requests
func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err := s.manager.Delete(req.Name); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// handleShow handles model info requests
func (s *Server) handleShow(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	info, err := s.manager.GetModelInfo(req.Name)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, info)
}

// handleCopy handles model copy requests
func (s *Server) handleCopy(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "copy not yet supported"})
}

// handleTags lists available models
func (s *Server) handleTags(w http.ResponseWriter, r *http.Request) {
	localModels, err := s.manager.ListLocal()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	type tagModel struct {
		Name       string `json:"name"`
		Model      string `json:"model"`
		ModifiedAt string `json:"modified_at"`
		Size       int64  `json:"size"`
		Digest     string `json:"digest"`
		Details    struct {
			ParentModel   string `json:"parent_model"`
			Format        string `json:"format"`
			Family        string `json:"family"`
			Families      []string `json:"families"`
			ParameterSize string `json:"parameter_size"`
			Quantization  string `json:"quantization_level"`
		} `json:"details"`
	}

	var models []tagModel
	for _, m := range localModels {
		tm := tagModel{
			Name:       m.Name,
			Model:      m.Name,
			ModifiedAt: m.ModifiedAt.Format(time.RFC3339),
			Size:       m.Size,
		}
		tm.Details.Format = "gguf"
		tm.Details.Family = m.Family
		tm.Details.Families = []string{m.Family}
		tm.Details.ParameterSize = m.Parameters
		models = append(models, tm)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"models": models,
	})
}

// handlePS shows running models
func (s *Server) handlePS(w http.ResponseWriter, r *http.Request) {
	models := []interface{}{}
	if s.model != "" && s.llama.IsRunning() {
		models = append(models, map[string]interface{}{
			"name":     s.model,
			"model":    s.model,
			"size":     0,
			"digest":   "",
			"expires_at": time.Now().Add(5 * time.Minute).Format(time.RFC3339),
		})
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"models": models})
}

// handleSearch searches models
func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "query parameter 'q' is required"})
		return
	}

	results := s.registry.SearchModels(query)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"models": results,
	})
}

// handleInfo returns system info
func (s *Server) handleInfo(w http.ResponseWriter, r *http.Request) {
	info := map[string]interface{}{
		"version":     config.Version,
		"org":         config.OrgName,
		"models_dir":  s.cfg.ModelsDir,
		"binaries_dir": s.cfg.BinariesDir,
		"api_url":     s.cfg.APIBaseURL(),
		"system":      llama.GetSystemInfo(),
	}
	writeJSON(w, http.StatusOK, info)
}

// handleSystem returns system info
func (s *Server) handleSystem(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"system": llama.GetSystemInfo(),
	})
}

// OpenAI-compatible endpoints

// handleV1Models handles OpenAI-compatible model listing
func (s *Server) handleV1Models(w http.ResponseWriter, r *http.Request) {
	localModels, _ := s.manager.ListLocal()

	type v1Model struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Created int64  `json:"created"`
		OwnedBy string `json:"owned_by"`
	}

	var models []v1Model
	for _, m := range localModels {
		models = append(models, v1Model{
			ID:      m.Name,
			Object:  "model",
			Created: m.ModifiedAt.Unix(),
			OwnedBy: "legiontret",
		})
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"object": "list",
		"data":   models,
	})
}

// handleV1Chat handles OpenAI-compatible chat completions
func (s *Server) handleV1Chat(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Convert to Ollama format and forward
	s.handleChat(w, r)
}

// handleV1Completions handles OpenAI-compatible completions
func (s *Server) handleV1Completions(w http.ResponseWriter, r *http.Request) {
	s.handleGenerate(w, r)
}

// Helper functions

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// streamProxy proxies a request with streaming response
func (s *Server) streamProxy(w http.ResponseWriter, r *http.Request, url string, payload interface{}) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	proxyReq, err := http.NewRequestWithContext(r.Context(), "POST", url, strings.NewReader(string(payloadBytes)))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	proxyReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": fmt.Sprintf("failed to connect to model server: %v", err)})
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		io.Copy(w, resp.Body)
		return
	}

	buf := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			w.Write(buf[:n])
			flusher.Flush()
		}
		if err != nil {
			break
		}
	}
}

// proxyRequest proxies a non-streaming request
func (s *Server) proxyRequest(w http.ResponseWriter, r *http.Request, url string, payload interface{}) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	proxyReq, err := http.NewRequestWithContext(r.Context(), "POST", url, strings.NewReader(string(payloadBytes)))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	proxyReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": fmt.Sprintf("failed to connect to model server: %v", err)})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

// logRequest is a middleware that logs requests
func logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next(w, r)
	}
}
