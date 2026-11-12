package view

type Field struct {
	Name  string
	Value string
}

type Gauge struct {
	Label   string
	Percent int
}

type RenderData struct {
	Time      string
	Timestamp int64
	System    System
	CPU       CPU
	Memory    Memory
	Disk      Disk
}

type Uptime struct {
	Total string
	Idle  string
}

type System struct {
	Info    []Field
	LoadAvg []Field
	Uptime  Uptime
}

type CPU struct {
	Total Gauge
	Cores []Gauge
}

type Memory struct {
	Info []Field
	RAM  Gauge
	Swap Gauge
}

type Mount struct {
	Path string
	Info []Field
}

type Disk struct {
	Mounts []Mount
}
