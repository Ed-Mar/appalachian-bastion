package model

import (
	"backend/internal"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerMissingNameReturnsErr(t *testing.T) {
	server := Server{
		Description: "W H A T S H A P P I N G vol.8",
	}

	validator := internal.NewValidation()
	err := validator.Validate(server)
	assert.Len(t, err, 1)
}

//TODO this error is returning nil need to handle it
func TestValidServerDoesNOTReturnsErr(t *testing.T) {
	server := Server{
		ID:          413,
		Name:        "Y33T in the Shade",
		Description: "PHONK 808",
	}

	validator := internal.NewValidation()
	err := validator.Validate(server)
	assert.Len(t, err, 1)
}

func TestServersToJSON(t *testing.T) {
	ps := []*Server{
		&Server{
			Name: "S2",
		},
	}

	b := bytes.NewBufferString("")
	err := internal.ToJSON(ps, b)
	assert.NoError(t, err)
}
