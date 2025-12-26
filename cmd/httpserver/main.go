package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/absurek/go-http-from-tcp/internal/constants"
	"github.com/absurek/go-http-from-tcp/internal/request"
	"github.com/absurek/go-http-from-tcp/internal/response"
	"github.com/absurek/go-http-from-tcp/internal/server"
)

func handler(w *response.Writer, req *request.Request) {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
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
	case "/myproblem":
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
	default:
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
