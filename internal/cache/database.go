package cache

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func OpenDatabase() (*sql.DB, error) {
	cachePath, err := getCachePath()
	if err != nil {
		return nil, fmt.Errorf("cannot get cache path: %w", err)
	}

	dbPath := filepath.Join(cachePath, "tunectl.sqlite")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func initDatabase(db *sql.DB) error {
	ctx := context.Background()

	stmts := []string{
		`CREATE TABLE IF NOT EXISTS genres (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
		);`,

		`CREATE TABLE IF NOT EXISTS artists (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
		);`,

		`CREATE TABLE IF NOT EXISTS artist_genres (
			artist_id INTEGER NOT NULL,
			genre_id INTEGER NOT NULL,
			PRIMARY KEY (artist_id, genre_id),
			FOREIGN KEY (artist_id) REFERENCES artists(id) ON DELETE CASCADE,
			FOREIGN KEY (genre_id) REFERENCES genres(id) ON DELETE CASCADE
		);`,

		`CREATE TABLE IF NOT EXISTS albums (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			album_mbid TEXT UNIQUE,
			artist_id INTEGER,
			FOREIGN KEY (artist_id) REFERENCES artists(id) ON DELETE SET NULL
		);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_albums_name_artist
			ON albums(name, artist_id);`,

		`CREATE TABLE IF NOT EXISTS songs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			mbid TEXT UNIQUE,
			album_id INTEGER,
			FOREIGN KEY (album_id) REFERENCES albums(id) ON DELETE SET NULL
		);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_songs_title_album
			ON songs(title, album_id);`,

		`CREATE TABLE IF NOT EXISTS song_artists (
			song_id INTEGER NOT NULL,
			artist_id INTEGER NOT NULL,
			PRIMARY KEY (song_id, artist_id),
			FOREIGN KEY (song_id) REFERENCES songs(id) ON DELETE CASCADE,
			FOREIGN KEY (artist_id) REFERENCES artists(id) ON DELETE CASCADE
		);`,
	}

	for _, stmt := range stmts {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("init database: %w", err)
		}
	}

	return nil
}
