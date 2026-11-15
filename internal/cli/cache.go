package cli

import (
	"akerraps/tunectl/internal/cache"
	"fmt"

	urfave "github.com/urfave/cli/v2"
)

func cacheCommand() *urfave.Command {
	return &urfave.Command{
		Name:  "cache",
		Usage: "Manage cache",

		Flags: []urfave.Flag{
			&urfave.BoolFlag{
				Name:  "ytdlp",
				Usage: "Clear the cached yt-dlp binary",
			},
			&urfave.BoolFlag{
				Name:  "db",
				Usage: "Clear the cached metadata database (artists, albums, songs, genres)",
			},
		},

		Action: cacheAction,
	}
}

func cacheAction(c *urfave.Context) error {
	ytdlp := c.Bool("ytdlp")
	db := c.Bool("db")

	if !ytdlp && !db {
		return fmt.Errorf("no valid flag provided, use --ytdlp and/or --db")
	}

	if ytdlp {
		if err := cache.ClearYtDlp(); err != nil {
			return err
		}
	}

	if db {
		if err := cache.ClearDatabase(); err != nil {
			return err
		}
	}

	return nil
}
