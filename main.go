package main

import (
	"fmt"
	"log"
	"net/http"
)

func myHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Set the Content-Type header to "text/plain; charset=utf-8"
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	// 2. Set the status code to 200
	w.WriteHeader(http.StatusOK)
	// 3. Write "OK" to the body
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func main() {
	const filepath = "."
	const port = "8080"
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepath))))
	mux.HandleFunc("/healthz", myHandler)

	fmt.Println("Server starting on port 8080...")
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}

}
