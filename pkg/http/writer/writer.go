package writer

import (
	"encoding/json"
	"net/http"

	"backend_gen/pkg/http/header"
)

func WriteJson(w http.ResponseWriter, data any) {
	header.AddJSONContentType(w.Header())

	encoder := json.NewEncoder(w)

	if err := encoder.Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func WriteStatusOK(w http.ResponseWriter) {
	header.AddJSONContentType(w.Header())
	w.WriteHeader(http.StatusOK)
}
