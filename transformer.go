package main

import (
	"encoding/base64"
	"encoding/xml"

	"github.com/pborman/uuid"
)

func transformGenre(tmeTerm term, taxonomyName string) genre {
	tmeIdentifier := buildTmeIdentifier(tmeTerm.RawID, taxonomyName)
	uuid := uuid.NewMD5(uuid.UUID{}, []byte(tmeIdentifier)).String()

	return genre{
		UUID:                   uuid,
		PrefLabel:              tmeTerm.CanonicalName,
		AlternativeIdentifiers: alternativeIdentifiers{TME: []string{tmeIdentifier}, Uuids: []string{uuid}},
		Types: genreTypes,
	}
}

func buildTmeIdentifier(rawID string, tmeTermTaxonomyName string) string {
	id := base64.StdEncoding.EncodeToString([]byte(rawID))
	taxonomyName := base64.StdEncoding.EncodeToString([]byte(tmeTermTaxonomyName))
	return id + "-" + taxonomyName
}

type genreTransformer struct {
}

func (*genreTransformer) UnMarshallTaxonomy(contents []byte) ([]interface{}, error) {
	taxonomy := taxonomy{}
	err := xml.Unmarshal(contents, &taxonomy)
	if err != nil {
		return nil, err
	}
	interfaces := make([]interface{}, len(taxonomy.Terms))
	for i, d := range taxonomy.Terms {
		interfaces[i] = d
	}
	return interfaces, nil
}

func (*genreTransformer) UnMarshallTerm(content []byte) (interface{}, error) {
	dummyTerm := term{}
	err := xml.Unmarshal(content, &dummyTerm)
	if err != nil {
		return term{}, err
	}
	return dummyTerm, nil
}
