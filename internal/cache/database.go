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

func initDatabase(db *sql.DB) (err error) {
	ctx := context.Background()

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS genres (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
		);
	`)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS artists (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			name TEXT NOT NULL UNIQUE
		);
	`)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS artist_genres (
			artist_id INTEGER NOT NULL,
			genre_id INTEGER NOT NULL,
			PRIMARY KEY (artist_id, genre_id),
			FOREIGN KEY (artist_id) REFERENCES artists(id) ON DELETE CASCADE,
			FOREIGN KEY (genre_id) REFERENCES genres(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS albums (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			artist_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			album_mbid TEXT NOT NULL UNIQUE,
			FOREIGN KEY (artist_id) REFERENCES artists(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS songs (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			album_id INTEGER,
			title TEXT NOT NULL, 
			mbid TEXT NOT NULL,
			FOREIGN KEY (album_id) REFERENCES albums(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS song_artists (
			song_id INTEGER NOT NULL,
			artist_id INTEGER NOT NULL,
			PRIMARY KEY (song_id, artist_id),
			FOREIGN KEY (song_id) REFERENCES songs(id) ON DELETE CASCADE,
			FOREIGN KEY (artist_id) REFERENCES artists(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		return err
	}

	return nil
}