package database

import (
	"akerraps/tunectl/internal/cache"
	"akerraps/tunectl/internal/types"
	"context"
)

func WriteSong(info *types.TrackInfo) (err error) {

	db, err := cache.OpenDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO songs (title, mbid, artist_mbid, album_mbid) VALUES (?,?,?,?);`,
		info.Title, nullIfEmpty(info.MBID), nullIfEmpty(info.ArtistMBID), nullIfEmpty(info.AlbumMBID))
	if err != nil {
		return err
	}

	return nil
}
