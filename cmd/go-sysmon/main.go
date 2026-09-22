package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
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
		fmt.Printf("go-sysmon %s (commit %s, built %s)\n", version, commit, date)
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
