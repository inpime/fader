package interfaces

import "bufio"
import "bytes"

var (
	ImageFile string = "image"
	TextFile  string = "text"
	RawFile   string = "raw"
)

func GetUserTypeFromContentType(t string) string {
	switch t {
	case "image/jpeg",
		"image/pjpeg",
		"image/png",
		"image/vnd.microsoft.icon",
		"image/gif":
		return ImageFile

	case "text/css",
		"text/plain",
		"text/javascript",
		"text/html",
		"text/toml",
		"application/javascript",
		"application/json",
		"application/soap+xml",
		"application/xhtml+xml",
		"text/csv",
		"text/x-jquery-tmpl",
		"text/php",
		"application/x-javascript":
		return TextFile

	default:
		return RawFile
	}
}

// untils for mporter

func readLine(scanner *bufio.Reader) ([]byte, error) {
	line, hasMore, err := scanner.ReadLine()

	if err != nil {
		return []byte{}, err
	}

	if !hasMore {
		return line, nil
	}

	buf := bytes.NewBuffer([]byte{})
	buf.Write(line)

	for true {
		line, hasMore, err = scanner.ReadLine()

		buf.Write(line)

		if !hasMore {
			break
		}
	}

	return buf.Bytes(), err
}
