package system

import (
	"errors"
	"fmt"
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

var ErrEmptyFile = errors.New("file is empty")

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

	if len(hostname) == 0 {
		return Info{}, fmt.Errorf("parsing file error %q: %w", hostnameFile, ErrEmptyFile)
	}

	if len(ostype) == 0 {
		return Info{}, fmt.Errorf("parsing file error %q: %w", ostypeFile, ErrEmptyFile)
	}

	if len(osrelease) == 0 {
		return Info{}, fmt.Errorf("parsing file error %q: %w", osreleaseFile, ErrEmptyFile)
	}

	return Info{
		Arch:     runtime.GOARCH,
		Hostname: hostname[0],
		OS:       ostype[0],
		Kernel:   osrelease[0],
	}, nil
}
