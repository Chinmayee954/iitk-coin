package main

import (
	// "encoding/json"
	 "fmt"
	 "net/http"
	  "time"
	   "strconv"
	//   "context"
	//   "log"

	// "github.com/dgrijalva/jwt-go"
     _ "github.com/mattn/go-sqlite3"
	 "database/sql"
	//    "golang.org/x/crypto/bcrypt"
)


func History(response http.ResponseWriter, request *http.Request) {
	database , _ := sql.Open("sqlite3", "./history.db")
        row, _ := database.Query("SELECT id, rollno1, rollno2, coins, time FROM records")
     
	var id int	
	var rollno1 string
	var rollno2 string
    var coins string
    var time string

	   for row.Next() {
          row.Scan(&id, &rollno1,&rollno2, &coins, &time)
		   fmt.Println(strconv.Itoa(id) + ": " + rollno1 + " " + rollno2 + " " + coins + " " + time)
		}
}


func AddHistory(rollno1 string, rollno2 string, coins string) {
	database , _ := sql.Open("sqlite3", "./history.db")
statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS records (id INTEGER PRIMARY KEY, rollno1 TEXT, rollno2 TEXT, coins TEXT, time TEXT)")
statement.Exec()

stmt, _ := database.Prepare("INSERT INTO records (rollno1, rollno2, coins, time) VALUES (?, ?, ?, ?)")
recordtime := time.Now()
stmt.Exec(rollno1, rollno2, coins, recordtime);

}