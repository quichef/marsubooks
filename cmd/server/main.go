package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"bookshelf/internal/db"
	"bookshelf/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/books.db"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	database, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer database.Close()

	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"starRange":  func(n int) []int { s := make([]int, n); return s },
		"emptyRange": func(n int) []int { s := make([]int, 5-n); return s },
	}).ParseGlob("templates/*.html"))

	bh := &handlers.BookHandler{DB: database, Tmpl: tmpl}
	eh := &handlers.ExportHandler{DB: database}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	r.Get("/", bh.List)
	r.Get("/books/new", bh.New)
	r.Post("/books", bh.Create)
	r.Get("/books/{id}", bh.Detail)
	r.Get("/books/{id}/edit", bh.Edit)
	r.Post("/books/{id}", bh.Update)
	r.Post("/books/{id}/delete", bh.Delete)

	r.Get("/api/openlibrary", bh.OpenLibrarySearch)

	r.Get("/export/json", eh.JSON)
	r.Get("/export/csv", eh.CSV)

	log.Printf("bookshelf listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
