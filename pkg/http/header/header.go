package header

import "net/http"

const (
	ContentType     = "Content-Type"
	JSONContentType = "application/json"
)

func AddJSONContentType(header http.Header) {
	header.Set(ContentType, JSONContentType)
}
