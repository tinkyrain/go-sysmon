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
	info.Arch = runtime.GOARCH

	fileToValue := map[string]*string{
		hostnameFile:  &info.Hostname,
		ostypeFile:    &info.OS,
		osreleaseFile: &info.Kernel,
	}

	for file, value := range fileToValue {
		row, err := r.fs.ScanRow(file)

		if err != nil {
			return Info{}, fmt.Errorf("scan %q: %w", file, err)
		}

		if len(row) == 0 {
			return Info{}, fmt.Errorf("read %q: %w", file, ErrEmptyFile)
		}

		*value = row
	}

	return info, nil
}
