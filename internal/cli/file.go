package cli

import (
	"akerraps/tunectl/internal/fetcher"
	"akerraps/tunectl/internal/reader"
	"akerraps/tunectl/internal/types"
	"fmt"

	urfave "github.com/urfave/cli/v2"
)

func fileCommand() *urfave.Command {
	return &urfave.Command{
		Name:  "file",
		Usage: "Download songs from a file",

		Flags: []urfave.Flag{
			&urfave.BoolFlag{
				Name:               "csv",
				Usage:              "Read tracks from CSV",
				DisableDefaultText: true,
			},
			&urfave.BoolFlag{
				Name:               "json",
				Usage:              "Read tracks from JSON",
				DisableDefaultText: true,
			},
			&urfave.StringFlag{
				Name:    "data",
				Aliases: []string{"d"},
				Usage:   "Input file path",
			},
			outputFlag(),
			envFlag(),
			noAPIFlag(),
			rewriteMetadataFlag(),
			parallelFlag(),
		},

		Action: fileAction,
	}
}

func fileAction(c *urfave.Context) error {
	opts := buildOptions(c)

	file := c.String("data")
	if file == "" {
		return fmt.Errorf("you must specify a file with --data")
	}

	tracks, err := loadTracks(c, file)
	if err != nil {
		return fmt.Errorf("failed to read input file %s: %w", file, err)
	}

	return fetcher.FetchAudio(tracks, opts)
}

func loadTracks(c *urfave.Context, file string) ([]types.TrackInfo, error) {
	csv := c.Bool("csv")
	json := c.Bool("json")

	if csv == json {
		return nil, fmt.Errorf("you must choose exactly one file type (--csv or --json)")
	}

	if csv {
		return reader.ReadCsvFile(file)
	}

	return reader.ReadJsonFile(file)
}
