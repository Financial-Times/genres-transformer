package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransform(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		term    term
		genre genre
	}{
		{"Trasform term to genre", term{CanonicalName: "Obituary", ID: "Ng==-R2VucmVz"}, genre{UUID: "2c4a7847-11ad-308b-b634-4c962708261c", CanonicalName: "Obituary", TmeIdentifier: "Ng==-R2VucmVz", Type: "Genre"}},
	}

	for _, test := range tests {
		expectedGenre := transformGenre(test.term)
		assert.Equal(test.genre, expectedGenre, fmt.Sprintf("%s: Expected genre incorrect", test.name))
	}

}
