package context

import (
	"bytes"
	"io"
)

// FileData the data file
type FileData struct {
	Data        []byte
	Name        string
	ContentType string
}

func (f FileData) IsEmpty() bool {
	return len(f.Data) == 0
}

// FormFileData file data from the form
func (c Context) FormFileData(name string) *FileData {
	file, err := c.FormFile(name)
	if err != nil {
		return &FileData{}
	}

	// fileName := file.Filename
	fileOpen, _ := file.Open()
	defer fileOpen.Close()

	fileContent := bytes.NewBuffer([]byte{})
	_, err = io.Copy(fileContent, fileOpen)

	if err != nil {
		return &FileData{}
	}

	return &FileData{
		Data:        fileContent.Bytes(),
		Name:        file.Filename,
		ContentType: file.Header.Get("Content-Type"),
	}
}
