package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

var difficultLevel = map[string]int{
	"0": 30,
	"1": 50,
	"2": 70,
}

type Sudoku struct {
	grid      [9][9]int
	userGrid  [9][9]int
	available [81][]int
}

func displayGrid(grid [9][9]int) {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			fmt.Printf("%v", grid[i][j])
			if (j+1)%3 == 0 {
				fmt.Printf("|")
			}
		}
		fmt.Println()
	}
}

func checkViolation(grid [9][9]int, row, column, element int) bool {

	for i := 0; i < 9; i++ {

		if i != column && grid[row][i] == element {
			return true
		}
		if i != row && grid[i][column] == element {
			return true
		}
	}

	blockCheck := checkBlockViolation(row, column, element, grid)
	if blockCheck {
		return true
	}
	return false

}

// For block violation
func checkBlockViolation(row, column, element int, grid [9][9]int) bool {
	var x int
	var y int

	x = (row / 3) * 3
	y = (column / 3) * 3

	// check block violation
	for i := x; i < x+3; i++ {
		for j := y; j < y+3; j++ {
			if element == grid[i][j] {
				return true
			}
		}
	}
	return false
}

//To convert row, column position into 0-80 position
func getIndex(row, column int) int {
	return (row*9 + column)
}

//To generate random number
func getRandomNumber(value int) int {

	return rand.Intn(value)
}

// To remove element from available if not needed for backtracking
func removeFromAvailable(available []int, element int) []int {
	var index int
	for i := 0; i < len(available); i++ {
		if available[i] == element {
			index = i
			break
		}
	}
	return append(available[:index], available[index+1:]...)
}

// For generating valid sudoku grid using backtracking
/*
- Start from 1st position in grid and select a random 1-9 number
- check for row, col, block violation. if happens then remove element from available slice for that position and try again
- If no violations then move to next position and repeat above
- If available slice of element is exausted then bracktrack to previous position and repeat above

*/
func (s *Sudoku) generateGrid() error {
	var j int
	var counter int = 0
	for i := 0; i < 9; i++ {
		j = 0
		for {
			counter++
			if j >= 9 {
				break
			}

			if j == -1 {
				return fmt.Errorf("lock")
			}

			index := getIndex(i, j)

			if len(s.available[index]) != 0 {
				randomIndex := getRandomNumber(len(s.available[index]))

				element := s.available[index][randomIndex]

				check := checkViolation(s.grid, i, j, element)
				if check == false {
					s.grid[i][j] = element
					s.available[index] = removeFromAvailable(s.available[index], element)
					j++
				} else {
					s.available[index] = removeFromAvailable(s.available[index], element)
				}
			} else {
				for i := 0; i < 9; i++ {
					s.available[index] = append(s.available[index], i+1)

				}
				b := j
				if b != 0 {
					s.grid[i][b-1] = 0
				}
				j = j - 1
			}
		}
	}
	return nil

}

// For generating random gap grid for user from actual generated grid
// To DO = If position is already blank then skip it by checking it is visited earlier or convert row,col position to index or use perm method in rand
func (s *Sudoku) getGridForUser(blankBoxes int) {
	s.userGrid = s.grid
	for i := 0; i < blankBoxes; i++ {
		row := getRandomNumber(9)
		col := getRandomNumber(9)
		s.userGrid[row][col] = 0
	}
}

// For initializing the available slice of element 1-9 for each position on grid
func (s *Sudoku) initializeAvailable() {
	for i := 0; i < 81; i++ {
		for j := 0; j < 9; j++ {
			s.available[i] = append(s.available[i], j+1)
		}
	}
}

//Converting 2d array to array of string for sending to UI
func getStringArray(userGrid [9][9]int) string {
	var str string
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			str = str + strconv.Itoa(userGrid[i][j])
		}
	}
	return str
}

// For checking if user won or lost
func (s *Sudoku) checkWin() bool {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if s.grid[i][j] != s.userGrid[i][j] {
				return false
			}
		}
	}
	return true
}

// To save score in database
func saveScore(userTime time.Duration, name string) {
	hours := int(userTime / time.Hour)
	minutes := int(userTime / time.Minute)
	seconds := int(userTime / time.Second)
	seconds = seconds - minutes*60
	current := time.Now()
	date := current.Format("2006-01-02")
	usertime := strconv.Itoa(hours) + ":" + strconv.Itoa(minutes) + ":" + strconv.Itoa(seconds)

	db, err := sql.Open("mysql", "mayur:mayur1092@tcp(127.0.0.1:3306)/sudoku")

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	sql := "INSERT INTO Scores(Name, Time, Date) VALUES (?,?,?)"

	insert, err := db.Query(sql, name, usertime, date)

	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
}

type Score struct {
	Name string `json:"Name"`
	Time string `json:"Time"`
}

func getTopScores() string {
	var top []Score
	db, err := sql.Open("mysql", "mayur:mayur1092@tcp(127.0.0.1:3306)/sudoku")

	// if there is an error opening the connection, handle it
	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()

	results, err := db.Query("SELECT Name, Time FROM Scores Order by Time LIMIT 5")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	for results.Next() {
		var tag Score
		// for each row, scan the result into our tag composite object
		err = results.Scan(&tag.Name, &tag.Time)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		top = append(top, tag)
	}
	jsonData, err := json.Marshal(top)
	if err != nil {
		fmt.Println("error:", err)
	}
	return string(jsonData)
}
