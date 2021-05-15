package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)
	return validate.Struct(p)
}

func validateSKU(fl validator.FieldLevel) bool {
	// sku is of the format sdfb-jashdf-sjhdf
	re := regexp.MustCompile(`^[a-z]+-[a-z]+-[a-z]+$`)
	matches := re.FindAllString(fl.Field().String(), -1)

	return len(matches) == 1
}

func AddProduct(p *Product) {
	p.ID = getNextId()
	productList = append(productList, p)
}

func getNextId() int {
	lId := productList[len(productList)-1].ID
	lId += 1
	return lId
}

type Products []*Product

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func GetProducts() Products {
	return productList
}

var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Strong and short coffee without milk",
		Price:       1.99,
		SKU:         "fbj324",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}

var ErrProductNotFound = fmt.Errorf("Product not found")

func PutProduct(id int, prod *Product) error {
	pos, err := findProduct(id)
	if err != nil {
		return ErrProductNotFound
	}

	prod.ID = id
	productList[pos] = prod
	return nil
}

func findProduct(id int) (int, error) {
	for i, v := range productList {
		if v.ID == id {
			return i, nil
		}
	}
	return -1, ErrProductNotFound
}
