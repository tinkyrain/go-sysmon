package ui

import "github.com/tinkyrain/go-sysmon/internal/metrics"

type RenderData struct {
	Time   string
	System System
}

type System struct {
	Info struct {
		Hostname string
		OS       string
		Kernel   string
		Arch     string
	}
	LoadAvg struct {
		OneMin       float64
		FiveMin      float64
		FifteenMin   float64
		RunningProcs uint64
		TotalProcs   uint64
		LastPID      uint64
	}
	Uptime struct {
		Total float64
		Idle  float64
	}
}

func NewRenderData(snapshot metrics.Snapshot) {

}
