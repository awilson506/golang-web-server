package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	api "github.com/awilson506/golang-web-server/pkg"
)

//struct to hold server components and middleware stats
type Server struct {
	client *api.Client
	server *http.Server
	mux    *http.ServeMux
	stats  *apiStats
	wg     *sync.WaitGroup
}

//struct to hold api usage data
type apiStats struct {
	Count     int     `json:"total"`
	Average   float32 `json:"average"`
	totalTime float32
}

//middleware handler struct
type StatsLogger struct {
	handler http.Handler
	stats   *apiStats
}

/* Handle all incoming requests and
track usage if it's a hash request */
func (l *StatsLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	l.handler.ServeHTTP(w, r)

	if r.URL.Path != "/hash" {
		return
	}

	l.stats.totalTime = l.stats.totalTime + float32(time.Since(start))/float32(time.Millisecond)
	l.stats.Average = l.stats.totalTime / float32(l.stats.Count)
}

/* Get an instance of the server and setup
stats struct for tracking hash api calls */
func NewServer() *Server {
	s := &Server{
		client: api.New(),
		mux:    http.NewServeMux(),
		stats:  &apiStats{},
		wg:     &sync.WaitGroup{},
	}

	s.server = &http.Server{
		Addr:    ":8080",
		Handler: &StatsLogger{s.mux, s.stats},
	}
	s.mux.HandleFunc("/hash", s.hashPasswordHandler)
	s.mux.HandleFunc("/hash/", s.getPasswordHandler)
	s.mux.HandleFunc("/stats", s.statsHandler)
	s.mux.HandleFunc("/shutdown", s.shutdownHandler)

	return s
}

//start up the server
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

/* Shutdown the server gracefully, this will allow
any inflight processes finish out */
func (s *Server) GracefulShutdown(wg *sync.WaitGroup) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg.Wait()
	s.server.Shutdown(ctx)

	log.Printf("Server shutdown gracefully")
}

// Update the counter for the hash password stats
func (s *Server) UpdatePasswordCount() int {
	s.stats.Count++
	return s.stats.Count
}

/* If we hit validation issues write the proper
status code and return the validation issues */
func (s *Server) WriteErrorResponse(w http.ResponseWriter, errors map[string]string) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(errors)
}

/* Send the plain text password off to be hashed
in a go routine */
func (s *Server) hashPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		input, ok := api.ValidateHashRequest(r.PostFormValue("password"))

		if !ok {
			s.WriteErrorResponse(w, input.Errors)
			return
		}

		count := s.UpdatePasswordCount()
		s.client.HandlePassword(s.wg, r.FormValue("password"), count)

		json.NewEncoder(w).Encode(count)
	}
}

/* Retrieve a password that was stored and
hashed */
func (s *Server) getPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		id := strings.TrimPrefix(r.URL.Path, "/hash/")
		hashId, msg, ok := api.ValidateHashGetRequest(id)

		if !ok {
			s.WriteErrorResponse(w, msg.Errors)
			return
		}

		json.NewEncoder(w).Encode(s.client.Get(hashId))
	}
}

/* Get the stats on the hash api calls and
processing time */
func (s *Server) statsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		json.NewEncoder(w).Encode(s.stats)
	}
}

//Shutdown the server
func (s *Server) shutdownHandler(w http.ResponseWriter, r *http.Request) {
	go s.GracefulShutdown(s.wg)
	w.Write([]byte("OK"))
}
