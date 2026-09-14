package ui

import (
	"fmt"
	"go-sysmon/internal/metrics"
)

func Preview(snapshot metrics.Snapshot) {
	fmt.Println("======================" + snapshot.Time.Format("2006-01-02 15:04:05") + "======================")

	fmt.Println("--------" + " RAM " + "--------")
	fmt.Println("RAM Total: ", snapshot.Memory.Total)
	fmt.Println("RAM Available: ", snapshot.Memory.Available)
	fmt.Println("SWAP Total: ", snapshot.Memory.SwapTotal)
	fmt.Println("SWAP Available: ", snapshot.Memory.SwapAvailable)

	fmt.Println("--------" + " CPU " + "--------")
	fmt.Println("CPU Usage: ", snapshot.CPUUsage)

	for _, disk := range snapshot.Disks {
		fmt.Println("-------- Disk: " + disk.Mount + " --------")
		fmt.Println("Disk Total: ", disk.Total)
		fmt.Println("Disk Available: ", disk.Available)
	}
}
