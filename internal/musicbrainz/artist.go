package musicbrainz

import (
	"log/slog"
	"time"

	"akerraps/tunectl/internal/types"

	"github.com/michiwend/gomusicbrainz"
)

func getArtistInfo(info *types.TrackInfo) error {
	log := slog.With(
		"track_id", info.ID,
		"title", info.Title,
		"artists", info.Artists,
	)

	log.Debug(
		"fetching artist information",
		"artist_id", info.ArtistMBID,
	)

	client := createClient()

	var resp *gomusicbrainz.ArtistSearchResponse
	var err error

	const maxAttempts = 3

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		resp, err = client.SearchArtist("arid:"+info.ArtistMBID, 1, 0)

		if err == nil {
			break
		}

		log.Warn(
			"artist lookup failed, retrying",
			"attempt", attempt,
			"max_attempts", maxAttempts,
			"artist_id", info.ArtistMBID,
			"err", err,
		)

		time.Sleep(time.Duration(attempt) * 5 * time.Second)
	}

	if err != nil {
		log.Error(
			"artist lookup failed",
			"artist_id", info.ArtistMBID,
			"err", err,
		)

		return err
	}

	extractTags(resp.Artists[0].Tags, info)

	return nil
}
