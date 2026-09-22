package fetcher

import (
	"errors"
	"log/slog"
	"os"
	"sync"
	"sync/atomic"

	"akerraps/tunectl/internal/cache"
	"akerraps/tunectl/internal/types"
)

func songExists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err
}

func FetchAudio(tracks []types.TrackInfo, opts types.Options) error {
	bin, err := cache.InitCache()
	if err != nil {
		return err
	}

	slog.Debug(
		"using yt-dlp binary",
		"path", bin,
	)

	workers := opts.Parallelism

	total := len(tracks)
	jobs := make(chan types.TrackInfo)

	var wg sync.WaitGroup
	var completed atomic.Int64

	// Assign a stable ID to each track based on its position in the input.
	for i := range tracks {
		tracks[i].ID = i + 1
	}

	// Start the worker pool.
	for i := 0; i < workers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for song := range jobs {
				log := slog.With(
					"track_id", song.ID,
					"title", song.Title,
					"artists", song.Artists,
				)

				log.Debug("processing track")

				if err := processTrack(bin, song, opts); err != nil {
					log.Error(
						"failed to process track",
						"err", err,
					)
				}

				current := int(completed.Add(1))
				remaining := total - current

				log.Info(
					"track progress",
					"current", current,
					"total", total,
					"remaining", remaining,
				)

			}
		}()
	}

	// Send all tracks to the worker pool.
	for _, song := range tracks {
		jobs <- song
	}

	// Signal that there are no more tracks to process.
	close(jobs)

	// Wait for all workers to finish.
	wg.Wait()

	return nil
}
