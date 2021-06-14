
package main

import (
	"log"
	"net/http"
	"database/sql"
	 _ "github.com/mattn/go-sqlite3"
	 "fmt"
	 "strconv"
)

func main() {

database, _ := sql.Open("sqlite3", "./signup_rollno.db")
statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS SignUpRoll (id INTEGER PRIMARY KEY, rollno TEXT, password TEXT)")
 statement.Exec()


    database1, _ := sql.Open("sqlite3", "./data.db")
    statement1, _ := database1.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, rollno TEXT, coins TEXT)")
    statement1.Exec()

	 statement1, _ = database1.Prepare("INSERT INTO people (rollno, coins) VALUES (?, ?)")
    statement1.Exec("190232", "20")
    rows, _ := database1.Query("SELECT id, rollno, coins FROM people")
    var id int
    var rollno string
    var coins string
    for rows.Next() {
        rows.Scan(&id, &rollno, &coins)
        fmt.Println(strconv.Itoa(id) + ": " + rollno + " " + coins)
	}
 


	
	http.HandleFunc("/login", Login)
	http.HandleFunc("/secretpage", Home)
	http.HandleFunc("/refresh", Refresh)
	http.HandleFunc("/signup", SignUp)
	log.Fatal(http.ListenAndServe(":8080", nil))

	
}