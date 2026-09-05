package database

import (
	"akerraps/tunectl/internal/cache"
	"akerraps/tunectl/internal/types"
	"context"
)

func WriteSong (info *types.TrackInfo) (err error) {
	
	db,err:=cache.OpenDatabase()
	if err != nil {
		return err
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO songs (title, mbid) VALUES (?,?);`,
		info.Title, info.MBID)
	if err != nil {
		return err
	}
	
	return nil
}