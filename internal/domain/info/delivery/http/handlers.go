package http

import "net/http"

type InfoHandler struct {
}

func NewInfoHandler() *InfoHandler {
	return &InfoHandler{}
}

func (i *InfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
