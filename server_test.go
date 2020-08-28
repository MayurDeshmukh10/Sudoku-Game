package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestHomeHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "localhost:3000", nil)
	if err != nil {
		t.Fatalf("Could not create request : %v", err)
	}
	rec := httptest.NewRecorder()

	homeHandler(rec, req)

	res := rec.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status ok; got %v", res.StatusCode)
	}
}

func TestNewGameHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(newGameHandler))
	defer server.Close()

	t.Run("Should return valid json of leaderboard", func(t *testing.T) {

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)

		err = ws.WriteMessage(websocket.TextMessage, []byte("0"))
		if err != nil {
			t.Fatalf("could not send message ws connection %v", err)
		}
		_, recvData, err := ws.ReadMessage()
		if err != nil {
			t.Fatalf("could read message from server due to %v", err)
		}
		fmt.Println(string(recvData))
		ws.Close()
	})

	t.Run("Should generate grid for difficulty level easy", func(t *testing.T) {

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)

		if err != nil {
			t.Fatalf("could not connect to websocket due to %v", err)
		}

		err = ws.WriteMessage(websocket.TextMessage, []byte("0"))
		if err != nil {
			t.Fatalf("could not send message ws connection %v", err)
		}

		_, _, err = ws.ReadMessage()

		var flag int = 0
		_, recvData, err := ws.ReadMessage()

		if err != nil {
			t.Fatalf("could read message from server due to %v", err)
		}
		for _, value := range recvData {
			if string(value) == "0" {
				flag = 1
			}
		}
		if flag == 0 {
			t.Errorf("expected complete grid but got zero value in grid")
		}
		ws.Close()
	})

	t.Run("Should generate grid for difficulty level medium", func(t *testing.T) {
		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)

		if err != nil {
			t.Fatalf("could not connect to websocket due to %v", err)
		}

		err = ws.WriteMessage(websocket.TextMessage, []byte("1"))
		if err != nil {
			t.Fatalf("could not send message ws connection %v", err)
		}
		_, _, err = ws.ReadMessage()
		var flag int = 0
		_, recvData, err := ws.ReadMessage()
		if err != nil {
			t.Fatalf("could read message from server due to %v", err)
		}
		for _, value := range recvData {
			if string(value) == "0" {
				flag = 1
			}
		}
		if flag == 0 {
			t.Errorf("expected complete grid but got zero value in grid")
		}
		ws.Close()
	})

	t.Run("Should generate grid for difficulty level hard", func(t *testing.T) {

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)

		if err != nil {
			t.Fatalf("could not connect to websocket due to %v", err)
		}

		err = ws.WriteMessage(websocket.TextMessage, []byte("2"))
		if err != nil {
			t.Fatalf("could not send message ws connection %v", err)
		}
		_, _, err = ws.ReadMessage()
		var flag int = 0
		_, recvData, err := ws.ReadMessage()

		if err != nil {
			t.Fatalf("could read message from server due to %v", err)
		}
		for _, value := range recvData {
			if string(value) == "0" {
				flag = 1
			}
		}
		if flag == 0 {
			t.Errorf("expected complete grid but got zero value in grid")
		}

		ws.Close()
	})

}

func TestCheckWin(t *testing.T) {
	t.Run("Should return true if user won the game", func(t *testing.T) {
		s := Sudoku{}
		s.initializeAvailable()
		s.generateGrid()
		s.userGrid = s.grid
		winStatus := s.checkWin()
		assert.Equal(t, true, winStatus)
	})
	t.Run("Should return false if user has not won the game yet", func(t *testing.T) {
		s := Sudoku{}
		s.initializeAvailable()
		s.generateGrid()
		winStatus := s.checkWin()
		assert.Equal(t, false, winStatus)
	})
}

// func TestSaveScore(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("error in creating mock database due to %v", err)
// 	}
// 	defer db.Close()

// 	current := time.Now()
// 	date := current.Format("2006-01-02")
// 	prep := mock.ExpectPrepare("^INSERT INTO Scores*")

// 	prep.ExpectExec().
// 		WithArgs("Test", 100*time.Millisecond, date).
// 		WillReturnResult(sqlmock.NewResult(1, 1))

// 	saveScore(db, 100*time.Microsecond, "Test")

// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("there were unfulfilled expections: %s", err)
// 	}
// }
