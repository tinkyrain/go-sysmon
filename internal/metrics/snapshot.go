package metrics

import "time"

type Snapshot struct {
	Time      time.Time
	Memory    Memory
	Disks     []Disk
	CPUUsages []CPUUsage
}
