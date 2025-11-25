package main

import (
	"log"
	"net/http"

	"go-notes-api/internal/notes"
)

func main() {
	// in memory store for testing
	store := notes.NewInMemoryStore()
	notes.SeedInitialNotes(store)

	handler := notes.NewHandler(store)

	mux := http.NewServeMux()
	mux.Handle("/notes/", http.HandlerFunc(handler.HandleNotes)) // Handle /notes/ endpoint
	mux.Handle("/notes/", http.HandlerFunc(handler.HandleNotesByID)) // Handle /notes/{id} endpoint

	handlerWithMiddleware := server.WithCors(mux, "http://localhost:5173")

	addr := ":8080"
	log.Println("Starting server on", addr)

	if err := http.ListenAndServe(addr, handlerWithMiddleware); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}