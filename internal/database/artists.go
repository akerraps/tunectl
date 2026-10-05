package database

import (
	"akerraps/tunectl/internal/types"
	"context"
	"database/sql"
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
