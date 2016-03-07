package main

import (
	"net/http"
)

type httpClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

type genreService interface {
	getGenres() ([]genreLink, bool)
	getGenreByUUID(uuid string) (genre, bool)
}

type genreServiceImpl struct {
	repository repository
	baseURL    string
	genresMap  map[string]genre
	genreLinks []genreLink
}

func newGenreService(repo repository, baseURL string) (genreService, error) {

	s := &genreServiceImpl{repository: repo, baseURL: baseURL}
	err := s.init()
	if err != nil {
		return &genreServiceImpl{}, err
	}
	return s, nil
}

func (s *genreServiceImpl) init() error {
	s.genresMap = make(map[string]genre)
	tax, err := s.repository.getGenresTaxonomy()
	if err != nil {
		return err
	}
	s.initGenresMap(tax.Terms)
	return nil
}

func (s *genreServiceImpl) getGenres() ([]genreLink, bool) {
	if len(s.genreLinks) > 0 {
		return s.genreLinks, true
	}
	return s.genreLinks, false
}

func (s *genreServiceImpl) getGenreByUUID(uuid string) (genre, bool) {
	genre, found := s.genresMap[uuid]
	return genre, found
}

func (s *genreServiceImpl) initGenresMap(terms []term) {
	for _, t := range terms {
		sub := transformGenre(t)
		s.genresMap[sub.UUID] = sub
		s.genreLinks = append(s.genreLinks, genreLink{APIURL: s.baseURL + sub.UUID})
		s.initGenresMap(t.Children.Terms)
	}
}
