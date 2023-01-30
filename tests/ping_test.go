package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/marcinlovescode/go-clean-fileupload/internal/app"
	"github.com/marcinlovescode/go-clean-fileupload/internal/pkg/logger"
	"github.com/marcinlovescode/go-clean-fileupload/tests/doubles"
)

func createStubLogger() logger.Logger {
	return &doubles.VoidLogger{}
}

func TestPingRoute(t *testing.T) {

	router := app.CreateHttpHandlers(createStubLogger())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/files/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Pong", w.Body.String())
}
