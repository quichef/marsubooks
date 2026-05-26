# Bookshelf

Journal de lecture personnel. Ajout, notation et recherche de livres avec autocomplétion OpenLibrary. Export JSON/CSV pour recommandations IA.

## Stack

- Go + chi router, templates HTML server-side
- SQLite (pure Go, pas de CGO)
- HTMX pour l'autocomplétion
- Docker

## Déploiement (VM Debian)

Prérequis : Docker avec le plugin Compose installé.

```bash
mkdir -p /app/books
cd /app/books
curl -O https://raw.githubusercontent.com/quichef/marsubooks/refs/heads/main/docker-compose.yml
docker compose up -d
```

L'application est disponible sur le port `8080`. La base de données SQLite est persistée dans `/app/books/data/books.db`.

## Mise à jour

```bash
cd /app/books
docker compose pull
docker compose up -d
```

## Développement local

```bash
mkdir -p data
go run ./cmd/server
```

Variables d'environnement (valeurs par défaut) :
- `DB_PATH=./data/books.db`
- `PORT=8080`

## Export pour recommandations IA

```bash
claude "Recommande-moi un livre : $(curl -s http://localhost:8080/export/json)"
```
