package metrics

import (
	"time"

	"github.com/tinkyrain/go-sysmon/internal/metrics/cpu"
	"github.com/tinkyrain/go-sysmon/internal/metrics/disk"
	"github.com/tinkyrain/go-sysmon/internal/metrics/memory"
	"github.com/tinkyrain/go-sysmon/internal/metrics/system"
)

type Snapshot struct {
	Time   time.Time
	System system.Stats
	CPU    cpu.Stats
	Memory memory.Stats
	Disk   disk.Stats
}
