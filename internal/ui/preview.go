package ui

import (
	"fmt"
	"go-sysmon/internal/metrics"
	"strconv"
)

func clearTerminal() {
	fmt.Print("\033[H\033[2J")
}

func Preview(snapshot metrics.Snapshot) {
	clearTerminal()

	fmt.Println("======================" + snapshot.Time.Format("2006-01-02 15:04:05") + "======================")

	fmt.Println("--------" + " RAM " + "--------")
	fmt.Println("RAM Total: ", snapshot.Memory.Total)
	fmt.Println("RAM Available: ", snapshot.Memory.Available)
	fmt.Println("SWAP Total: ", snapshot.Memory.SwapTotal)
	fmt.Println("SWAP Available: ", snapshot.Memory.SwapAvailable)

	fmt.Println("--------" + " CPU " + "--------")
	fmt.Println("CPU Usage: " + strconv.FormatFloat(snapshot.CPUUsage, 'f', -1, 64) + "%")

	for _, disk := range snapshot.Disks {
		fmt.Println("-------- Disk: " + disk.Mount + " --------")
		fmt.Println("Disk Total: ", disk.Total)
		fmt.Println("Disk Available: ", disk.Available)
	}
}
