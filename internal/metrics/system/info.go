package system

import (
	"runtime"

	"github.com/tinkyrain/go-sysmon/internal/procfs"
)

const (
	ostypeFile    = "sys/kernel/ostype"
	osreleaseFile = "sys/kernel/osrelease"
	hostnameFile  = "sys/kernel/hostname"
)

type infoReader struct {
	fs procfs.FileScanner
}

type Info struct {
	Hostname string
	OS       string
	Kernel   string
	Arch     string
}

func (r infoReader) read() (Info, error) {
	ostype, err := r.fs.ScanRows(ostypeFile)
	if err != nil {
		return Info{}, err
	}

	osrelease, err := r.fs.ScanRows(osreleaseFile)
	if err != nil {
		return Info{}, err
	}

	hostname, err := r.fs.ScanRows(hostnameFile)
	if err != nil {
		return Info{}, err
	}

	info := Info{}
	info.Arch = runtime.GOARCH
	info.Hostname = "not found"
	info.OS = "not found"
	info.Kernel = "not found"

	if len(hostname) != 0 {
		info.Hostname = hostname[0]
	}

	if len(ostype) != 0 {
		info.OS = ostype[0]
	}

	if len(osrelease) != 0 {
		info.Kernel = osrelease[0]
	}

	return info, nil
}
