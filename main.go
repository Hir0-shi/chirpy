package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{}

	mux := http.NewServeMux()

	// website endpoint
	fileServer := http.FileServer(http.Dir(filepathRoot))
	// handle /app
	mux.Handle("/app",
		apiCfg.middlewareMetricsInc(
			http.StripPrefix("/app", fileServer),
		),
	)
	// handle /app/
	mux.Handle("/app/",
		apiCfg.middlewareMetricsInc(
			http.StripPrefix("/app", fileServer),
		),
	)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /api/reset", apiCfg.handlerReset)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	// w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
	// build bytes without fmt.Sprintf -> []byte(...) roundtrip
	var b []byte
	b = fmt.Appendf(b, "Hits: %d", cfg.fileserverHits.Load())
	w.Write(b)
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
