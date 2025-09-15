package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"packsolver/packer"
)

type requestPayload struct {
	Amount int   `json:"amount"`
	Sizes  []int `json:"sizes"`
}

func main() {
	mux := http.NewServeMux()

	// Serve UI (index.html)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("static", "index.html"))
	})

	// API endpoint
	mux.HandleFunc("/api/pack", handlePack)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("listening on http://localhost:8080")
	log.Fatal(srv.ListenAndServe())
}

func handlePack(w http.ResponseWriter, r *http.Request) {
	var p requestPayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if p.Amount <= 0 || len(p.Sizes) == 0 {
		http.Error(w, "amount > 0 and at least one size required", http.StatusBadRequest)
		return
	}

	result, err := packer.Solve(p.Amount, p.Sizes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
