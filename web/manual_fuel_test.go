package web

import (
	"github.com/Team254/cheesy-arena/field"
	"github.com/Team254/cheesy-arena/game"
	"github.com/Team254/cheesy-arena/model"
	"github.com/Team254/cheesy-arena/websocket"
	"github.com/google/uuid"
	gorilla "github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func readFuelEvent(t *testing.T, conn *gorilla.Conn, kind string) websocket.Message {
	t.Helper()
	require.NoError(t, conn.SetReadDeadline(time.Now().Add(3*time.Second)))
	for {
		var message websocket.Message
		require.NoError(t, conn.ReadJSON(&message))
		if message.Type == kind {
			return message
		}
	}
}

func manualWebCommand(web *Web, delta int) model.FuelCommand {
	state := web.arena.FuelState()
	return model.FuelCommand{RequestId: uuid.NewString(), AttemptId: state.AttemptId, MatchId: state.MatchId, Alliance: "red", Phase: "auto", Delta: delta}
}

func startManualWebMatch(t *testing.T, web *Web) {
	command := manualWebCommand(web, 0)
	command.Mode = model.FuelManual
	require.Empty(t, web.arena.ApplyFuelCommand("setFuelMode", command, "test").Error)
	web.arena.MatchState = field.AutoPeriod
	web.arena.MatchStartTime = time.Now()
	web.arena.Update()
}

func TestManualFuelPanelAndAuthorization(t *testing.T) {
	web := setupTestWeb(t)
	page := web.getHttpResponse("/panels/manual_fuel")
	require.Equal(t, 200, page.Code, page.Body.String())
	assert.Contains(t, page.Body.String(), "red-auto-value")
	assert.Contains(t, page.Body.String(), "blue-teleop-value")
	assert.Contains(t, web.getHttpResponse("/").Body.String(), "/panels/manual_fuel")
	web.arena.EventSettings.AdminPassword = "secret"
	assert.Equal(t, 307, web.getHttpResponse("/panels/manual_fuel").Code)
	assert.Equal(t, 307, web.getHttpResponse("/panels/manual_fuel/history").Code)
	server, url := web.startTestServer()
	defer server.Close()
	_, response, err := gorilla.DefaultDialer.Dial(url+"/panels/manual_fuel/websocket", nil)
	require.Error(t, err)
	assert.Equal(t, 307, response.StatusCode)
}

func TestManualFuelWebsocketSharedStateReconnectAndAudit(t *testing.T) {
	web := setupTestWeb(t)
	startManualWebMatch(t, web)
	server, url := web.startTestServer()
	defer server.Close()
	connect := func() *gorilla.Conn {
		conn, _, err := gorilla.DefaultDialer.Dial(url+"/panels/manual_fuel/websocket", nil)
		require.NoError(t, err)
		readFuelEvent(t, conn, "manualFuelState")
		return conn
	}
	a, b := connect(), connect()
	defer a.Close()
	defer b.Close()
	command := manualWebCommand(web, 1)
	require.NoError(t, a.WriteJSON(websocket.Message{Type: "adjustFuel", Data: command}))
	ack := readFuelEvent(t, a, "fuelResult").Data.(map[string]any)
	assert.Empty(t, ack["Error"])
	state := readFuelEvent(t, b, "manualFuelState").Data.(map[string]any)
	assert.Equal(t, float64(1), state["Red"].(map[string]any)["Auto"])
	require.NoError(t, a.WriteJSON(websocket.Message{Type: "adjustFuel", Data: command}))
	ack = readFuelEvent(t, a, "fuelResult").Data.(map[string]any)
	assert.Equal(t, true, ack["Duplicate"])
	c := connect()
	defer c.Close()
	require.NoError(t, c.WriteJSON(websocket.Message{Type: "getFuelState", Data: model.FuelCommand{}}))
	state = readFuelEvent(t, c, "manualFuelState").Data.(map[string]any)
	assert.Equal(t, float64(1), state["Red"].(map[string]any)["Auto"])
	history := web.getHttpResponse("/panels/manual_fuel/history?attemptId=" + command.AttemptId)
	assert.Equal(t, 200, history.Code)
	assert.Contains(t, history.Body.String(), command.RequestId)
	assert.Contains(t, history.Body.String(), "open-access")
}

func TestManualFuelRejectsExpiredSessionAndFractionalDelta(t *testing.T) {
	web := setupTestWeb(t)
	startManualWebMatch(t, web)
	web.arena.EventSettings.AdminPassword = "secret"
	session := &model.UserSession{Token: uuid.NewString(), Username: "admin", CreatedAt: time.Now()}
	require.NoError(t, web.arena.Database.CreateUserSession(session))
	server, url := web.startTestServer()
	defer server.Close()
	headers := http.Header{"Cookie": []string{"session_token=" + session.Token}}
	conn, _, err := gorilla.DefaultDialer.Dial(url+"/panels/manual_fuel/websocket", headers)
	require.NoError(t, err)
	defer conn.Close()
	readFuelEvent(t, conn, "manualFuelState")
	require.NoError(t, conn.WriteJSON(websocket.Message{Type: "adjustFuel", Data: map[string]any{"Delta": 1.5}}))
	assert.Contains(t, readFuelEvent(t, conn, "error").Data, "Invalid command fields")
	require.NoError(t, web.arena.Database.TruncateUserSessions())
	require.NoError(t, conn.WriteJSON(websocket.Message{Type: "adjustFuel", Data: manualWebCommand(web, 1)}))
	ack := readFuelEvent(t, conn, "fuelResult").Data.(map[string]any)
	assert.Contains(t, ack["Error"], "expired")
	assert.Zero(t, web.arena.FuelState().Red.Auto)
}

func TestManualFuelResultPublishingAndReview(t *testing.T) {
	web := setupTestWeb(t)
	match := &model.Match{Type: model.Qualification, TypeOrder: 1, ShortName: "Q1", LongName: "Qualification 1", Red1: 111, Red2: 112, Red3: 113, Blue1: 211, Blue2: 212, Blue3: 213}
	for _, id := range []int{111, 112, 113, 211, 212, 213} {
		require.NoError(t, web.arena.Database.CreateTeam(&model.Team{Id: id}))
	}
	require.NoError(t, web.arena.Database.CreateMatch(match))
	require.NoError(t, web.arena.LoadMatch(match))
	startManualWebMatch(t, web)
	for i := 0; i < 20; i++ {
		require.Empty(t, web.arena.ApplyFuelCommand("adjustFuel", manualWebCommand(web, 5), "admin").Error)
	}
	web.arena.MatchState = field.PostMatch
	web.arena.MatchStartTime = time.Now().Add(-game.GetDurationToTeleopEnd() - time.Second)
	web.arena.Update()
	require.NoError(t, web.commitCurrentMatchScore())
	result, err := web.arena.Database.GetMatchResultForMatch(match.Id)
	require.NoError(t, err)
	require.NotNil(t, result.FuelControl)
	assert.Equal(t, 100, result.RedScoreSummary().AutoFuelPoints)
	assert.True(t, result.RedScoreSummary().EnergizedBonusRankingPoint)
	rankings, err := web.arena.Database.GetAllRankings()
	require.NoError(t, err)
	require.Len(t, rankings, 6)
	assert.Equal(t, 4, rankings[0].RankingPoints)
	assert.Equal(t, 100, web.arena.SavedMatchResult.RedScoreSummary().Score)
	assert.NotEmpty(t, web.arena.ApplyFuelCommand("adjustFuel", manualWebCommand(web, 1), "admin").Error)
	// Historical review keeps source metadata even when the client omits it.
	metadata := result.FuelControl
	result.RedScore.Hub.ShiftCounts[game.ShiftAuto] = 95
	require.NoError(t, model.ReconcileFuelMetadata(result, metadata))
	require.NoError(t, web.commitMatchScore(match, result, true))
	result, err = web.arena.Database.GetMatchResultForMatch(match.Id)
	require.NoError(t, err)
	assert.Equal(t, 95, result.RedScoreSummary().Score)
	assert.Equal(t, 95, result.FuelControl.Red.Adjustment[game.ShiftAuto])
}

func TestManualFuelDoesNotSerializeLiveScoreAliases(t *testing.T) {
	web := setupTestWeb(t)
	startManualWebMatch(t, web)
	snapshot := web.getCurrentMatchResult()
	assert.NotNil(t, snapshot.RedScore.Fouls) // The existing result editor expects an array.
	require.Empty(t, web.arena.ApplyFuelCommand("adjustFuel", manualWebCommand(web, 1), "a").Error)
	assert.Zero(t, snapshot.RedScoreSummary().Score)
	request := httptest.NewRequest("GET", "/panels/manual_fuel", nil)
	recorder := httptest.NewRecorder()
	web.newHandler().ServeHTTP(recorder, request)
	assert.Equal(t, 200, recorder.Code)
}
