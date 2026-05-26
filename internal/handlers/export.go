package handlers

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"bookshelf/internal/models"
)

type ExportHandler struct {
	DB *sql.DB
}

type exportEnvelope struct {
	ExportedAt string        `json:"exported_at"`
	Books      []exportBook  `json:"books"`
}

type exportBook struct {
	Title    string `json:"title"`
	Author   string `json:"author"`
	Genre    string `json:"genre,omitempty"`
	DateRead string `json:"date_read,omitempty"`
	Rating   int    `json:"rating,omitempty"`
	Review   string `json:"review,omitempty"`
}

func (h *ExportHandler) JSON(w http.ResponseWriter, r *http.Request) {
	books, err := models.ListBooks(h.DB)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	env := exportEnvelope{
		ExportedAt: time.Now().Format("2006-01-02"),
	}
	for _, b := range books {
		env.Books = append(env.Books, exportBook{
			Title:    b.Title,
			Author:   b.Author,
			Genre:    b.Genre,
			DateRead: b.DateRead,
			Rating:   b.Rating,
			Review:   b.Review,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="books-%s.json"`, env.ExportedAt))
	json.NewEncoder(w).Encode(env)
}

func (h *ExportHandler) CSV(w http.ResponseWriter, r *http.Request) {
	books, err := models.ListBooks(h.DB)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	date := time.Now().Format("2006-01-02")
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="books-%s.csv"`, date))

	cw := csv.NewWriter(w)
	cw.Write([]string{"title", "author", "genre", "date_read", "rating", "review"})
	for _, b := range books {
		cw.Write([]string{
			b.Title, b.Author, b.Genre, b.DateRead,
			fmt.Sprintf("%d", b.Rating), b.Review,
		})
	}
	cw.Flush()
}
