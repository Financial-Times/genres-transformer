package main

type genre struct {
	UUID          string `json:"uuid"`
	CanonicalName string `json:"canonicalName"`
	TmeIdentifier string `json:"tmeIdentifier,omitempty"`
	Type          string `json:"type"`
}

type genreLink struct {
	APIURL string `json:"apiUrl"`
}
