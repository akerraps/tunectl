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
		ctx := context.Background()

		artistIDs := make([]int64, 0, len(song.Artists))
		for _, artist := range song.Artists {
			var id int64
			err := tx.QueryRowContext(ctx,
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
			err := tx.QueryRowContext(ctx,
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
				if _, err := tx.ExecContext(ctx,
					`INSERT INTO artist_genres (artist_id, genre_id) VALUES (?, ?)
					 ON CONFLICT DO NOTHING;`,
					artistID, genreID); err != nil {
					return fmt.Errorf("link artist %d to genre %d: %w", artistID, genreID, err)
				}
			}
		}

		var albumArtistID sql.NullInt64
		if len(artistIDs) > 0 {
			albumArtistID = sql.NullInt64{Int64: artistIDs[0], Valid: true}
		}
		albumMBID := sql.NullString{String: song.AlbumMBID, Valid: song.AlbumMBID != ""}

		var albumID int64
		err := tx.QueryRowContext(ctx,
			`INSERT INTO albums (name, album_mbid, artist_id) VALUES (?, ?, ?)
			 ON CONFLICT DO UPDATE SET
			     name = excluded.name,
			     album_mbid = COALESCE(excluded.album_mbid, albums.album_mbid)
			 RETURNING id;`,
			song.Album, albumMBID, albumArtistID).Scan(&albumID)
		if err != nil {
			return fmt.Errorf("upsert album %q: %w", song.Album, err)
		}

		songMBID := sql.NullString{String: song.MBID, Valid: song.MBID != ""}

		var songID int64
		err = tx.QueryRowContext(ctx,
			`INSERT INTO songs (title, mbid, album_id) VALUES (?, ?, ?)
			 ON CONFLICT DO UPDATE SET
			     title = excluded.title,
			     mbid = COALESCE(excluded.mbid, songs.mbid),
			     album_id = excluded.album_id
			 RETURNING id;`,
			song.Title, songMBID, albumID).Scan(&songID)
		if err != nil {
			return fmt.Errorf("upsert song %q: %w", song.Title, err)
		}

		for _, artistID := range artistIDs {
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO song_artists (song_id, artist_id) VALUES (?, ?)
				 ON CONFLICT DO NOTHING;`,
				songID, artistID); err != nil {
				return fmt.Errorf("link song %d to artist %d: %w", songID, artistID, err)
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
