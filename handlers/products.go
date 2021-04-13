package handlers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

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

	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	if r.Method == http.MethodPut {
		// expect ID in the URI
		reg := regexp.MustCompile(`/([0-9]+)`)
		params := reg.FindAllStringSubmatch(r.URL.Path, -1)
		p.l.Println(params)

		if len(params) != 1 {
			p.l.Println("More than one id")
			http.Error(rw, "Invalid parameters", http.StatusBadRequest)
			return
		}
		if len(params[0]) != 2 {
			p.l.Println("More than control group match")
			http.Error(rw, "Invalid parameters", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(params[0][1])
		if err != nil {
			p.l.Println("Cannot parse number")
			http.Error(rw, "Invalid parameters", http.StatusBadRequest)
			return
		}
		p.l.Println("Got id:", id)

		p.updateProducts(id, rw, r)
		return
	}

	// catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products...")
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

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product...")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Oops, some error occured while parsing the data.", http.StatusBadRequest)
		return
	}

	p.l.Printf("Prod: %#v", prod)
	data.AddProduct(prod)
}

func (p *Products) updateProducts(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Product...")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Oops, some error occured while parsing the data.", http.StatusBadRequest)
		return
	}

	p.l.Printf("Prod: %#v", prod)
	err = data.PutProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(rw, "Patch successful...!")
}
