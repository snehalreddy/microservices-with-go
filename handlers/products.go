// Package classification of Product API
//
// Documentation for Product API
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
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

// A list of products are returned as response
// swagger:response productsResponse
type productsResponseWrapper struct {
	// All products in the system
	// in:body
	Body []data.Product
}

// swagger:response noContent
type productsNoContent struct {
}

// swagger:parameters deleteProducts
type productIDParameterWrapper struct {
	// The id of the product to delete from the database
	// in:path
	// required:true
	ID int `json:"id"`
}

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// swagger:route GET /products products listProducts
// Returns a list of products
// Responses:
// 	200: productsResponse

// GetProducts returns the products from the data store
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

// swagger:route DELETE /product/{id} products deleteProducts
// Returns a list of products
// Responses:
// 	201: noContent

// DeleteProducts deletes a product from the data store
func (p *Products) DeleteProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, strErr := strconv.Atoi(vars["id"])
	if strErr != nil {
		http.Error(rw, "Oops, cannot parse id...", http.StatusBadRequest)
		return
	}

	p.l.Println("Handle DELETE product with id:", id)

	err := data.DeleteProduct(id)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Product with id:%d, doesn't exist.", id), http.StatusBadRequest)
		return
	}
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
