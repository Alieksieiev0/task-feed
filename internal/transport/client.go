package transport

import (
	"bufio"
	"context"
	"fmt"
)

func NewHttpClient(w *bufio.Writer, cancel context.CancelFunc) *HttpClient {
	return &HttpClient{w: w, cancel: cancel}
}

type HttpClient struct {
	w      *bufio.Writer
	cancel context.CancelFunc
}

func (h *HttpClient) Write(data string) error {
	fmt.Fprint(h.w, data, "\n")
	return h.w.Flush()
}

func (h *HttpClient) Cancel() {
	h.cancel()
}
