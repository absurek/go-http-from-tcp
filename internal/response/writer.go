package response

import (
	"bytes"
	"fmt"

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

func (w *Writer) WriteTrailers(headers headers.Headers) error {
	return WriteHeaders(&w.buffer, headers)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	chunkSize := len(p)

	nTotal := 0
	n, err := fmt.Fprintf(&w.buffer, "%x\r\n", chunkSize)
	if err != nil {
		return nTotal, err
	}
	nTotal += n

	n, err = w.buffer.Write(p)
	if err != nil {
		return nTotal, err
	}
	nTotal += n

	n, err = w.buffer.Write([]byte("\r\n"))
	if err != nil {
		return nTotal, err
	}
	nTotal += n

	return nTotal, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	n, err := w.buffer.Write([]byte("0\r\n"))
	if err != nil {
		return n, err
	}

	return n, nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	return w.buffer.Write(p)
}
