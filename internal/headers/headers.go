package headers

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/absurek/go-http-from-tcp/internal/constants"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	idx := bytes.Index(data, []byte(constants.CRLF))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return 2, true, nil
	}

	headerText := string(data[:idx])
	key, value, err := headerFromString(headerText)
	if err != nil {
		return 0, false, err
	}

	h.Add(key, value)
	return idx + 2, false, nil
}

func (h Headers) Get(key string) string {
	return h[strings.ToLower(key)]
}

func (h Headers) Set(key, value string) {
	h[key] = value
}

func (h Headers) Add(key, value string) {
	key = strings.ToLower(key)

	current, ok := h[key]
	if ok {
		value = fmt.Sprintf("%s, %s", current, value)
	}

	h[key] = value
}

func headerFromString(str string) (string, string, error) {
	trimmed := strings.Trim(str, " ")

	idx := strings.Index(trimmed, ":")
	if idx == -1 || trimmed[idx-1] == ' ' {
		fmt.Println("returning error")
		return "", "", errors.New("malformed header")
	}

	key := strings.ToLower(trimmed[:idx])
	if !validTokens(key) {
		return "", "", errors.New("invalid token in header key")
	}

	value := strings.TrimLeft(trimmed[idx+1:], " ")

	return key, value, nil
}

var tokenChars = []rune{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}

func validTokens(data string) bool {
	for _, c := range data {
		if !isTokenChar(c) {
			return false
		}
	}

	return true
}

func isTokenChar(c rune) bool {
	if c >= 'A' && c <= 'Z' ||
		c >= 'a' && c <= 'z' ||
		c >= '0' && c <= '9' {

		return true
	}

	return slices.Contains(tokenChars, c)
}
