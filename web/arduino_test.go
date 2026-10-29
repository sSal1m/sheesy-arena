package web

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestArduinoPanelAndIsolatedSimulation(t *testing.T) {
	web := setupTestWeb(t)
	assert.Equal(t, 200, web.getHttpResponse("/panels/arduino").Code)
	assert.Contains(t, web.getHttpResponse("/").Body.String(), "/panels/arduino")
	status := web.getHttpResponse("/panels/arduino/status")
	assert.Equal(t, 200, status.Code)
	assert.NotContains(t, status.Body.String(), "Key")
	before := web.arena.ScoreResultSnapshot()
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/panels/arduino/simulate", strings.NewReader(`{"Alliance":"red","Phase":"auto","Count":5}`))
	web.newHandler().ServeHTTP(recorder, req)
	assert.Equal(t, 200, recorder.Code)
	var result map[string]any
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &result))
	assert.Equal(t, true, result["Simulation"])
	assert.Equal(t, before, web.arena.ScoreResultSnapshot())
	recorder = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/panels/arduino/simulate", strings.NewReader(`{"Alliance":"red","Phase":"auto","Count":1,"MatchId":0}`))
	web.newHandler().ServeHTTP(recorder, req)
	assert.Equal(t, 400, recorder.Code)
	web.arena.EventSettings.AdminPassword = "private"
	for _, path := range []string{"/panels/arduino", "/panels/arduino/status"} {
		assert.Equal(t, http.StatusTemporaryRedirect, web.getHttpResponse(path).Code)
	}
	recorder = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/panels/arduino/simulate", strings.NewReader(`{"Alliance":"red","Phase":"auto","Count":1}`))
	web.newHandler().ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusTemporaryRedirect, recorder.Code)
}
