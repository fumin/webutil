package webutil

import (
	"net/http"
	"time"
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

// A Sse is a wrapper over a Server-Sent Events request.
//
// Aside from providing helpers to send events in the SSE format, it also starts
// a background goroutine that sends heartbeats to the client on the other side,
// in an attempt to detect connection failures.
// Users are notified of such failures by listening on the ConnClosed channel.
//
// Note that users are expected to call Close in order to cleanup the background
// goroutine and other resources including time.Tickers when they do not need
// the Sse anymore.
// Moreover, not calling Close might result in crashes since the background
// goroutine calls Sse.Write when sending heartbeats, which in turn calls Flush.
// This causes the http package to *panic* since in this case,
// http.Flusher.Flush is called after the http handlers return.
// https://groups.google.com/d/msg/Golang-Nuts/qcjLQ4O8Pc4/BrgYJF4mENMJ
type Sse struct {
	w          http.ResponseWriter
	stopTicker chan bool
	ConnClosed chan bool
}

func NewServerSideEventWriter(w http.ResponseWriter, heartbeat string, d time.Duration) Sse {
	headers := w.Header()
	headers.Set("Content-Type", "text/event-stream")
	headers.Set("Cache-Control", "no-cache")
	headers.Set("Connection", "keep-alive")

	ticker := time.NewTicker(d)
  sse := Sse{w: w, stopTicker: make(chan bool, 1), ConnClosed: make(chan bool, 1)}
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				err := sse.EventWrite(heartbeat, make([]byte, 0))
				if err != nil {
          select {
          case sse.ConnClosed <- true:
          default:
          }
					return
				}
			case <-sse.stopTicker:
				return
			}
		}
	}()
	return sse
}

func (sse Sse) Close() {
  select {
  case sse.stopTicker <- true:
  default:
  }
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
