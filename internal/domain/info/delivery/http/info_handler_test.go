package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfoHandler_ServeHTTP(t *testing.T) {
	infoHandler := NewInfoHandler()

	req := &http.Request{}
	w := httptest.NewRecorder()

	infoHandler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	assert.Empty(t, w.Body.String())
}
