package procfs

import (
	"bufio"
	"os"
	"path/filepath"
)

type FileScanner struct {
	path string
}

func New(path string) FileScanner {
	return FileScanner{path: path}
}

func (fs FileScanner) ScanRows(filename string) ([]string, error) {
	lines := []string{}
	f, err := os.Open(filepath.Join(fs.path, filename))
	if err != nil {
		return lines, err
	}
	defer func() {
		_ = f.Close()
	}()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
