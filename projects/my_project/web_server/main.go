package main

import (
	"my_project/logger" // import logger
	"net/http"
)

func main() {
	l := new(logger.Logger) // create and use a new logger
	l.LogError("Not Found")
	http.Handle("/hello", http.HandlerFunc(handle))
	http.ListenAndServe(":5500", nil)
}

func handle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("world"))
}
