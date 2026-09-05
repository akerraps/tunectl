package database

import (
	"akerraps/tunectl/internal/types"
	"log/slog"
)

func WriteDatabase(song *types.TrackInfo) {
	log := slog.With(
		"track_id", song.ID,
		"title", song.Title,
		"artists", song.Artists,
	)
	err := WriteSong(song)
	if err != nil {
		log.Error(
			"failed to write song to database",
			"err", err,
		)
	}

	err = WriteArtist(song)
	if err != nil {
		log.Error(
			"failed to write artist to database",
			"err", err,
		)
	}
}