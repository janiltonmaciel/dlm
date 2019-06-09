package manager

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type cacheFile struct {
	tmpdir   string
	disabled bool
}

func (c cacheFile) get(url string, ttl ...time.Duration) ([]byte, error) {
	if c.disabled {
		return nil, errors.New("Disabled cache")
	}

	tmpfn := filepath.Join(c.tmpdir, slug(url))
	file, err := os.Open(tmpfn)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if len(ttl) > 0 {
		if stat, err := file.Stat(); err != nil {
			return nil, err
		} else {
			elapsed := time.Since(stat.ModTime())
			// fmt.Printf("File time: %+v - elapsed: %+v", stat.ModTime(), elapsed)
			if elapsed > ttl[0] {
				return nil, err
			}
		}
	}

	return ioutil.ReadAll(file)
}

func (c cacheFile) set(content []byte, url string) error {
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

	if _, err = fmt.Fprintln(f, string(content)); err != nil {
		return err
	}

	return nil
}

var (
	cache cacheFile
)

func init() {
	cache = NewCache("dlm")
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
