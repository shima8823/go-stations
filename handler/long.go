package handler

import (
	"net/http"
	"time"
)

type longHandler struct{}

func NewLong() *longHandler {
	return &longHandler{}
}

func (h *longHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	time.Sleep(10 * time.Second)
	w.Write([]byte("Finished long request!"))
}
