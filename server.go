package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/urfave/negroni"
)

var dev = flag.Bool("dev", false, "developement mode")

var db, err = sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@tcp(127.0.0.1:3306)/"+DB_NAME)

// Handling routes
func initRouter() (router *mux.Router) {
	router = mux.NewRouter()

	router.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./assets/"))))
	router.HandleFunc("/", homeHandler).Methods(http.MethodGet)
	router.HandleFunc("/ws", newGameHandler).Methods(http.MethodGet)

	return
}

//Handler for render the sudoku page
func homeHandler(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "game.html")
}

//Handler for new game
func newGameHandler(rw http.ResponseWriter, req *http.Request) {
	// Start Timer for current game
	start := time.Now()

	type Score struct {
		name string
		time []int
	}
	c, err := upgrader.Upgrade(rw, req, nil)
	if err != nil {
		log.Print("Upgrade : ", err)
	}

	// To get difficuly level from UI
	_, recvLevel, err := c.ReadMessage()
	if err != nil {
		fmt.Println(err)
	}
	blankBoxes := difficultLevel[string(recvLevel)]
	c.WriteMessage(websocket.TextMessage, []byte(getTopScores(db)))

	s := Sudoku{}
	s.initializeAvailable()
	err = s.generateGrid()
	if err != nil {
		fmt.Println("Lock")
		s.grid = [9][9]int{}
		s.available = [81][]int{}
		s.initializeAvailable()
		s.generateGrid()
	}

	s.getGridForUser(blankBoxes)

	if *dev {
		fmt.Println("Answer")
		displayGrid(s.grid)
	}

	str := getStringArray(s.userGrid)

	// Send userGrid to UI
	c.WriteMessage(websocket.TextMessage, []byte(str))

	for {
		// score := Score{}
		var userData map[string]int

		_, recvData, err := c.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}

		//Extracting data from UI
		_ = json.Unmarshal(recvData, &userData)
		value := userData["value"]
		row := userData["row"]
		col := userData["col"]

		//To Do - create func for directly checking by comparing grid position value and user entered value
		blockCheck := checkViolation(s.userGrid, row, col, value)
		if blockCheck {
			c.WriteMessage(websocket.TextMessage, []byte("violation"))
		} else {
			s.userGrid[row][col] = value
			win := s.checkWin()
			if win {
				c.WriteMessage(websocket.TextMessage, []byte("win"))
				userTiming := time.Since(start)
				// Getting player name
				_, nameData, _ := c.ReadMessage()
				name := string(nameData)
				saveScore(db, userTiming, name)
				break
			}
		}
	}

}

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	// var grid = [9][9]int{
	// 	{5, 0, 0, 0, 2, 7, 0, 0, 1},
	// 	{8, 2, 0, 0, 0, 0, 0, 7, 5},
	// 	{6, 0, 2, 0, 3, 0, 9, 4, 0},
	// 	{1, 5, 0, 4, 9, 0, 0, 0, 3},
	// 	{0, 8, 0, 7, 0, 0, 0, 0, 9},
	// 	{0, 0, 0, 2, 1, 8, 0, 0, 0},
	// 	{4, 0, 0, 9, 0, 2, 0, 0, 7},
	// 	{9, 2, 8, 3, 0, 0, 0, 1, 6},
	// 	{0, 6, 3, 1, 8, 5, 0, 0, 0},
	// }

	router := initRouter()
	server := negroni.Classic()
	server.UseHandler(router)

	server.Run(":3000")
	defer db.Close()

}
