package handlers

import (
	"log"
	"net/http"

	"github.com/snehalreddy/MicroGoIntro/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	// catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	// d, err := json.Marshal(lp)
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Oops, some error occured while parsing the data.", http.StatusInternalServerError)
		return
	}

	// fmt.Fprintf(rw, "product_list: %v", d)
	// rw.Write(d)
}
