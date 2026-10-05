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
	info := Info{}
	var err error

	info.Arch = runtime.GOARCH

	info.OS, err = r.readRow(ostypeFile)
	if err != nil {
		return Info{}, err
	}

	info.Kernel, err = r.readRow(osreleaseFile)
	if err != nil {
		return Info{}, err
	}

	info.Hostname, err = r.readRow(hostnameFile)
	if err != nil {
		return Info{}, err
	}

	return info, nil
}

func (r infoReader) readRow(file string) (string, error) {
	row, err := r.fs.ScanRow(file)
	if err != nil {
		return "", fmt.Errorf("scan %q: %w", file, err)
	}

	if len(row) == 0 {
		return "", fmt.Errorf("read %q: %w", file, ErrEmptyFile)
	}

	return row, nil
}
