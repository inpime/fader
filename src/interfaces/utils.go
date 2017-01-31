package interfaces

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
