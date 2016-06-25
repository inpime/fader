package store

import (
	"github.com/Sirupsen/logrus"
	"utils"
)

func (f File) MMeta() utils.M {
	return utils.M(f.Meta())
}

func (f *File) SetContentType(v string) *File {
	utils.M(f.Meta()).Set(ContentTypeKey, v)
	return f
}

func (f File) ContentType() string {

	return utils.Map(f.Meta()).String(ContentTypeKey)
}

func (f *File) MMapData() utils.M {
	return utils.M(f.MapData())
}

func (f *File) SetMapData(v map[string]interface{}) *File {
	f.File.SetMapData(v)
	return f
}

func (f File) TextData() string {

	rawData := f.RawData().Bytes()

	if len(rawData) > 1024*1024*10 {
		logrus.WithField("length", len(rawData)).Warning("file raw data as text: to long")
		return ""
	}

	return string(rawData)
}

func (f *File) SetTextData(src string) *File {
	f.RawData().Write([]byte(src))
	return f
}

func (f *File) SetRawData(src []byte) *File {
	f.RawData().Write(src)
	return f
}

func (f File) IsImage() bool {
	return getTypeNameFromContentType(f.ContentType()) == "image"
}

func (f File) IsText() bool {
	return getTypeNameFromContentType(f.ContentType()) == "text"
}

func (f File) IsRaw() bool {
	return getTypeNameFromContentType(f.ContentType()) == "raw"
}

// for pongo2 (must have exactly 1 output argument)

func (f File) SetName(name string) File {
	f.File.SetName(name)
	return f
}
