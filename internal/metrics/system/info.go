package system

import (
	"runtime"

	"github.com/tinkyrain/go-sysmon/internal/procfs"
)

const (
	osTypeFilepath    = "sys/kernel/ostype"
	osReleaseFilepath = "sys/kernel/osrelease"
	hostNameFilepath  = "sys/kernel/hostname"
)

type InfoReader struct {
	fs procfs.FileScanner
}

type Info struct {
	Hostname string
	OS       string
	Kernel   string
	Arch     string
}

func (r *InfoReader) Read() (Info, error) {
	osTypeFile, err := r.fs.ScanRows(osTypeFilepath)
	if err != nil {
		return Info{}, err
	}

	osReleaseFile, err := r.fs.ScanRows(osReleaseFilepath)
	if err != nil {
		return Info{}, err
	}

	hostnameFile, err := r.fs.ScanRows(hostNameFilepath)
	if err != nil {
		return Info{}, err
	}

	info := Info{}
	info.Arch = runtime.GOARCH

	if len(hostnameFile) == 0 && hostnameFile[0] == "" {
		info.Hostname = "not found"
	} else {
		info.Hostname = hostnameFile[0]
	}

	if len(osTypeFile) == 0 && osTypeFile[0] == "" {
		info.OS = "not found"
	} else {
		info.OS = osTypeFile[0]
	}

	if len(osReleaseFile) == 0 && osReleaseFile[0] == "" {
		info.Kernel = "not found"
	} else {
		info.Kernel = osReleaseFile[0]
	}

	return info, nil
}
