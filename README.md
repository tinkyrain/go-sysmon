# Terminal system monitor for Linux
System info, CPU, memory and disk usage on four tabs, one keypress apart.

[![lint-and-test](https://github.com/tinkyrain/go-sysmon/actions/workflows/lint-and-test.yml/badge.svg)](https://github.com/tinkyrain/go-sysmon/actions/workflows/lint-and-test.yml)

<img src="https://github.com/tinkyrain/go-sysmon/blob/main/preview.gif?raw=true" width="900" alt="go-sysmon work example">

## Install

Download the `.deb` for your architecture (amd64 or arm64) from
[Releases](https://github.com/tinkyrain/go-sysmon/releases) and install it:

```bash
go install github.com/tinkyrain/go-sysmon/cmd/go-sysmon@latest   # or
sudo apt install ./go-sysmon_*.deb
```

Or build from source (Go 1.26+):

```bash
git clone https://github.com/tinkyrain/go-sysmon.git
cd go-sysmon
go build -o go-sysmon ./cmd/go-sysmon
```

## Usage

```bash
go-sysmon            # start the dashboard
go-sysmon --version  # print version
```

| Key | Action |
|---|---|
| `←` `→` | previous / next tab |
| `Tab` | next tab, wraps to the first |
| `q` / `Ctrl+C` | quit |

### Tabs

- **System** — hostname, OS, kernel, architecture, load average and uptime
- **CPU** — total load plus the busiest cores, as many as the window fits
- **Memory** — RAM and swap, used percentage and absolute sizes
- **Disk** — total and available space per mounted filesystem

The layout follows the terminal: resize the window and the pages re-fit to it.

## Configuration

Optional. Create `~/.config/go-sysmon/config.toml`:

```toml
interval  = "3s"     # how often metrics refresh
proc_root = "/proc"  # path to procfs
```

The values above are the defaults.

## Testing

```bash
go test -race ./...                        # what CI runs
go test ./... -coverprofile=cover.out      # the table below
go tool cover -func=cover.out | tail -1    # the total
```

Total statement coverage is **35.3%**. In package order, as `go test` reports it:

| Package | Coverage |
|---|---|
| `cmd/go-sysmon` | 0.0% |
| `internal/app` | 100.0% |
| `internal/config` | 100.0% |
| `internal/metrics` | 100.0% |
| `internal/metrics/cpu` | 100.0% |
| `internal/metrics/disk` | 100.0% |
| `internal/metrics/memory` | 100.0% |
| `internal/metrics/system` | 100.0% |
| `internal/procfs` | 100.0% |
| `internal/statfs` | 100.0% |
| `internal/ui` | 0.0% |
| `internal/ui/components` | 0.0% |
| `internal/ui/pages` | 2.2% |
| `internal/ui/view` | 100.0% |
## License

MIT
