CREATE TABLE IF NOT EXISTS books (
    id         INTEGER  PRIMARY KEY AUTOINCREMENT,
    title      TEXT     NOT NULL,
    author     TEXT     NOT NULL,
    isbn       TEXT     DEFAULT '',
    genre      TEXT     DEFAULT '',
    cover_url  TEXT     DEFAULT '',
    date_read  TEXT     DEFAULT '',
    rating     INTEGER  DEFAULT 0 CHECK(rating BETWEEN 0 AND 5),
    review     TEXT     DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
