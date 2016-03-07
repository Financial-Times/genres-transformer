package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

const testUUID = "bba39990-c78d-3629-ae83-808c333c6dbc"
const getGenresResponse = "[{\"apiUrl\":\"http://localhost:8080/transformers/genres/bba39990-c78d-3629-ae83-808c333c6dbc\"}]\n"
const getGenresByUUIDResponse = "{\"uuid\":\"bba39990-c78d-3629-ae83-808c333c6dbc\",\"canonicalName\":\"Obituary\",\"tmeIdentifier\":\"Ng==-R2VucmVz\",\"type\":\"Genre\"}\n"

func TestHandlers(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name         string
		req          *http.Request
		dummyService genreService
		statusCode   int
		contentType  string // Contents of the Content-Type header
		body         string
	}{
		{"Success - get genre by uuid", newRequest("GET", fmt.Sprintf("/transformers/genres/%s", testUUID)), &dummyService{found: true, genres: []genre{genre{UUID: testUUID, CanonicalName: "Obituary", TmeIdentifier: "Ng==-R2VucmVz", Type: "Genre"}}}, http.StatusOK, "application/json", getGenresByUUIDResponse},
		{"Not found - get genre by uuid", newRequest("GET", fmt.Sprintf("/transformers/genres/%s", testUUID)), &dummyService{found: false, genres: []genre{genre{}}}, http.StatusNotFound, "application/json", ""},
		{"Success - get genres", newRequest("GET", "/transformers/genres"), &dummyService{found: true, genres: []genre{genre{UUID: testUUID}}}, http.StatusOK, "application/json", getGenresResponse},
		{"Not found - get genres", newRequest("GET", "/transformers/genres"), &dummyService{found: false, genres: []genre{}}, http.StatusNotFound, "application/json", ""},
	}

	for _, test := range tests {
		rec := httptest.NewRecorder()
		router(test.dummyService).ServeHTTP(rec, test.req)
		assert.True(test.statusCode == rec.Code, fmt.Sprintf("%s: Wrong response code, was %d, should be %d", test.name, rec.Code, test.statusCode))
		assert.Equal(test.body, rec.Body.String(), fmt.Sprintf("%s: Wrong body", test.name))
	}
}

func newRequest(method, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	return req
}

func router(s genreService) *mux.Router {
	m := mux.NewRouter()
	h := newGenresHandler(s)
	m.HandleFunc("/transformers/genres", h.getGenres).Methods("GET")
	m.HandleFunc("/transformers/genres/{uuid}", h.getGenreByUUID).Methods("GET")
	return m
}

type dummyService struct {
	found  bool
	genres []genre
}

func (s *dummyService) getGenres() ([]genreLink, bool) {
	var genreLinks []genreLink
	for _, sub := range s.genres {
		genreLinks = append(genreLinks, genreLink{APIURL: "http://localhost:8080/transformers/genres/" + sub.UUID})
	}
	return genreLinks, s.found
}

func (s *dummyService) getGenreByUUID(uuid string) (genre, bool) {
	return s.genres[0], s.found
}
