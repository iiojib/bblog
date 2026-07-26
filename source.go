package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"syscall"
)

const permissions = 0666

type sourceType uint

const (
	invalid sourceType = iota
	stdin
	fifo
)

type source struct {
	io.ReadCloser

	path      string
	kind      sourceType
	own       bool
	closeOnce sync.Once
	closeErr  error
}

func (s *source) Close() error {
	s.closeOnce.Do(func() {
		if s.own {
			cleanup(s.path)
		}

		s.closeErr = s.ReadCloser.Close()
	})

	return s.closeErr
}

func (s *source) String() string {
	switch s.kind {
	case stdin:
		return "stdin"
	case fifo:
		return s.path
	default:
		return "invalid source"
	}
}

func cleanup(p string) {
	if err := os.Remove(p); err != nil {
		log.Printf("remove %s: %v", p, err)
	}
}

func newSource(path string) (*source, error) {
	if path == "" {
		return &source{
			ReadCloser: os.Stdin,
			kind:       stdin,
			own:        false,
		}, nil
	}

	own := false
	info, err := os.Stat(path)

	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("cannot access %s: %v", path, err)
		}

		if mkErr := syscall.Mkfifo(path, permissions); mkErr != nil {
			return nil, fmt.Errorf("mkfifo %s: %v", path, mkErr)
		}

		if chErr := os.Chmod(path, permissions); chErr != nil {
			os.Remove(path)
			return nil, fmt.Errorf("chmod %s: %v", path, chErr)
		}

		own = true
	} else if (info.Mode() & os.ModeNamedPipe) == 0 {
		return nil, fmt.Errorf("%s is not a FIFO", path)
	}

	f, err := os.OpenFile(path, os.O_RDWR, os.ModeNamedPipe)
	if err != nil {
		if own {
			cleanup(path)
		}

		return nil, fmt.Errorf("open %s: %v", path, err)
	}

	return &source{
		ReadCloser: f,
		path:       path,
		kind:       fifo,
		own:        own,
	}, nil
}
