package database

import (
	"akerraps/tunectl/internal/cache"
	"akerraps/tunectl/internal/types"
	"context"
)

func WriteArtist (song *types.TrackInfo) (err error) {
	
	db,err:=cache.OpenDatabase()
	if err != nil {
		return  err
	}

	for _, artist := range song.Artists {
		_, err = db.ExecContext(context.Background(),
			`INSERT INTO artists (name) VALUES (?);`,
			artist)
		if err != nil {
			return err
		}
	}
	
	return nil
}