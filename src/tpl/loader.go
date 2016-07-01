package tpl

import (
	"bytes"
	"fmt"
	"io"
	"store"
	"strings"
)

type Pongo2BoltdDBLoader struct {
	// TODO: logger
}

func MustNewBoltdDBLoader() *Pongo2BoltdDBLoader {
	loader, err := NewBoltdDBLoader()
	if err != nil {
		panic(err)
	}
	return loader
}

func NewBoltdDBLoader() (*Pongo2BoltdDBLoader, error) {
	loader := &Pongo2BoltdDBLoader{}

	return loader, nil
}

func (l *Pongo2BoltdDBLoader) Get(path string) (io.Reader, error) {
	_sep := "/"
	_path := strings.Split(path, _sep)

	if path == "" {
		return bytes.NewReader([]byte{}), nil
	}

	if len(_path) < 2 {

		return nil, fmt.Errorf("pongo2loader: invalid path %q", path)
	}

	bucketName := _path[0]
	fileName := strings.Join(_path[1:], _sep)

	file, err := store.LoadOrNewFile(bucketName, fileName)

	if err != nil {

		return nil, err
	}

	return bytes.NewReader(file.RawData().Bytes()), nil
}

// implement pongo2.TemplateLoader

func (l *Pongo2BoltdDBLoader) Abs(base, name string) string {
	return name
}
