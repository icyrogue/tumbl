package tumbl

import (
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	getter "github.com/hashicorp/go-getter"
)

type puller struct {
	link    string
	mtx     *sync.RWMutex
	options *Options
}

func NewPuller(link string, options *Options) *puller {
	return &puller{
		options: options,
		link:    link,
		mtx:     &sync.RWMutex{},
	}
}

func (p *puller) SetURL(link string) error {
	_, err := url.Parse(link)
	if err != nil {
		return err
	}
	if p.link == link {
		return nil
	}
	p.mtx.Lock()
	defer p.mtx.Unlock()
	p.link = link
	return os.RemoveAll(p.options.Dst)
}

func (p *puller) Pull() ([]string, error) {
	p.mtx.RLock()
	defer p.mtx.RUnlock()
	if p.options == nil {
		return nil, fmt.Errorf("Pull: got nil dst option")
	}
	if err := getter.Get(p.options.Dst, p.link); err != nil {
		return nil, err
	}
	return findExecutable(p.options.Dst)
}

const exts = ".sh,.bash,.py"

func findExecutable(dir string) (executables []string, err error) {
	err = filepath.WalkDir(dir, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info.Type()
		if info.IsDir() && strings.Contains(info.Name(), ".") {
			return fs.SkipDir
		}
		if !strings.Contains(exts, filepath.Ext(info.Name())) {
			return nil
		}
		executables = append(executables, filepath.Base(path))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return executables, nil
}
