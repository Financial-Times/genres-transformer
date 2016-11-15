package main

import (
	"encoding/json"
	"fmt"
	"github.com/Financial-Times/go-fthealth/v1a"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type genresHandler struct {
	service genreService
}

// HealthCheck does something
func (h *genresHandler) HealthCheck() v1a.Check {
	return v1a.Check{
		BusinessImpact:   "Unable to respond to request for the genre data from TME",
		Name:             "Check connectivity to TME",
		PanicGuide:       "https://sites.google.com/a/ft.com/ft-technology-service-transition/home/run-book-library/genres-transfomer",
		Severity:         1,
		TechnicalSummary: "Cannot connect to TME to be able to supply genres",
		Checker:          h.checker,
	}
}

// Checker does more stuff
func (h *genresHandler) checker() (string, error) {
	err := h.service.checkConnectivity()
	if err == nil {
		return "Connectivity to TME is ok", err
	}
	return "Error connecting to TME", err
}

func newGenresHandler(service genreService) genresHandler {
	return genresHandler{service: service}
}

func (h *genresHandler) getGenres(writer http.ResponseWriter, req *http.Request) {
	obj, found := h.service.getGenres()
	writeJSONResponse(obj, found, writer)
}

func (h *genresHandler) getGenreByUUID(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	uuid := vars["uuid"]

	obj, found := h.service.getGenreByUUID(uuid)
	writeJSONResponse(obj, found, writer)
}

func writeJSONResponse(obj interface{}, found bool, writer http.ResponseWriter) {
	writer.Header().Add("Content-Type", "application/json")

	if !found {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	enc := json.NewEncoder(writer)
	if err := enc.Encode(obj); err != nil {
		log.Errorf("Error on json encoding=%v\n", err)
		writeJSONError(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func writeJSONError(w http.ResponseWriter, errorMsg string, statusCode int) {
	w.WriteHeader(statusCode)
	fmt.Fprintln(w, fmt.Sprintf("{\"message\": \"%s\"}", errorMsg))
}

//GoodToGo returns a 503 if the healthcheck fails - suitable for use from varnish to check availability of a node
func (h *genresHandler) GoodToGo(writer http.ResponseWriter, req *http.Request) {
	if _, err := h.checker(); err != nil {
		writer.WriteHeader(http.StatusServiceUnavailable)
	}
}

func (h *genresHandler) getCount(writer http.ResponseWriter, req *http.Request) {
	count := h.service.getGenreCount()
	_, err := writer.Write([]byte(strconv.Itoa(count)))
	if err != nil {
		log.Warnf("Couldn't write count to HTTP response. count=%d %v\n", count, err)
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *genresHandler) getIds(writer http.ResponseWriter, req *http.Request) {
	ids := h.service.getGenreIds()
	writer.Header().Add("Content-Type", "text/plain")
	if len(ids) == 0 {
		writer.WriteHeader(http.StatusOK)
		return
	}
	enc := json.NewEncoder(writer)
	type genreId struct {
		ID string `json:"id"`
	}
	for _, id := range ids {
		rID := genreId{ID: id}
		err := enc.Encode(rID)
		if err != nil {
			log.Warnf("Couldn't encode to HTTP response genre with uuid=%s %v\n", id, err)
			continue
		}
	}
}

func (h *genresHandler) reload(writer http.ResponseWriter, req *http.Request) {
	err := h.service.reload()
	if err != nil {
		log.Warnf("Problem reloading terms from TME: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
	}
}
