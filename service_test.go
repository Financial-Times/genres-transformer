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
		terms   []term
		genres  []genreLink
		found   bool
		err     error
	}{
		{"Success", "localhost:8080/transformers/genres/",
			[]term{term{CanonicalName: "Z_Archive", RawID: "b8337559-ac08-3404-9025-bad51ebe2fc7"}, term{CanonicalName: "Feature", RawID: "mNGQ2MWQ0NDMtMDc5Mi00NWExLTlkMGQtNWZhZjk0NGExOWU2-Z2VucVz"}},
			[]genreLink{genreLink{APIURL: "localhost:8080/transformers/genres/60351af9-c53d-3007-a891-978877d9c75e"},
				genreLink{APIURL: "localhost:8080/transformers/genres/9bb6ab36-6102-3bf0-896f-2600a63f51c8"}}, true, nil},
		{"Error on init", "localhost:8080/transformers/genres/", []term{}, []genreLink(nil), false, errors.New("Error getting taxonomy")},
	}

	for _, test := range tests {
		repo := dummyRepo{terms: test.terms, err: test.err}
		service, err := newGenreService(&repo, test.baseURL, "Genres", 10000)
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
		terms []term
		uuid  string
		genre genre
		found bool
		err   error
	}{
		{"Success", []term{term{CanonicalName: "Z_Archive", RawID: "b8337559-ac08-3404-9025-bad51ebe2fc7"}, term{CanonicalName: "Feature", RawID: "NGQ2MWQ0NDMtMDc5Mi00NWExLTlkMGQtNWZhZjk0NGExOWU2-Z2VucmVz"}},
			"60351af9-c53d-3007-a891-978877d9c75e", getDummyGenre("60351af9-c53d-3007-a891-978877d9c75e", "Z_Archive", "YjgzMzc1NTktYWMwOC0zNDA0LTkwMjUtYmFkNTFlYmUyZmM3-R2VucmVz"), true, nil},
		{"Not found", []term{term{CanonicalName: "Z_Archive", RawID: "845dc7d7-ae89-4fed-a819-9edcbb3fe507"}, term{CanonicalName: "Feature", RawID: "NGQ2MWdefsdfsfcmVz"}},
			"some uuid", genre{}, false, nil},
		{"Error on init", []term{}, "some uuid", genre{}, false, errors.New("Error getting taxonomy")},
	}

	for _, test := range tests {
		repo := dummyRepo{terms: test.terms, err: test.err}
		service, err := newGenreService(&repo, "", "Genres", 10000)
		expectedGenre, found := service.getGenreByUUID(test.uuid)
		assert.Equal(test.genre, expectedGenre, fmt.Sprintf("%s: Expected genre incorrect", test.name))
		assert.Equal(test.found, found)
		assert.Equal(test.err, err)
	}
}

type dummyRepo struct {
	terms []term
	err   error
}

func (d *dummyRepo) GetTmeTermsFromIndex(startRecord int) ([]interface{}, error) {
	if startRecord > 0 {
		return nil, d.err
	}
	var interfaces = make([]interface{}, len(d.terms))
	for i, data := range d.terms {
		interfaces[i] = data
	}
	return interfaces, d.err
}
func (d *dummyRepo) GetTmeTermById(uuid string) (interface{}, error) {
	return d.terms[0], d.err
}

func getDummyGenre(uuid string, prefLabel string, tmeId string) genre {
	return genre{
		UUID:                   uuid,
		PrefLabel:              prefLabel,
		PrimaryType:            primaryType,
		TypeHierarchy:          genreTypes,
		AlternativeIdentifiers: alternativeIdentifiers{TME: []string{tmeId}, Uuids: []string{uuid}}}
}
