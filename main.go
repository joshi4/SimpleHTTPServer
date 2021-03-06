package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func init() {
	log.SetPrefix("[Info]")
}

// From Jackie Li's answer at https://groups.google.com/forum/#!topic/golang-nuts/_44ehpmFOjU
type WrapHTTPHandler struct {
	m http.Handler
}

func (h *WrapHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lw := &loggedResponse{ResponseWriter: w, status: 200}
	h.m.ServeHTTP(lw, r)
	log.Printf("%s %s - %d\n", r.Method, r.URL, lw.status)
}

type loggedResponse struct {
	http.ResponseWriter
	status int
}

func (l *loggedResponse) WriteHeader(status int) {
	l.status = status
	l.ResponseWriter.WriteHeader(status)
}

// From https://golang.org/doc/articles/wiki/
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		pwd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		fn(w, req, pwd)
	}
}

func okHandler(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(200)
}

func simpleHandler(w http.ResponseWriter, req *http.Request, pwd string) {
	http.ServeFile(w, req, fmt.Sprintf("%s%s", pwd, req.URL.Path))
}

func main() {
	port := ":9090"
	if len(os.Args) >= 2 {
		port = fmt.Sprintf(":%s", os.Args[1])
	}

	http.HandleFunc("/", makeHandler(simpleHandler))
	http.HandleFunc("/healthz", okHandler)

	log.Printf("Listening on port %s\n", port)
	err := http.ListenAndServe(port, &WrapHTTPHandler{http.DefaultServeMux})
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
