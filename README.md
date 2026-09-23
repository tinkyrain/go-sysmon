# Terminal system monitor for Linux
CPU, memory and disk usage on one screen.

<img src="https://github.com/tinkyrain/go-sysmon/blob/main/preview.gif?raw=true" width="900" alt="go-sysmon work example">

## Install

Download the `.deb` for your architecture (amd64 or arm64) from
[Releases](https://github.com/tinkyrain/go-sysmon/releases) and install it:

```bash
go install github.com/tinkyrain/go-sysmon/cmd/go-sysmon@latest   # or
sudo apt install ./go-sysmon_*.deb
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
