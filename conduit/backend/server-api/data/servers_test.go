package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Server{
		ID:   69,
		Name: "Y33t From the Mts",
	}
	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
