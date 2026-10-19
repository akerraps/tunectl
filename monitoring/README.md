# Tunectl Grafana Dashboard

Grafana dashboard for monitoring the `tunectl` container: download progress, errors, MusicBrainz retries, logs, and CPU / memory / network usage.

## Requirements

- The container must be named **`tunectl`**.
- **Prometheus** with container metrics (cAdvisor) that have the `name` and `instance` labels.
- **Loki** with the container logs, labeled `container` and `instance`.
- **Grafana** with a Prometheus and a Loki data source (recent version, the dashboard uses the v2 schema).

[Alloy](https://grafana.com/docs/alloy/) is the recommended way to ship both metrics and logs.

## Labels

| Source | Label | Value |
|---|---|---|
| Prometheus (cAdvisor) | `name` | `tunectl` |
| Prometheus | `instance` | your host name |
| Loki | `container` | `tunectl` |
| Loki | `instance` | **same value** as in Prometheus |

The `Host` dropdown is built from the `instance` label of the `container_*` metrics, and the same value filters the Loki queries. If they differ, the log panels will be empty.

## Setup

1. Run Prometheus with remote write enabled: `--web.enable-remote-write-receiver`.
2. Run Loki and Grafana, and add both as data sources.
3. Run Alloy with `/var/run/docker.sock` mounted, and configure it to:
   - collect cAdvisor metrics and send them to Prometheus (`prometheus.remote_write`),
   - collect Docker logs and send them to Loki (`loki.source.docker` + `loki.write`),
   - set the same `instance` label on both metrics and logs.
4. In Grafana: **Dashboards → New → Import**, and upload the dashboard JSON.
5. Pick your host in the `Host` dropdown.

## Notes

- Select a specific host, not "All". The Loki queries use an exact match on `instance`.
- `tunectl` must log in logfmt with an uppercase `level` (`level=ERROR`). The stats and error panels parse these fields.
- Docker log files are the source, so make sure the container uses the default `json-file` or `local` logging driver.