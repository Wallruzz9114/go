package server

import "net/http"

// HandleIndex ...
func (server *Server) HandleIndex(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Length", "12")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	w.Write([]byte("Hello World!"))
}
