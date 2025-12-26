package response

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/absurek/go-http-from-tcp/internal/constants"
	"github.com/absurek/go-http-from-tcp/internal/headers"
)

type StatusCode int

const (
	StatusOk                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

var statusPhraseMap = map[StatusCode]string{
	StatusOk:                  "OK",
	StatusBadRequest:          "Bad Request",
	StatusInternalServerError: "Internal Server Error",
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	phrase, ok := statusPhraseMap[statusCode]
	if !ok {
		return errors.New("unkown status code")
	}

	_, err := fmt.Fprintf(w, "HTTP/1.1 %d %s%s", statusCode, phrase, constants.CRLF)
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()

	headers.Set("Content-Type", "text/plain")
	headers.Set("Content-Length", strconv.Itoa(contentLen))
	headers.Set("Connection", "close")

	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		_, err := fmt.Fprintf(w, "%s: %s%s", key, value, constants.CRLF)
		if err != nil {
			return err
		}
	}

	_, err := w.Write([]byte(constants.CRLF))
	return err
}
