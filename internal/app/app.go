package app

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"event-service/internal/model"
	"event-service/internal/store"
	"event-service/internal/worker"
	"time"
)

// Config holds the application configuration
type Config struct {
	Port              string
	Env               string
	ProcessingDelayMs int
}

// App represents the HTTP application
type App struct {
	config    Config
	store     *store.Store
	worker    *worker.Worker
	startTime time.Time
	server    *http.Server
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() Config {
	port := getEnv("PORT", "8080")
	env := getEnv("ENV", "dev")
	processingDelayMs := getEnvAsInt("PROCESSING_DELAY_MS", 1000)

	return Config{
		Port:              port,
		Env:               env,
		ProcessingDelayMs: processingDelayMs,
	}
}

// New creates a new application instance
func New(config Config) *App {
	st := store.New()
	wkr := worker.New(st, config.ProcessingDelayMs)

	return &App{
		config:    config,
		store:     st,
		worker:    wkr,
		startTime: time.Now(),
	}
}

// Start starts the HTTP server and background worker
func (a *App) Start() error {
	a.worker.Start()

	mux := http.NewServeMux()
	mux.HandleFunc("/events", a.handleEvents)
	mux.HandleFunc("/health", a.handleHealth)
	mux.HandleFunc("/ready", a.handleReady)

	a.server = &http.Server{
		Addr:    ":" + a.config.Port,
		Handler: mux,
	}

	log.Printf("Starting server on port %s (env: %s)", a.config.Port, a.config.Env)
	return a.server.ListenAndServe()
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown() {
	log.Println("Shutting down application...")
	a.worker.Stop()
	if a.server != nil {
		a.server.Close()
	}
}

// handleEvents handles POST /events
func (a *App) handleEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.EventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.EventID == "" {
		http.Error(w, "event_id is required", http.StatusBadRequest)
		return
	}

	// Check for idempotency
	if a.store.Exists(req.EventID) {
		log.Printf("Event already exists: %s", req.EventID)
		w.WriteHeader(http.StatusConflict)
		return
	}

	// Create and save event
	event := &model.Event{
		EventID: req.EventID,
		Payload: req.Payload,
		Status:  model.StatusAccepted,
	}
	a.store.Save(event)

	// Enqueue for background processing
	a.worker.Enqueue(event)

	log.Printf("Event accepted: %s", req.EventID)
	w.WriteHeader(http.StatusAccepted)
}

// handleHealth handles GET /health
func (a *App) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uptime := time.Since(a.startTime).String()
	resp := model.HealthResponse{
		Status: "ok",
		Uptime: uptime,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleReady handles GET /ready
func (a *App) handleReady(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !a.worker.IsRunning() {
		resp := model.ReadyResponse{
			Status: "not ready",
			Ready:  false,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := model.ReadyResponse{
		Status: "ready",
		Ready:  true,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Helper functions for environment variable parsing

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Invalid value for %s: %s, using default: %d", key, valueStr, defaultValue)
		return defaultValue
	}
	return value
}

// GetServer returns the underlying HTTP server (useful for testing or custom shutdown)
func (a *App) GetServer() *http.Server {
	return a.server
}

// Worker returns the worker instance (useful for testing)
func (a *App) Worker() *worker.Worker {
	return a.worker
}

// Store returns the store instance (useful for testing)
func (a *App) Store() *store.Store {
	return a.store
}
