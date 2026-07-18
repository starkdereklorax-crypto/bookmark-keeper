package main

import (
	"bookmark-keeper/internal/handlers"
	"bookmark-keeper/internal/storage"
	"log"
	"net/http"
	"os"
)

func main() {
	store := storage.NewStorage()
	handler := handlers.NewHandler(store)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /bookmarks/{id}", handler.GetBookmark)
	mux.HandleFunc("POST /bookmarks", handler.CreateBookmark)
	mux.HandleFunc("GET /bookmarks", handler.ListBookmarks)
	mux.HandleFunc("DELETE /bookmarks/{id}", handler.DeleteBookmark)
	mux.Handle("/", http.FileServer(http.Dir("web")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
