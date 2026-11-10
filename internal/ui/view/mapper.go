package view

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/tinkyrain/go-sysmon/internal/metrics"
	"github.com/tinkyrain/go-sysmon/internal/metrics/cpu"
	"github.com/tinkyrain/go-sysmon/internal/metrics/disk"
	"github.com/tinkyrain/go-sysmon/internal/metrics/memory"
	"github.com/tinkyrain/go-sysmon/internal/metrics/system"
)

func NewRenderData(snapshot metrics.Snapshot) RenderData {
	return RenderData{
		Time:   snapshot.Time.Format(timeLayout),
		System: systemFrom(snapshot.System),
		CPU:    cpuFrom(snapshot.CPU),
		Memory: memoryFrom(snapshot.Memory),
		Disk:   diskFrom(snapshot.Disk),
	}
}

func systemFrom(stats system.Stats) System {
	return System{
		Info: []Field{
			{Name: "Hostname", Value: stats.Info.Hostname},
			{Name: "OS", Value: stats.Info.OS},
			{Name: "Kernel", Value: stats.Info.Kernel},
			{Name: "Arch", Value: stats.Info.Arch},
		},
		LoadAvg: []Field{
			{Name: "One min", Value: fmt.Sprintf("%.2f", stats.LoadAvg.OneMin)},
			{Name: "Five min", Value: fmt.Sprintf("%.2f", stats.LoadAvg.FiveMin)},
			{Name: "Fifteen min", Value: fmt.Sprintf("%.2f", stats.LoadAvg.FifteenMin)},
			{Name: "Running procs", Value: strconv.FormatUint(stats.LoadAvg.RunningProcs, 10)},
			{Name: "Total procs", Value: strconv.FormatUint(stats.LoadAvg.TotalProcs, 10)},
			{Name: "Last run PID", Value: strconv.FormatUint(stats.LoadAvg.LastPID, 10)},
		},
		Uptime: Uptime{
			Total: HumanizeDuration(stats.Uptime.Total),
			Idle:  HumanizeDuration(stats.Uptime.Idle),
		},
	}
}

func cpuFrom(stats cpu.Stats) CPU {
	cores := make([]cpu.CoreUsage, len(stats.Usage.Cores))
	copy(cores, stats.Usage.Cores)

	sort.SliceStable(cores, func(i, j int) bool {
		if cores[i].Usage != cores[j].Usage {
			return cores[i].Usage > cores[j].Usage
		}
		return cores[i].ID < cores[j].ID
	})

	gauges := make([]Gauge, 0, len(cores))
	for _, core := range cores {
		gauges = append(gauges, Gauge{Label: core.ID, Percent: int(core.Usage)})
	}

	return CPU{
		Total: Gauge{Label: "Total", Percent: int(stats.Usage.Total)},
		Cores: gauges,
	}
}

func memoryFrom(stats memory.Stats) Memory {
	info := stats.Info

	return Memory{
		Info: []Field{
			{Name: "Total", Value: HumanizeSize(float64(info.Total))},
			{Name: "Available", Value: HumanizeSize(float64(info.Available))},
			{Name: "SWAP Total", Value: HumanizeSize(float64(info.SwapTotal))},
			{Name: "SWAP Available", Value: HumanizeSize(float64(info.SwapAvailable))},
		},
		RAM:  Gauge{Label: "Memory", Percent: usedPercent(info.Total, info.Available)},
		Swap: Gauge{Label: "SWAP", Percent: usedPercent(info.SwapTotal, info.SwapAvailable)},
	}
}

func diskFrom(stats disk.Stats) Disk {
	mounts := make([]Mount, 0, len(stats.Mounts))

	for _, mount := range stats.Mounts {
		mounts = append(mounts, Mount{
			Path: mount.Path,
			Info: []Field{
				{Name: "Total", Value: HumanizeSize(float64(mount.Total))},
				{Name: "Available", Value: HumanizeSize(float64(mount.Available))},
			},
		})
	}

	return Disk{Mounts: mounts}
}
