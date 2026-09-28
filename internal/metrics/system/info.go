package system

// const (
// 	osTypeFilepath    = "sys/kernel/ostype"
// 	osReleaseFilepath = "sys/kernel/osrelease"
// 	hostNameFilepath  = "sys/kernel/hostname"
// )

type InfoReader struct {
	// fs procfs.FileScanner
}

type Info struct {
	Hostname string
	OS       string
	Kernel   string
	Arch     string
}

func (r *InfoReader) Read() (Info, error) {
	// 	osTypeFile, err := r.fs.ScanRows(osTypeFilepath)
	// 	if err != nil {
	// 		return Info{}, err
	// 	}
	// 	osReleaseFile, err := r.fs.ScanRows(osReleaseFilepath)
	// 	if err != nil {
	// 		return Info{}, err
	// 	}
	// 	hostnameFile, err := r.fs.ScanRows(hostNameFilepath)
	// 	if err != nil {
	// 		return Info{}, err
	// 	}
	//
	// 	info := Info{}
	// 	info.Arch = runtime.GOARCH

	return Info{}, nil
}
