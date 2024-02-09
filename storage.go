package tumbl

import (
	"errors"
	"os"
	"path/filepath"
)

type storage struct {
	options *Options
}

func NewStorage(opts *Options) (_ *storage, err error) {
	if opts == nil {
		return nil, errors.New("storage: got nil options")
	}
	if err = os.MkdirAll(opts.Dst, 0777); err != nil {
		return nil, err
	}
	return &storage{
		options: opts,
	}, nil
}

func (s *storage) Files() (files []string, err error) {
	return filepath.Glob(s.options.Dst)
}
