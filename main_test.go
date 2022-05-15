package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {

	router := setUpRouter()

	w := httptest.NewRecorder()

	httpReq, _ := http.NewRequest("GET", "/ping", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, 200, w.Code)

	assert.Equal(t, "pong", w.Body.String())

}
