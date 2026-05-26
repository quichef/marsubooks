package models

import (
	"database/sql"
	"time"
)

type Book struct {
	ID        int64
	Title     string
	Author    string
	ISBN      string
	Genre     string
	CoverURL  string
	DateRead  string
	Rating    int
	Review    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func ListBooks(db *sql.DB) ([]Book, error) {
	rows, err := db.Query(`
		SELECT id, title, author, isbn, genre, cover_url, date_read, rating, review, created_at, updated_at
		FROM books ORDER BY date_read DESC, created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.ISBN, &b.Genre, &b.CoverURL,
			&b.DateRead, &b.Rating, &b.Review, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, rows.Err()
}

func GetBook(db *sql.DB, id int64) (*Book, error) {
	var b Book
	err := db.QueryRow(`
		SELECT id, title, author, isbn, genre, cover_url, date_read, rating, review, created_at, updated_at
		FROM books WHERE id = ?
	`, id).Scan(&b.ID, &b.Title, &b.Author, &b.ISBN, &b.Genre, &b.CoverURL,
		&b.DateRead, &b.Rating, &b.Review, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func CreateBook(db *sql.DB, b *Book) (int64, error) {
	res, err := db.Exec(`
		INSERT INTO books (title, author, isbn, genre, cover_url, date_read, rating, review)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, b.Title, b.Author, b.ISBN, b.Genre, b.CoverURL, b.DateRead, b.Rating, b.Review)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func UpdateBook(db *sql.DB, b *Book) error {
	_, err := db.Exec(`
		UPDATE books SET title=?, author=?, isbn=?, genre=?, cover_url=?, date_read=?, rating=?, review=?,
		updated_at=CURRENT_TIMESTAMP WHERE id=?
	`, b.Title, b.Author, b.ISBN, b.Genre, b.CoverURL, b.DateRead, b.Rating, b.Review, b.ID)
	return err
}

func DeleteBook(db *sql.DB, id int64) error {
	_, err := db.Exec(`DELETE FROM books WHERE id = ?`, id)
	return err
}
