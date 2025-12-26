package response

import (
	"bytes"

	"github.com/absurek/go-http-from-tcp/internal/headers"
)

type Writer struct {
	buffer bytes.Buffer
}

func (w *Writer) Bytes() []byte {
	return w.buffer.Bytes()
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	return WriteStatusLine(&w.buffer, statusCode)
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	return WriteHeaders(&w.buffer, headers)
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	return w.buffer.Write(p)
}
