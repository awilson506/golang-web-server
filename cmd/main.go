package main

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

//TODO better way to handle real time hash endpoint counter
type Server struct {
	client api.Client
	server *http.Server
	mux    *http.ServeMux
	stats  *apiStats
}

type apiStats struct {
	Count     int     `json:"total"`
	Average   float32 `json:"average"`
	totalTime float32
}

type StatsLogger struct {
	handler http.Handler
	stats   *apiStats
}

func (l *StatsLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	l.handler.ServeHTTP(w, r)

	if r.URL.Path != "/hash" {
		return
	}

	l.stats.Count++
	l.stats.totalTime = l.stats.totalTime + float32(time.Since(start))/float32(time.Millisecond)
	l.stats.Average = l.stats.totalTime / float32(l.stats.Count)
}

//construct a new stats StatsLogger
func NewStatsLogger(handlerToWrap http.Handler, stats *apiStats) *StatsLogger {
	return &StatsLogger{handlerToWrap, stats}
}

func main() {
	var wg sync.WaitGroup

	s := &Server{
		client: api.New(),
		mux:    http.NewServeMux(),
		stats:  &apiStats{},
	}

	s.server = &http.Server{
		Addr:    ":8080",
		Handler: NewStatsLogger(s.mux, s.stats),
	}

	s.mux.HandleFunc("/hash", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			input, ok := api.ValidateHashRequest(r.PostFormValue("password"))

			if !ok {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(input.Errors)
				return
			}

			count := s.client.UpdatePasswordCount()
			s.client.HandlePassword(&wg, r.FormValue("password"), count)

			json.NewEncoder(w).Encode(count)
		}
	})

	s.mux.HandleFunc("/hash/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			id := strings.TrimPrefix(r.URL.Path, "/hash/")
			hashId, msg, ok := api.ValidateHashGetRequest(id)

			if !ok {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(msg.Errors)
				return
			}

			json.NewEncoder(w).Encode(s.client.Get(hashId))
		}
	})

	s.mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			json.NewEncoder(w).Encode(s.stats)
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s.mux.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
		cancel()
	})

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	select {
	case <-ctx.Done():
		//wait for all of the inflight tasks
		wg.Wait()
		s.server.Shutdown(ctx)
	}

	log.Printf("Server shutdown gracefully")
}
