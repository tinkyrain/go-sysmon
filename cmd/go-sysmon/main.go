package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/tinkyrain/go-sysmon/internal/app"
	"github.com/tinkyrain/go-sysmon/internal/config"
	"github.com/tinkyrain/go-sysmon/internal/ui"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()
	if *showVersion {
		fmt.Println(versionString())
		return
	}

	// Context with stop notify
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop() // stop needle for collect context resources

	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "go-sysmon error:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	configFile, err := config.DefaultPath()
	if err != nil {
		return err
	}

	cfg, err := config.Load(configFile)
	if err != nil {
		return err
	}

	dash := ui.NewDashboard()
	defer dash.Stop()

	return app.Run(ctx, cfg, dash.Render)
}

func versionString() string {
	v, c, d := version, commit, date
	if v == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok {
			if info.Main.Version != "" && info.Main.Version != "(devel)" {
				v = info.Main.Version
			}
			for _, s := range info.Settings {
				switch s.Key {
				case "vcs.revision":
					c = s.Value
				case "vcs.time":
					d = s.Value
				}
			}
		}
	}
	return fmt.Sprintf("go-sysmon %s (commit %s, built %s)", v, c, d)
}
