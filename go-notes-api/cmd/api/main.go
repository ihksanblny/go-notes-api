package main

import (
	"database/sql"
	"log"
	"net/http"

	"go-notes-api/internal/notes"
	"go-notes-api/internal/server"

	_ "modernc.org/sqlite"
)

func main() {
	//Load config
	cfg := server.LoadConfig()

	//Buat file sqlite
	db, err := sql.Open("sqlite", "file:notes.db?_foreign_keys=on")
	if err != nil {
		log.Fatal("failed to open database:", err)
	}
	defer db.Close()

	if err := notes.InitSchema(db); err != nil {
		log.Fatal("failed to init schema:", err)
	}

	store := notes.NewSQLiteStore(db)
	notes.SeedInitialNotes(store)

	handler := notes.NewHandler(store)

	mux := http.NewServeMux()
	mux.Handle("/notes", http.HandlerFunc(handler.HandleNotes))      // Handle /notes/ endpoint
	mux.Handle("/notes/", http.HandlerFunc(handler.HandleNotesByID)) // Handle /notes/{id} endpoint
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}) //cek healthz

	handlerWithMiddleware := server.WithCors(mux, cfg.AllowedOrigins)
	handlerWithMiddleware = server.WithLogging(handlerWithMiddleware)

	addr := ":" + cfg.Port
	log.Println("Starting server on", addr)

	if err := http.ListenAndServe(addr, handlerWithMiddleware); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
