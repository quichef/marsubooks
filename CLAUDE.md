# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run locally (must be run from project root so templates/ and static/ are found)
mkdir -p data
go run ./cmd/server

# Build
go build -o bookshelf ./cmd/server

# After adding/changing dependencies
go mod tidy

# Docker (preferred for deployment on Debian VM)
docker compose up --build
docker compose up -d   # background
```

Environment variables (defaults shown):
- `DB_PATH=./data/books.db`
- `PORT=8080`

## Architecture

Single binary Go server. No build step for frontend (no npm/webpack).

**Request flow:** `chi router` → `handlers.BookHandler` or `handlers.ExportHandler` → `models` (raw `database/sql`) → SQLite

**Template rendering:** All HTML is server-rendered. `main.go` parses all `templates/*.html` at startup into a single `*template.Template` set. Every page route calls `ExecuteTemplate(w, "layout.html", data)` with a `Page` field (`"index"`, `"form"`, `"detail"`); `layout.html` uses `{{if eq .Page "..."}}` to include the right sub-template.

Two custom template funcs registered in `main.go`: `starRange` and `emptyRange` — used to render star ratings as ★/☆ in templates.

**Database:** SQLite via `modernc.org/sqlite` (pure Go, no CGO). Schema is inlined in `internal/db/db.go:migrate()` and applied automatically on startup. The `data/` directory is git-ignored; when running locally it must exist before starting.

**OpenLibrary autocomplete:** The form page uses HTMX to call `/api/openlibrary?q=...` (min 3 chars) which proxies to `openlibrary.org/search.json`. The returned HTML fragment (`search_results.html`) renders clickable rows; clicking one fills the form fields via a `data-*` attribute + JS event listener in `form.html`.

**Export for AI recommendations:** `/export/json` and `/export/csv` produce downloadable files. The JSON envelope format (`exported_at` + `books[]`) is designed to be passed directly to Claude CLI:
```bash
claude "Recommande-moi un livre : $(curl -s http://localhost:8080/export/json)"
```

**DELETE workaround:** HTML forms don't support DELETE/PUT. Delete is routed as `POST /books/{id}/delete`.

## Key constraints

- The server must be started from the project root — it looks for `templates/` and `static/` relative to the working directory (not embedded).
- `modernc.org/sqlite` requires `CGO_ENABLED=0` for the Docker build (already set in Dockerfile).
- The SQLite WAL files (`data/*.db-shm`, `data/*.db-wal`) are normal — SQLite WAL mode is enabled at startup.
