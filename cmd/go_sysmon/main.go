package main

import (
	"fmt"
	"go-sysmon/internal/config"
	"go-sysmon/internal/metrics"
	"go-sysmon/internal/ui"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "go-sysmon error:", err)
		os.Exit(1)
	}
}

func run() error {
	configFile, err := config.DefaultPath()
	if err != nil {
		return err
	}
	cfg, err := config.Load(configFile)
	if err != nil {
		return err
	}

	collector := metrics.New(cfg.ProcRoot)
	snapshot, err := collector.Collect()
	if err != nil {
		return err
	}
	ui.Preview(snapshot)

	return nil
}
