package musicbrainz

import (
	"fmt"
	"log/slog"
	"time"

	"akerraps/tunectl/internal/types"

	"github.com/michiwend/gomusicbrainz"
)

// ArtistGenreLookup resolves genres already known for an artist, so callers
// can skip the MusicBrainz artist lookup when the data is already cached.
type ArtistGenreLookup func(artist string) (genres []string, found bool, err error)

func searchQuery(song string, artists []string) (query string) {
	if len(artists) == 0 {
		query = fmt.Sprintf(`recording:"%s"`, song)
	} else {
		query = fmt.Sprintf(`recording:%s AND artist:%s`, song, artists[0])
	}

	return query
}

func enrichTrack(info *types.TrackInfo, resp *gomusicbrainz.RecordingSearchResponse, lookupArtistGenres ArtistGenreLookup) error {
	log := slog.With(
		"track_id", info.ID,
		"title", info.Title,
		"artists", info.Artists,
	)

	recording := resp.Recordings[0]

	info.Artists = info.Artists[:0]
	info.Title = recording.Title

	for _, credit := range recording.ArtistCredit.NameCredits {
		info.Artists = append(info.Artists, credit.Artist.Name)
	}

	info.ArtistMBID = string(
		recording.ArtistCredit.NameCredits[0].Artist.ID,
	)

	info.MBID = string(recording.ID)

	if info.Album == "" {
		if err := getAlbumInfo(info); err != nil {
			return err
		}
	}

	var genres []string
	var found bool

	if lookupArtistGenres != nil {
		var err error

		genres, found, err = lookupArtistGenres(info.Artists[0])
		if err != nil {
			log.Warn("failed to check artist cache", "artist", info.Artists[0], "err", err)
		}
	}

	if found {
		log.Debug(
			"artist found in cache, skipping api lookup",
			"artist", info.Artists[0],
			"genres", genres,
		)
		info.Genres = genres
	} else if err := getArtistInfo(info); err != nil {
		return err
	}

	log.Debug(
		"metadata enriched",
		"title", info.Title,
		"artists", info.Artists,
		"album", info.Album,
		"genres", info.Genres,
		"mbid", info.MBID,
	)

	return nil
}

func GetSongInfo(song *types.TrackInfo, lookupArtistGenres ArtistGenreLookup) error {
	log := slog.With(
		"track_id", song.ID,
		"title", song.Title,
		"artists", song.Artists,
	)

	log.Debug("searching musicbrainz recording")

	client := createClient()

	query := searchQuery(song.Title, song.Artists)

	var resp *gomusicbrainz.RecordingSearchResponse
	var err error

	const maxAttempts = 3

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		resp, err = client.SearchRecording(query, 1, 0)

		if err == nil {
			break
		}

		log.Warn(
			"musicbrainz recording search failed, retrying",
			"attempt", attempt,
			"max_attempts", maxAttempts,
			"query", query,
			"err", err,
		)

		time.Sleep(time.Duration(attempt) * 5 * time.Second)
	}

	if err != nil {
		return fmt.Errorf(
			`cannot fetch "%s - %v" song data: %w`,
			song.Title,
			song.Artists,
			err,
		)
	}

	if len(resp.Recordings) == 0 {
		log.Warn("no information found for song")

		return nil
	}

	log.Debug(
		"recording found",
		"title", resp.Recordings[0].Title,
		"artists", song.Artists,
	)

	if err := enrichTrack(song, resp, lookupArtistGenres); err != nil {
		return err
	}

	return nil
}
