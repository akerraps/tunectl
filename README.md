# tunectl

Terminal music downloader written in Go.

Uses `yt-dlp` for downloading and MusicBrainz for optional metadata lookup.

## Installation

Build the binary:

```bash
    go build -o tunectl ./cmd/tunectl
```

Or run it directly:

```bash
    go run ./cmd/tunectl
```

### Dependencies

#### Development

Required to build and develop TuneCTL:

* `go`
* `gcc`
* `gopls`
* `delve`

#### Runtime

Required by `yt-dlp` to download and process media:

* `python3`
* `ffmpeg`
* `nodejs`
* `deno`

## Run with Docker

TuneCTL can also be built and run using Docker.

The provided `Dockerfile` uses a multi-stage build. The first stage clones the repository and compiles TuneCTL, while the second stage creates the runtime image with the dependencies required by `yt-dlp`.

Build the image locally:

```bash
docker build --no-cache -t tunectl .
```

The repository and version can be customized using build arguments:

```bash
docker build \
  --build-arg TUNECTL_VERSION=main \
  -t tunectl .
```

By default, the Dockerfile builds the main branch of the configured repository.

### Run TuneCTL

```bash
docker run --rm tunectl songs "Alchemy;Philip Sayce"
```

To persist downloaded music, mount an output directory:

```bash
docker run --rm \
  -v "$(pwd):/out" \
  tunectl songs -o /out "Alchemy;Philip Sayce"
```

Input files can also be mounted into the container:

```bash
docker run --rm \
  -v "$(pwd)/music:/out" \
  -v "$(pwd)/tunectl-files:/input" \
  tunectl file --csv --data /input/ -o /out
```

The Docker image runs TuneCTL as a non-root user.

## Configuration

Configuration can be provided through an environment file and CLI options.

By default, tunectl looks for a .env file in the current working directory:

```env
    TUNECTL_OUTPUT_DIR=/home/user/Music
    TUNECTL_NO_API=false
    TUNECTL_REWRITE_METADATA=false
    TUNECTL_PARALLELISM=1
```

A different environment file can be specified with the `--environment` (`-e`) option:

```bash
    tunectl songs -e ~/.env.tunectl "Title;Artist"
```

CLI options override values loaded from the environment file.

Configuration priority:

1- Default values
2- Environment file
3- CLI options

### Options

| Option               | Alias | Description                                          | Default   |
| -------------------- | ----- | ---------------------------------------------------- | --------- |
| `--environment`      | `-e`  | Path to the environment file                         | `.env`    |
| `--output`           | `-o`  | Output directory                                     | `~/Music` |
| `--no-api`           |       | Disable MusicBrainz metadata lookup                  | `false`   |
| `--rewrite-metadata` |       | Rewrite metadata of existing songs using MusicBrainz | `false`   |
| `--paralelism`       | `-n`  | Number of parallel downloads                         | `1`       |

## Download songs

The main command is `songs`.

```bash
    tunectl songs "Title;Artist"
```

A genre can optionally be specified:

```bash
    tunectl songs "Title;Artist;Genre"
```

Multiple artists can be specified as a comma-separated list:

```bash
    tunectl songs "Title;Artist 1, Artist 2"
```

Multiple songs can be downloaded in a single command:

```bash
    tunectl songs \
      "Song 1;Artist 1" \
      "Song 2;Artist 2" \
      "Song 3;Artist 3"
```

### Examples

Download a song:

```bash
    tunectl songs "Alchemy;Philip Sayce"
```

Specify a genre:

```bash
    tunectl songs "Alchemy;Philip Sayce;Blues"
```

Specify an output directory:

```bash
    tunectl songs -o ~/Music "Alchemy;Philip Sayce"
```

Download multiple tracks in parallel:

```bash
    tunectl songs -n 4 \
      "Song 1;Artist 1" \
      "Song 2;Artist 2"
```

## Metadata lookup

By default, `tunectl` uses the MusicBrainz API to resolve track metadata before downloading.

For example:

```bash
    tunectl songs "Alchemy;Philip sayce"
```

MusicBrainz can resolve the artist and track information before generating the download query.

To disable metadata lookup:

```bash
    tunectl songs --no-api "Alchemy;Philip Sayce"
```

When the API is disabled, the original input is used directly.

> [!WARNING]
> MusicBrainz resolution is not guaranteed to return the intended recording. An incorrect match can result in downloading a different track.

### Existing songs

By default, `tunectl` skips songs that already exist in the output directory.

To rewrite the metadata of an existing song, use `--rewrite-metadata`:

```bash
    tunectl songs --rewrite-metadata "Alchemy;Philip Sayce"
```

## File input

Tracks can also be loaded from CSV or JSON files using the `file` command.

The `--data` option accepts either a **single file** or a **directory**. When a directory is provided, all files with the corresponding extension are processed.

### CSV

Single file:

```bash
    tunectl file --csv --data songs.csv
```

Directory:

```bash
    tunectl file --csv --data ./songs/
```

Example:

```csv
Title,Artists,Album,Genre
What the Funk,"Gustavo Mota, Naizon",What the Funk,Electronic
```

The required columns are:

- `Title`
- `Artists`
- `Album`
- `Genre`

Multiple artists and genres can be separated by commas.

Additional columns are ignored.

### JSON

Single file:

```bash
    tunectl file --json --data songs.json
```

Directory:

```bash
    tunectl file --json --data ./songs/
```

Example:

```json
[
  {
    "Title": "Title",
    "Artists": ["Artist 1", "Artist 2"],
    "Album": "Album",
    "Genre": "Rock"
  }
]
```

`Genre` is optional and is used as part of the track metadata.

## Cache

`tunectl` automatically downloads and caches the `yt-dlp` binary.

Clear the cached binary with:

```bash
    tunectl cache --clear
```

The next execution will download `yt-dlp` again.

This can be useful if the cached binary is corrupted or outdated.

## Debugging

Enable debug logs with:

```bash
    tunectl --debug songs "Alchemy;Philip Sayce"
```

This displays information about configuration, metadata lookup, downloads, and other internal operations.

## Command reference

### `songs`

Download one or more songs.

```
    tunectl songs [OPTIONS] SONG...
```

### `file`

Download songs from a CSV or JSON file.

```bash
    tunectl file --csv --data FILE
    tunectl file --json --data FILE
```

### `cache`

Manage the local `yt-dlp` cache.

```bash
    tunectl cache --clear
```
