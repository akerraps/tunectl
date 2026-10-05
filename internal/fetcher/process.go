package fetcher

import (
	"fmt"
	"log/slog"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"

	"akerraps/tunectl/internal/musicbrainz"
	"akerraps/tunectl/internal/types"

	"golang.org/x/text/unicode/norm"
)

func sanitize(s string) string {
	s = norm.NFD.String(s)

	var b strings.Builder

	for _, r := range s {
		if unicode.Is(unicode.Mn, r) {
			continue
		}

		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '_' {
			b.WriteRune(r)
		}
	}

	s = strings.ToLower(strings.TrimSpace(b.String()))

	return strings.ReplaceAll(s, " ", "_")
}

func outputPath(song *types.TrackInfo, outputDir string) string {
	name := song.Title
	if len(song.Artists) > 0 {
		name += "_" + song.Artists[0]
	}

	fileName := sanitize(name) + ".mp3"

	var artist string
	if len(song.Artists) > 0 {
		artist = song.Artists[0]
	}

	album := sanitize(song.Album)

	switch {
	case artist == "":
		return filepath.Join(outputDir, "unknown", fileName)
	case album == "":
		return filepath.Join(outputDir, artist, "unknown", fileName)
	default:
		return filepath.Join(outputDir, artist, album, fileName)
	}
}

func processTrack(bin string, song *types.TrackInfo, opts types.Options) error {
	log := slog.With(
		"track_id", song.ID,
		"title", song.Title,
		"artists", song.Artists,
	)

	output := outputPath(song, opts.OutputDir)

	exists, err := songExists(output)
	if err != nil {
		return err
	}

	// Skip existing songs unless metadata rewriting was requested.
	if exists && !opts.RewriteMetadata {
		log.Info("song already exists, skipping")
		return nil
	}

	if exists {
		log.Debug("song already exists, rewriting metadata")
	}

	// Fetch metadata only when the API is enabled.
	if !opts.NoAPI {
		log.Debug("fetching metadata from api")

		if err := musicbrainz.GetSongInfo(song); err != nil {
			log.Warn(
				"failed to fetch metadata, using original data",
				"err", err,
			)
		}
	}

	// Existing song: only rewrite metadata.
	if exists {
		writeMetadata(output, song)

		log.Info("metadata rewritten")
		return nil
	}

	search := fmt.Sprintf(
		`ytsearch:%s %s song`,
		song.Title,
		strings.Join(song.Artists, " "),
	)

	log.Debug(
		"starting download",
		"query", search,
	)

	cmd := exec.Command(
		bin,
		//"--extractor-args", "youtube:player_client=web_embedded",
		"-x",
		"--restrict-filenames",
		"--quiet",
		"--no-warnings",
		"--embed-thumbnail",
		"--no-playlist",
		"-t", "mp3",
		"--audio-quality", "128K",
		search,
		"-o", output,
	)

	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(
			"download failed",
			"err", err,
			"output", string(cmdOutput),
		)

		return nil
	}

	log.Info("download completed")

	writeMetadata(output, song)

	return nil
}
