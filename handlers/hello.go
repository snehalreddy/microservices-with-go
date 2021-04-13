package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Hello, world!")

	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// rw.WriteHeader(http.StatusBadRequest)
		// rw.Write([]byte("Oops, bad request"))

		// convenience method from http can be used instead of the above two lines
		http.Error(rw, "Oops, bad request", http.StatusBadRequest)
		return
	}

	// log.Printf("Data: %s\n", d)
	fmt.Fprintf(rw, "Hello, %s.\n", d)
}
