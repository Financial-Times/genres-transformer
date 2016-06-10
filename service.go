package main

import (
	"github.com/Financial-Times/tme-reader/tmereader"
	log "github.com/Sirupsen/logrus"
	"net/http"
)

type httpClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

type genreService interface {
	getGenres() ([]genreLink, bool)
	getGenreByUUID(uuid string) (genre, bool)
	checkConnectivity() error
}

type genreServiceImpl struct {
	repository    tmereader.Repository
	baseURL       string
	genresMap     map[string]genre
	genreLinks    []genreLink
	taxonomyName  string
	maxTmeRecords int
}

func newGenreService(repo tmereader.Repository, baseURL string, taxonomyName string, maxTmeRecords int) (genreService, error) {
	s := &genreServiceImpl{repository: repo, baseURL: baseURL, taxonomyName: taxonomyName, maxTmeRecords: maxTmeRecords}
	err := s.init()
	if err != nil {
		return &genreServiceImpl{}, err
	}
	return s, nil
}

func (s *genreServiceImpl) init() error {
	s.genresMap = make(map[string]genre)
	responseCount := 0
	log.Printf("Fetching genres from TME\n")
	for {
		terms, err := s.repository.GetTmeTermsFromIndex(responseCount)
		if err != nil {
			return err
		}

		if len(terms) < 1 {
			log.Printf("Finished fetching genres from TME\n")
			break
		}
		s.initGenresMap(terms)
		responseCount += s.maxTmeRecords
	}
	log.Printf("Added %d genre links\n", len(s.genreLinks))

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

func (s *genreServiceImpl) checkConnectivity() error {
	// TODO: Can we just hit an endpoint to check if TME is available? Or do we need to make sure we get genre taxonmies back?
	//	_, err := s.repository.GetTmeTermsFromIndex()
	//	if err != nil {
	//		return err
	//	}
	return nil
}

func (s *genreServiceImpl) initGenresMap(terms []interface{}) {
	for _, iTerm := range terms {
		t := iTerm.(term)
		top := transformGenre(t, s.taxonomyName)
		s.genresMap[top.UUID] = top
		s.genreLinks = append(s.genreLinks, genreLink{APIURL: s.baseURL + top.UUID})
	}
}
