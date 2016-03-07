package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGenres(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		baseURL string
		tax     taxonomy
		genres  []genreLink
		found   bool
		err     error
	}{
		{"Success", "localhost:8080/transformers/genres/",
			taxonomy{Terms: []term{term{CanonicalName: "Obituary", ID: "Ng==-R2VucmVz"}}},
			[]genreLink{genreLink{APIURL: "localhost:8080/transformers/genres/2c4a7847-11ad-308b-b634-4c962708261c"}}, true, nil},
		{"Error on init", "localhost:8080/transformers/genres/", taxonomy{}, []genreLink(nil), false, errors.New("Error getting taxonomy")},
	}

	for _, test := range tests {
		repo := dummyRepo{tax: test.tax, err: test.err}
		service, err := newGenreService(&repo, test.baseURL)
		expectedGenres, found := service.getGenres()
		assert.Equal(test.genres, expectedGenres, fmt.Sprintf("%s: Expected genres link incorrect", test.name))
		assert.Equal(test.found, found)
		assert.Equal(test.err, err)
	}
}

func TestGetGenreByUuid(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name  string
		tax   taxonomy
		uuid  string
		genre genre
		found bool
		err   error
	}{
		{"Success", taxonomy{Terms: []term{term{CanonicalName: "Obituary", ID: "Ng==-R2VucmVz"}}},
			"2c4a7847-11ad-308b-b634-4c962708261c", genre{UUID: "2c4a7847-11ad-308b-b634-4c962708261c", CanonicalName: "Obituary", TmeIdentifier: "Ng==-R2VucmVz", Type: "Genre"}, true, nil},
		{"Not found", taxonomy{Terms: []term{term{CanonicalName: "Obituary", ID: "Ng==-R2VucmVz"}}},
			"some uuid", genre{}, false, nil},
		{"Error on init", taxonomy{}, "some uuid", genre{}, false, errors.New("Error getting taxonomy")},
	}

	for _, test := range tests {
		repo := dummyRepo{tax: test.tax, err: test.err}
		service, err := newGenreService(&repo, "")
		expectedGenre, found := service.getGenreByUUID(test.uuid)
		assert.Equal(test.genre, expectedGenre, fmt.Sprintf("%s: Expected genre incorrect", test.name))
		assert.Equal(test.found, found)
		assert.Equal(test.err, err)
	}
}

type dummyRepo struct {
	tax taxonomy
	err error
}

func (d *dummyRepo) getGenresTaxonomy() (taxonomy, error) {
	return d.tax, d.err
}
