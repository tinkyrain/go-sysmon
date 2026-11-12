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
	fs procfs.Reader
}

type Info struct {
	Hostname string
	OS       string
	Kernel   string
	Arch     string
}

var ErrEmptyInfoData = errors.New("info data is empty")

func (r infoReader) read() (Info, error) {
	info := Info{}
	info.Arch = runtime.GOARCH

	fileToValue := map[string]*string{
		hostnameFile:  &info.Hostname,
		ostypeFile:    &info.OS,
		osreleaseFile: &info.Kernel,
	}

	for file, value := range fileToValue {
		line, err := r.fs.ReadLine(file)

		if err != nil {
			return Info{}, fmt.Errorf("read %q: %w", file, err)
		}

		if len(line) == 0 {
			return Info{}, fmt.Errorf("read %q: %w", file, ErrEmptyInfoData)
		}

		*value = line
	}

	return info, nil
}
