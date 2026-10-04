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
	f, err := os.Open(filepath.Join(fs.path, filename))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	lines := []string{}
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
