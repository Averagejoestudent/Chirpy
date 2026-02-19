package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Averagejoestudent/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform		string
}

func main() {
	const filepath = "."
	const port = "8080"
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)

	cfg := Config{
		db: dbQueries,
		platform: os.Getenv("PLATFORM"),
	}

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepath)))
	mux.Handle("/app/",cfg.middlewareMetricsInc(fileServerHandler))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", cfg.chripsHandler)
	mux.HandleFunc("POST /api/users", cfg.userHandler)
	mux.HandleFunc("GET /api/chirps", cfg.GetchripsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.GetOnechripsHandler)

	fmt.Println("Server starting on port 8080...")
	err = server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}

}
