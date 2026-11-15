package database

import (
	"akerraps/tunectl/internal/cache"
	"akerraps/tunectl/internal/types"
	"context"
	"database/sql"
	"errors"
)

func WriteArtist(song *types.TrackInfo) (err error) {

	return ExecTx(context.Background(), func(tx *sql.Tx) error {
		for _, artist := range song.Artists {
			_, err := tx.ExecContext(context.Background(),
				`INSERT INTO artists (name) VALUES (?)
				ON CONFLICT DO NOTHING;`,
				artist)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// GetArtistGenres returns the genres already cached for an artist. found is
// false when the artist has never been written to the cache, meaning its
// info still needs to be fetched from the API.
func GetArtistGenres(name string) (genres []string, found bool, err error) {
	db, err := cache.OpenDatabase()
	if err != nil {
		return nil, false, err
	}
	defer db.Close()

	ctx := context.Background()

	var artistID int64
	err = db.QueryRowContext(ctx,
		`SELECT id FROM artists WHERE name = ?;`, name).Scan(&artistID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	rows, err := db.QueryContext(ctx,
		`SELECT g.name FROM genres g
		 JOIN artist_genres ag ON ag.genre_id = g.id
		 WHERE ag.artist_id = ?;`, artistID)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	for rows.Next() {
		var genre string
		if err := rows.Scan(&genre); err != nil {
			return nil, false, err
		}
		genres = append(genres, genre)
	}
	if err := rows.Err(); err != nil {
		return nil, false, err
	}

	return genres, true, nil
}
