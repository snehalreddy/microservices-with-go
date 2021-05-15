package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name: "Awesome product",
		Price: 1,
		SKU: "skdfjs-sdfkj-aaksdj",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}