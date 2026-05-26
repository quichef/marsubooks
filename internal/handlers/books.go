package handlers

import (
	"bytes"
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"bookshelf/internal/models"
	"bookshelf/internal/openlibrary"

	"github.com/go-chi/chi/v5"
	"github.com/yuin/goldmark"
)

type BookHandler struct {
	DB   *sql.DB
	Tmpl *template.Template
}

func (h *BookHandler) List(w http.ResponseWriter, r *http.Request) {
	books, err := models.ListBooks(h.DB)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	h.Tmpl.ExecuteTemplate(w, "layout.html", map[string]any{
		"Title": "Ma bibliothèque",
		"Books": books,
		"Page":  "index",
	})
}

func (h *BookHandler) New(w http.ResponseWriter, r *http.Request) {
	h.Tmpl.ExecuteTemplate(w, "layout.html", map[string]any{
		"Title": "Ajouter un livre",
		"Book":  &models.Book{},
		"Page":  "form",
	})
}

func (h *BookHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	rating, _ := strconv.Atoi(r.FormValue("rating"))
	b := &models.Book{
		Title:    r.FormValue("title"),
		Author:   r.FormValue("author"),
		ISBN:     r.FormValue("isbn"),
		Genre:    r.FormValue("genre"),
		CoverURL: r.FormValue("cover_url"),
		DateRead: r.FormValue("date_read"),
		Rating:   rating,
		Review:   r.FormValue("review"),
	}
	id, err := models.CreateBook(h.DB, b)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/books/"+strconv.FormatInt(id, 10), http.StatusSeeOther)
}

func (h *BookHandler) Detail(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	b, err := models.GetBook(h.DB, id)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var buf bytes.Buffer
	goldmark.Convert([]byte(b.Review), &buf)
	h.Tmpl.ExecuteTemplate(w, "layout.html", map[string]any{
		"Title":      b.Title,
		"Book":       b,
		"ReviewHTML": template.HTML(buf.String()),
		"Page":       "detail",
	})
}

func (h *BookHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	b, err := models.GetBook(h.DB, id)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	h.Tmpl.ExecuteTemplate(w, "layout.html", map[string]any{
		"Title": "Modifier — " + b.Title,
		"Book":  b,
		"Page":  "form",
	})
}

func (h *BookHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	rating, _ := strconv.Atoi(r.FormValue("rating"))
	b := &models.Book{
		ID:       id,
		Title:    r.FormValue("title"),
		Author:   r.FormValue("author"),
		ISBN:     r.FormValue("isbn"),
		Genre:    r.FormValue("genre"),
		CoverURL: r.FormValue("cover_url"),
		DateRead: r.FormValue("date_read"),
		Rating:   rating,
		Review:   r.FormValue("review"),
	}
	if err := models.UpdateBook(h.DB, b); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/books/"+strconv.FormatInt(id, 10), http.StatusSeeOther)
}

func (h *BookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err := models.DeleteBook(h.DB, id); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *BookHandler) OpenLibrarySearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if len(q) < 3 {
		w.Write([]byte(""))
		return
	}
	results, err := openlibrary.Search(q)
	if err != nil {
		log.Printf("openlibrary search: %v", err)
		w.Write([]byte(""))
		return
	}
	h.Tmpl.ExecuteTemplate(w, "search_results.html", results)
}
