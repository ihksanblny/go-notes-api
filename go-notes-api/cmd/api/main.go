package main

import (
	"log"
	"net/http"

	"go-notes-api/internal/notes"
	"go-notes-api/internal/server"
)

func main() {
	//Load config
	cfg := server.LoadConfig()
	
	// in memory store for testing
	store := notes.NewInMemoryStore()
	notes.SeedInitialNotes(store)

	handler := notes.NewHandler(store)

	mux := http.NewServeMux()
	mux.Handle("/notes", http.HandlerFunc(handler.HandleNotes)) // Handle /notes/ endpoint
	mux.Handle("/notes/", http.HandlerFunc(handler.HandleNotesByID)) // Handle /notes/{id} endpoint
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request){
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}) //cek healthz

	handlerWithMiddleware := server.WithCors(mux, cfg.AllowedOrigins)
	handlerWithMiddleware = server.WithLogging(handlerWithMiddleware)

	addr := ":"+cfg.Port
	log.Println("Starting server on", addr)

	if err := http.ListenAndServe(addr, handlerWithMiddleware); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}