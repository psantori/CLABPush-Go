package server

import (
	"log"
	"net/http"
)

// Handler can be "subclassed" by other handlers.
type Handler struct {
	Ctx *AppContext
}

// RenderHTTPStatus is an utility function that will write the status code in
// the header and log both the status code and the optional error.
func (h *Handler) RenderHTTPStatus(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	log.Printf("%d - %s", status, http.StatusText(status))
	if err != nil {
		log.Println(err)
	}
}
