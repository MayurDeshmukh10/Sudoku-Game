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
### Setting Enviroment Variables

```
Add following in bashrc file
export DATABASE_USERNAME="Your database username"
export DATABASE_PASSWORD="Your database password"
export DATABASE_NAME="Your database name"
```

### To run
```
go build *.go
./server

Go to - localhost:9000
```
