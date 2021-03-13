package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Hello handler
type Hello struct {
	l *log.Logger
}

// NewHello function that takes logger and returns Hello reference
func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	log.Println("Hello, world!")
	d, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(rw, "Oops", http.StatusBadRequest)
		//rw.WriteHeader(http.StatusBadRequest)
		//rw.Write([]byte("Oops"))
		return
	}

	//log.Printf("Data: %s\n", d)
	fmt.Fprintf(rw, "Data: %s\n", d)
}
