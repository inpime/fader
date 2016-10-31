package interfaces

type FileUserType string

var (
	ImageFile FileUserType = "image"
	TextFile  FileUserType = "text"
	RawFile   FileUserType = "raw"
)

func getUserTypeFromContentType(t string) FileUserType {
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
