package main

import (
	"github.com/gorilla/handlers"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type NotFoundRedirectRespWr struct {
	http.ResponseWriter // We embed http.ResponseWriter
	status              int
}

func (w *NotFoundRedirectRespWr) WriteHeader(status int) {
	w.status = status // Store the status for our own use
	if status != http.StatusNotFound {
		w.ResponseWriter.WriteHeader(status)
	}
}

func (w *NotFoundRedirectRespWr) Write(p []byte) (int, error) {
	if w.status != http.StatusNotFound {
		return w.ResponseWriter.Write(p)
	}
	return len(p), nil // Lie that we successfully written it
}

func wrapHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nfrw := &NotFoundRedirectRespWr{ResponseWriter: w}
		h.ServeHTTP(nfrw, r)
		if nfrw.status == http.StatusNotFound {
			newR := r
			newR.URL.Path = "/"
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			s, err := ioutil.ReadFile("files/index.html")
			if err != nil {
				_, err := w.Write([]byte("error"))
				if err != nil {
					log.Fatal("ListenAndServe: ", err)
				}
			} else {
				_, err := w.Write(s)
				if err != nil {
					log.Fatal("ListenAndServe: ", err)
				}
			}
		}
	}
}

func main() {
	port := "2333"
	http.Handle("/", http.StripPrefix("/", wrapHandler(http.FileServer(http.Dir("files")))))
	err := http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, http.DefaultServeMux))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
