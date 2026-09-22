package database

import (
	"akerraps/tunectl/internal/cache"
	"akerraps/tunectl/internal/types"
	"context"
	"database/sql"
	"fmt"
)

func WriteDatabase(song *types.TrackInfo) error {
	return ExecTx(context.Background(), func(tx *sql.Tx) error {

		artistIDs := make([]int64, 0, len(song.Artists))
		for _, artist := range song.Artists {
			var id int64
			err := tx.QueryRowContext(context.Background(),
				`INSERT INTO artists (name) VALUES (?)
				 ON CONFLICT(name) DO UPDATE SET name = excluded.name
				 RETURNING id;`,
				artist).Scan(&id)
			if err != nil {
				return fmt.Errorf("upsert artist %q: %w", artist, err)
			}
			artistIDs = append(artistIDs, id)
		}

		genreIDs := make([]int64, 0, len(song.Genres))
		for _, genre := range song.Genres {
			var id int64
			err := tx.QueryRowContext(context.Background(),
				`INSERT INTO genres (name) VALUES (?)
				 ON CONFLICT(name) DO UPDATE SET name = excluded.name
				 RETURNING id;`,
				genre).Scan(&id)
			if err != nil {
				return fmt.Errorf("upsert genre %q: %w", genre, err)
			}
			genreIDs = append(genreIDs, id)
		}

		for _, artistID := range artistIDs {
			for _, genreID := range genreIDs {
				_, err := tx.ExecContext(context.Background(),
					`INSERT INTO artist_genres (artist_id, genre_id) VALUES (?, ?)
					 ON CONFLICT DO NOTHING;`,
					artistID, genreID)
				if err != nil {
					return fmt.Errorf("link artist %d to genre %d: %w", artistID, genreID, err)
				}
			}
		}

		return nil
	})
}

func nullIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func ExecTx(ctx context.Context, fn func(tx *sql.Tx) error) (err error) {
	db, err := cache.OpenDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	err = fn(tx)
	return err
}
