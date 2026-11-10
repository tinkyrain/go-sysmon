package procfs

import (
	"bufio"
	"os"
	"path/filepath"
)

type Reader struct {
	path string
}

func New(path string) Reader {
	return Reader{path: path}
}

func (fs Reader) ReadLines(filename string) ([]string, error) {
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

func (fs Reader) ReadLine(filename string) (string, error) {
	f, err := os.Open(filepath.Join(fs.path, filename))
	if err != nil {
		return "", err
	}
	defer func() {
		_ = f.Close()
	}()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		return "", scanner.Err()
	}

	return scanner.Text(), nil
}
