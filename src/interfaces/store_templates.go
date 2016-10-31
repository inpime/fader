package interfaces

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"strings"
)

var _ TemplatesLoader = (*TemplatesStore)(nil)

const (
	DefaultPathSeparator = "/"
)

func NewTemplatesStore(manager FileManager) *TemplatesStore {
	return &TemplatesStore{
		manager: manager,
		logger:  log.New(os.Stderr, "[TEMPLATE_LOADER]", -1),
	}
}

type TemplatesStore struct {
	manager FileManager
	logger  *log.Logger
}

func (l *TemplatesStore) Abs(base, name string) string {
	return name
}

func (l *TemplatesStore) Get(path string) (io.Reader, error) {
	var file *File

	path = strings.TrimSpace(path)
	if len(path) <= 2 {
		l.logger.Println("[ERR] intvalid path, length less than 2,", path)
		return bytes.NewReader([]byte{}), errors.New("too short file path")
	}

	_path := strings.Split(path, DefaultPathSeparator)
	if len(_path) < 2 {
		l.logger.Println("[ERR] intvalid path, file name is empty,", _path)
		return bytes.NewReader([]byte{}), errors.New("file name is empty")
	}

	bucketName := _path[0]
	fileName := strings.Join(_path[1:], DefaultPathSeparator)

	file, err := l.manager.FindFileByName(bucketName, fileName,
		RawData, // TODO: access control
	)

	if err != nil {
		l.logger.Println("[ERR] find file by names,", bucketName, fileName)
		return nil, err
	}

	return bytes.NewReader(file.RawData), nil
}

// api.tempaltes#TemplatesLoader

type TemplatesLoader interface {
	Abs(base, name string) string
	Get(path string) (io.Reader, error)
}
