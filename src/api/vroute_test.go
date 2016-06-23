package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVRoute_buildUrl(t *testing.T) {
	router := NewRouter()

	router.Handle("/fc/{id:[a-zA-Z0-9._-]+}", NewHandlerFromString("a b")).Methods("GET").Name("op")

	url, err := router.Get("op").URLPath("id", "123")

	assert.NoError(t, err, "Build url")
	assert.Equal(t, url.String(), "/fc/123", "Build url")

	_, err = router.Get("op").URLPath("id", "123!")

	assert.NotEmpty(t, err, "Build invalid url")
}
