package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/snehalreddy/MicroGoIntro/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
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

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product...")

	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	data.AddProduct(prod)
}

func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, strErr := strconv.Atoi(vars["id"])
	if strErr != nil {
		http.Error(rw, "Oops, cannot parse id...", http.StatusBadRequest)
		return
	}

	p.l.Println("Handle PUT product with id:", id)

	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	err := data.PutProduct(id, prod)
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

type KeyProduct struct{}

func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}
		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Oops, some error occured while parsing the data.", http.StatusBadRequest)
			return
		}

		// validate the product
		valErr := prod.Validate()
		if valErr != nil {
			p.l.Println("[ERROR] validating product", valErr)
			http.Error(
				rw,
				fmt.Sprintf("Oops, some error occured while parsing the data, %s", valErr),
				http.StatusBadRequest,
			)
			return
		}

		p.l.Printf("Prod: %#v", prod)
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
