package manager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type cacheFile struct {
	tmpdir   string
	disabled bool
}

func (c cacheFile) get(url string) ([]string, error) {
	if c.disabled {
		return []string{}, errors.New("Disabled cache")
	}

	tmpfn := filepath.Join(c.tmpdir, slug(url))
	file, err := os.Open(tmpfn)
	if err != nil {
		return []string{}, err
	}
	defer file.Close()

	return LinesFromReader(file)
}

func (c cacheFile) set(content []string, url string) error {
	if c.disabled {
		return errors.New("Disabled cache")
	}

	tmpfn := filepath.Join(c.tmpdir, slug(url))
	f, err := os.Create(tmpfn)
	if err != nil {
		if f != nil {
			f.Close()
		}
		return err
	}

	for _, line := range content {
		if _, err = fmt.Fprintln(f, line); err != nil {
			return err
		}
	}

	return nil
}

var (
	cache cacheFile
)

func init() {
	cache = NewCache("dfm")
}

func NewCache(prefix string) cacheFile {
	dir := os.TempDir()
	tmpdir := filepath.Join(dir, prefix)
	if _, err := os.Stat(tmpdir); os.IsNotExist(err) {
		_ = os.MkdirAll(tmpdir, os.ModePerm)
	}

	_, err := os.Stat(tmpdir)
	disabled := os.IsNotExist(err)
	return cacheFile{
		tmpdir:   tmpdir,
		disabled: disabled,
	}
}
