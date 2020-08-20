package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/urfave/negroni"

	"github.com/gorilla/mux"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

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

var LEVEL int = 1

func checkViolation(grid [9][9]int, row, column, element int) bool {
	var x int
	var y int
	// load row
	var gridrow [9]int = grid[row]
	var gridcolumn [9]int

	// load column
	for i := 0; i < 9; i++ {
		gridcolumn[i] = grid[i][column]
	}

	// fmt.Println("Row : ", gridrow)
	// fmt.Println("Column : ", gridcolumn)

	// check row or column violation
	for i := 0; i < 9; i++ {
		if element == gridrow[i] {
			if i != column {
				// fmt.Println("Row violation")
				return true
			}
		}
		if element == gridcolumn[i] {
			if i != row {
				// fmt.Println("Column violation")
				return true
			}
		}
	}

	switch {
	case row < 3 && column < 3:
		x = 0
		y = 0
		// fmt.Println("Block 1")
	case row < 3 && (column > 2 && column < 6):
		x = 0
		y = 3
		// fmt.Println("Block 2")
	case row < 3 && (column > 5 && column < 9):
		x = 0
		y = 6
		// fmt.Println("Block 3")
	case (row > 2 && row < 6) && (column < 3):
		x = 3
		y = 0
		// fmt.Println("Block 4")
	case (row > 2 && row < 6) && (column > 2 && column < 6):
		x = 3
		y = 3
		// fmt.Println("Block 5")
	case (row > 2 && row < 6) && (column > 5 && column < 9):
		x = 3
		y = 6
		// fmt.Println("Block 6")
	case (row > 5) && (column < 3):
		x = 6
		y = 0
		// fmt.Println("Block 7")
	case (row > 5) && (column > 2 && column < 6):
		x = 6
		y = 3
		// fmt.Println("Block 8")
	case (row > 5) && (column > 5 && column < 9):
		x = 6
		y = 6
		// fmt.Println("Block 9")

	}

	// check block violation
	for i := x; i < x+3; i++ {
		for j := y; j < y+3; j++ {
			if element == grid[i][j] {
				// fmt.Println("Block Violation")
				return true
			}
		}
	}
	return false

}

// }1,1
// (1*9 + 1)
func getIndex(row, column int) int {
	return (row*9 + column)
}

func getRandomNumber(value int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(value)
}

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

func generateGrid(grid [9][9]int, available [81][]int) ([9][9]int, error) {
	// for
	var j int
	// var k int
	var counter int = 0
	for i := 0; i < 9; i++ {
		j = 0
		for {
			counter++
			if j >= 9 {
				break
			}

			if j == -1 {
				// j = 0
				// break
				return grid, fmt.Errorf("Gridlock")
			}

			index := getIndex(i, j)

			if len(available[index]) != 0 {
				randomIndex := getRandomNumber(len(available[index]))
				// fmt.Printf("Random Index %v %v ", index, randomIndex)

				element := available[index][randomIndex]

				// fmt.Println(element)
				check := checkViolation(grid, i, j, element)
				if check == false {
					grid[i][j] = element
					// displayGrid(grid)
					// fmt.Println(element)
					// new_grid = append(new_grid, element)
					available[index] = removeFromAvailable(available[index], element)
					j++
					// break
				} else {
					available[index] = removeFromAvailable(available[index], element)
				}
			} else {
				for i := 0; i < 9; i++ {
					available[index] = append(available[index], i+1)

				}
				// if j == -1 {
				// 	// j = 0
				// 	// break
				// 	return grid, fmt.Errorf("Gridlock")
				// }
				b := j

				if b != 0 {
					grid[i][b-1] = 0
				}

				// fmt.Println("available :", available[index])
				j = j - 1
			}
		}

	}
	fmt.Println("Counter :", counter)
	fmt.Println()

	return grid, nil

}

func getGridForUser(grid [9][9]int, blankBoxes int) [9][9]int {
	for i := 0; i < blankBoxes; i++ {
		row := getRandomNumber(9)
		col := getRandomNumber(9)
		grid[row][col] = 0
	}
	return grid
}

func (s *Sudoku) initializeAvailable() {
	for i := 0; i < 81; i++ {
		for j := 0; j < 9; j++ {
			s.available[i] = append(s.available[i], j+1)
		}
	}
}

func InitRouter() (router *mux.Router) {
	router = mux.NewRouter()

	router.HandleFunc("/", homeHandler).Methods(http.MethodGet)
	router.HandleFunc("/ws", newGameHandler).Methods(http.MethodGet)

	return
}

func homeHandler(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "home.html")
}

func newGameHandler(rw http.ResponseWriter, req *http.Request) {
	c, err := upgrader.Upgrade(rw, req, nil)
	if err != nil {
		log.Print("Upgrade : ", err)
	}
	s := Sudoku{}
	// s.initializeAvailable()

	c.WriteMessage(websocket.TextMessage, []byte("Hello from Server"))
}

var homeTemplate = template.Must

type Sudoku struct {
	grid      [9][9]int
	available [81][]int
}

func main() {
	// var available [81][]int
	s := Sudoku{}

	// for i := 0; i < 81; i++ {
	// 	for j := 0; j < 9; j++ {
	// 		s.available[i] = append(s.available[i], j+1)
	// 	}
	// }

	s.initializeAvailable()
	// for i := 0; i < 81; i++ {
	// 	for j := 0; j < 9; j++ {
	// 		fmt.Printf("%v", available[i][j])
	// 		// if (j+1)%3 == 0 {
	// 		// 	fmt.Printf("|")
	// 		// }
	// 	}
	// 	fmt.Println()
	// }

	// fmt.Println("available len", len(available))
	// var grid = [9][9]int{
	// 	{5, 0, 0, 0, 2, 7, 0, 0, 1},
	// 	{8, 0, 0, 0, 0, 0, 0, 7, 5},
	// 	{6, 0, 2, 0, 3, 0, 9, 4, 0},
	// 	{1, 5, 0, 4, 9, 0, 0, 0, 3},
	// 	{0, 8, 0, 7, 0, 0, 0, 0, 9},
	// 	{0, 0, 0, 2, 1, 8, 0, 0, 0},
	// 	{4, 0, 0, 9, 0, 2, 0, 0, 7},
	// 	{9, 2, 8, 3, 0, 0, 0, 1, 6},
	// 	{0, 6, 3, 1, 8, 5, 0, 0, 0},
	// }

	// var grid = [9][9]int{}
	// var userGrid [9][9]int
	// for i := 0; i < 100; i++ {

	fmt.Println("Welcome to sudoku")
	// displayGrid(grid)
	var err error
	s.grid, err = generateGrid(s.grid, s.available)
	// time.Sleep(5 * time.Second)
	if err != nil {
		s.grid = [9][9]int{}
		fmt.Println(err)
		s.grid, err = generateGrid(s.grid, s.available)
		// time.Sleep(5 * time.Second)
	}

	displayGrid(s.grid)

	// userGrid = getGridForUser(grid, 60)
	// fmt.Println("User Grid")
	// displayGrid(userGrid)

	router := InitRouter()
	server := negroni.Classic()
	server.UseHandler(router)

	server.Run(":8000")

}
