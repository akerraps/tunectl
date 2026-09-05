package musicbrainz

import (
	"encoding/xml"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"akerraps/tunectl/internal/types"
)

type recordingLookupResponse struct {
	Recording struct {
		Releases []release `xml:"release-list>release"`
	} `xml:"recording"`
}

type release struct {
	ID     string `xml:"id,attr"`
	Title  string `xml:"title"`
	Status string `xml:"status"`

	ReleaseGroup struct {
		ID          string `xml:"id,attr"`
		PrimaryType string `xml:"primary-type"`
	} `xml:"release-group"`
}

func getAlbumInfo(info *types.TrackInfo) error {
	log := slog.With(
		"track_id", info.ID,
		"title", info.Title,
		"artists", info.Artists,
		"mbid", info.MBID,
	)

	if info.MBID == "" {
		return fmt.Errorf("recording has no MBID")
	}

	const maxAttempts = 3

	var resp *recordingLookupResponse
	var err error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		resp, err = lookupRecordingReleases(info.MBID)

		if err == nil {
			break
		}

		log.Warn(
			"album lookup failed, retrying",
			"attempt", attempt,
			"max_attempts", maxAttempts,
			"err", err,
		)

		if attempt < maxAttempts {
			time.Sleep(time.Duration(attempt) * 5 * time.Second)
		}
	}

	if err != nil {
		return err
	}

	if len(resp.Recording.Releases) == 0 {
		log.Warn("no releases found")
		return nil
	}

	// Prefer an official album.
	for _, release := range resp.Recording.Releases {
		if release.Status == "Official" &&
			release.ReleaseGroup.PrimaryType == "Album" {

			info.Album = release.Title
			info.AlbumMBID = release.ReleaseGroup.ID

			log.Debug(
				"album found",
				"album", info.Album,
				"album_mbid", info.AlbumMBID,
				"release_id", release.ID,
				"status", release.Status,
			)

			return nil
		}
	}

	// If there is no official album, prefer any album.
	for _, release := range resp.Recording.Releases {
		if release.ReleaseGroup.PrimaryType == "Album" {
			info.Album = release.Title
			info.AlbumMBID = release.ReleaseGroup.ID

			log.Debug(
				"album found",
				"album", info.Album,
				"album_mbid", info.AlbumMBID,
				"release_id", release.ID,
				"status", release.Status,
			)

			return nil
		}
	}

	// Last resort: use the first release.
	firstRelease := resp.Recording.Releases[0]

	info.Album = firstRelease.Title
	info.AlbumMBID = firstRelease.ReleaseGroup.ID

	log.Debug(
		"release found, using release title as album",
		"album", info.Album,
		"album_mbid", info.AlbumMBID,
		"release_id", firstRelease.ID,
		"status", firstRelease.Status,
	)

	return nil
}

func lookupRecordingReleases(mbid string) (*recordingLookupResponse, error) {
	endpoint := fmt.Sprintf(
		"https://musicbrainz.org/ws/2/recording/%s",
		url.PathEscape(mbid),
	)

	query := url.Values{}

	// Request releases together with their release groups.
	query.Set("inc", "releases+release-groups")
	query.Set("fmt", "xml")

	req, err := http.NewRequest(
		http.MethodGet,
		endpoint+"?"+query.Encode(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set(
		"User-Agent",
		"Terminal music downloader & metadata tool/0.0.1-beta (https://github.com/akerraps/tunectl)",
	)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"musicbrainz returned HTTP status %s",
			res.Status,
		)
	}

	var data recordingLookupResponse

	if err := xml.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf(
			"failed to decode MusicBrainz response: %w",
			err,
		)
	}

	return &data, nil
}