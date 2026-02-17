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

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		cfg.fileserverHits.Add(1)
		
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	
	myint := cfg.fileserverHits.Load()
	
	w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>",myint)))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusOK)
	
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func handler(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
    w.Write(dat)
}



func main() {
	const filepath = "."
	const port = "8080"
	var apiCfg apiConfig

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepath)))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileServerHandler))
	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	fmt.Println("Server starting on port 8080...")
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}

}
