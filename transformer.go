package main

import (
	"github.com/pborman/uuid"
)

func transformGenre(t term) genre {
	return genre{
		UUID:          uuid.NewMD5(uuid.UUID{}, []byte(t.ID)).String(),
		CanonicalName: t.CanonicalName,
		TmeIdentifier: t.ID,
		Type:          "Genre",
	}
}
