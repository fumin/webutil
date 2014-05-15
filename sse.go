package webutil

import (
	"net/http"
)

// ByteWriter is a one-off utility type for writing to http.ResponseWriter
type ByteWriter struct {
	RespWriter http.ResponseWriter
	Err        error
}

func (w *ByteWriter) Write(b []byte) {
	if w.Err != nil {
		return
	}
	_, w.Err = w.RespWriter.Write(b)
}

// A Sse is a wrapper over a Server-Sent Events response.
type Sse struct {
	w          http.ResponseWriter
}

func NewServerSideEventsWriter(w http.ResponseWriter) Sse {
	headers := w.Header()
	headers.Set("Content-Type", "text/event-stream")
	headers.Set("Cache-Control", "no-cache")
	headers.Set("Connection", "keep-alive")
  return Sse{w: w}
}

func (sse Sse) Write(b []byte) error {
	bw := &ByteWriter{RespWriter: sse.w}
	bw.Write([]byte("data: "))
	bw.Write(b)
	bw.Write([]byte("\n\n"))
	if bw.Err != nil {
		return bw.Err
	}
	if f, ok := sse.w.(http.Flusher); ok {
		f.Flush()
	}
	return nil
}

func (sse Sse) EventWrite(event string, b []byte) error {
	bw := &ByteWriter{RespWriter: sse.w}
	bw.Write([]byte("event: "))
	bw.Write([]byte(event))
	bw.Write([]byte("\n"))
	if bw.Err != nil {
		return bw.Err
	}
	return sse.Write(b)
}
