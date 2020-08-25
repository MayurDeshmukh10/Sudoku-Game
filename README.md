# Sudoku-Game

### Dependencies
```
go get -u github.com/go-sql-driver/mysql
go get -u github.com/gorilla/mux
go get github.com/gorilla/websocket
go get github.com/urfave/negroni
```

### Import MySQL Database
```
CREATE DATABASE newdatabase;
mysql -u [username] -p newdatabase < database.sql
```
### To run
go run sudoku.go

Go to - localhost:9000
```
