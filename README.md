# Terminal system monitor for Linux
CPU, memory and disk usage on one screen.

## Install

Download the `.deb` for your architecture (amd64 or arm64) from
[Releases](https://github.com/tinkyrain/go-sysmon/releases) and install it:

```bash
sudo apt install ./go-sysmon_*_amd64.deb
```

Or build from source (Go 1.24+):

```bash
git clone https://github.com/tinkyrain/go-sysmon.git
cd go-sysmon
go build -o go-sysmon ./cmd/go_sysmon
```

## Usage

```bash
go-sysmon            # start the dashboard, Ctrl+C to quit
go-sysmon --version  # print version
```

## Configuration

Optional. Create `~/.config/go-sysmon/config.toml`:

```toml
interval  = "3s"     # how often metrics refresh
proc_root = "/proc"  # path to procfs
```

The values above are the defaults.

## License

MIT
