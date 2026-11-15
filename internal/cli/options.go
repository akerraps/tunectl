package cli

import (
	"akerraps/tunectl/internal/types"
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	urfave "github.com/urfave/cli/v2"
)

func buildOptions(c *urfave.Context) types.Options {
	opts := types.DefaultOptions()

	if c.IsSet("environment") {
		opts.EnvFile = c.String("environment")

		slog.Debug(
			"overriding environment file with CLI option",
			"env_file", opts.EnvFile,
		)
	}

	readEnv(&opts)

	if out := c.String("output"); out != "" {
		opts.OutputDir = out

		slog.Debug(
			"overriding output directory with CLI option",
			"output_dir", opts.OutputDir,
		)
	}

	if c.IsSet("parallelism") {
		opts.Parallelism = c.Int("parallelism")

		slog.Debug(
			"overriding parallelism with CLI option",
			"parallelism", opts.Parallelism,
		)
	}

	if c.IsSet("no-api") {
		opts.NoAPI = c.Bool("no-api")

		slog.Debug(
			"overriding no-api with CLI option",
			"no_api", opts.NoAPI,
		)
	}

	if c.IsSet("rewrite-metadata") {
		opts.RewriteMetadata = c.Bool("rewrite-metadata")

		slog.Debug(
			"overriding rewrite-metadata with CLI option",
			"rewrite-metadata", opts.RewriteMetadata,
		)
	}

	slog.Debug(
		"final options",
		"output_dir", opts.OutputDir,
		"no_api", opts.NoAPI,
		"rewrite-metadata", opts.RewriteMetadata,
		"parallelism", opts.Parallelism,
	)

	return opts
}

func readEnv(opts *types.Options) {
	if err := godotenv.Overload(opts.EnvFile); err != nil {
		slog.Debug(
			"failed to load environment file",
			"file", opts.EnvFile,
			"error", err,
		)
		return
	}

	if err := envconfig.Process("TUNECTL", opts); err != nil {
		slog.Debug(
			"failed to process environment configuration",
			"error", err,
		)
	}
}
