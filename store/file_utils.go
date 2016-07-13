package store

func getTypeNameFromContentType(t string) string {
	switch t {
	case "image/jpeg",
		"image/pjpeg",
		"image/png",
		"image/vnd.microsoft.icon",
		"image/gif":
		return "image"

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
		return "text"

	default:
		return "raw"
	}
}
