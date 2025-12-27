package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/absurek/go-http-from-tcp/internal/constants"
	"github.com/absurek/go-http-from-tcp/internal/headers"
	"github.com/absurek/go-http-from-tcp/internal/request"
	"github.com/absurek/go-http-from-tcp/internal/response"
	"github.com/absurek/go-http-from-tcp/internal/server"
)

const httpBinBaseURL string = "https://httpbin.org"

func httpBin(w *response.Writer, req *request.Request) {
	url := httpBinBaseURL + strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin")

	resp, err := http.Get(url)
	if err != nil {
		internalServerError(w, req)
		return
	}
	defer resp.Body.Close()

	w.WriteStatusLine(response.StatusOk)

	h := response.GetDefaultHeaders(0)
	h.Set("Transfer-Encoding", "chunked")
	h.Set("Trailer", "X-Content-Sha256")
	h.Add("Trailer", "X-Content-Length")
	h.Remove("Content-Length")
	w.WriteHeaders(h)

	const maxChunkSize = 1024
	buffer := make([]byte, maxChunkSize)
	var fullBody bytes.Buffer
	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			_, err = w.WriteChunkedBody(buffer[:n])
			if err != nil {
				log.Println("Error: writing chunked body:", err)
				break
			}

			fullBody.Write(buffer[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Error: reading response body:", err)
			break
		}
	}
	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		log.Println("Error: writing chunked body done:", err)
	}

	trailers := headers.NewHeaders()
	sum := sha256.Sum256(fullBody.Bytes())
	trailers.Set("X-Content-Sha256", hex.EncodeToString(sum[:]))
	trailers.Set("X-Content-Length", strconv.Itoa(len(fullBody.Bytes())))
	w.WriteTrailers(trailers)
}

func handler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		httpBin(w, req)
		return
	}

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		badRequest(w, req)
	case "/myproblem":
		internalServerError(w, req)
	case "/video":
		video(w, req)
	default:
		statusOk(w, req)
	}
}

func main() {
	server, err := server.Serve(constants.Port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", constants.Port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func badRequest(w *response.Writer, req *request.Request) {
	res := []byte(`
			<html>
				<head>
					<title>400 Bad Request</title>
				</head>
				<body>
					<h1>Bad Request</h1>
					<p>Your request honestly kinda sucked.</p>
				</body>
			</html>`)

	w.WriteStatusLine(response.StatusBadRequest)
	headers := response.GetDefaultHeaders(len(res))
	headers.Set("Content-Type", "text/html")

	w.WriteHeaders(headers)
	w.WriteBody(res)
}

func internalServerError(w *response.Writer, req *request.Request) {
	res := []byte(`
			<html>
				<head>
					<title>500 Internal Server Error</title>
				</head>
				<body>
					<h1>Internal Server Error</h1>
					<p>Okay, you know what? This one is on me.</p>
				</body>
			</html>`)

	w.WriteStatusLine(response.StatusInternalServerError)
	headers := response.GetDefaultHeaders(len(res))
	headers.Set("Content-Type", "text/html")

	w.WriteHeaders(headers)
	w.WriteBody(res)
}

func statusOk(w *response.Writer, req *request.Request) {
	res := []byte(`
			<html>
				<head>
					<title>200 OK</title>
				</head>
				<body>
					<h1>Success!</h1>
					<p>Your request was an absolute banger.</p>
				</body>
			</html>`)

	w.WriteStatusLine(response.StatusOk)
	headers := response.GetDefaultHeaders(len(res))
	headers.Set("Content-Type", "text/html")

	w.WriteHeaders(headers)
	w.WriteBody(res)
}

func video(w *response.Writer, req *request.Request) {
	file, err := os.ReadFile("./assets/vim.mp4")
	if err != nil {
		internalServerError(w, req)
		return
	}

	w.WriteStatusLine(response.StatusOk)
	headers := response.GetDefaultHeaders(len(file))
	headers.Set("Content-Type", "video/mp4")

	w.WriteHeaders(headers)
	w.WriteBody(file)
}
