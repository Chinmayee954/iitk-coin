
package main

import (
	"log"
	"net/http"
	// "database/sql"
	 _ "github.com/mattn/go-sqlite3"
	//  "fmt"
	//  "strconv"
)

func main() {

// database, _ := sql.Open("sqlite3", "./signup_rollno.db")
// statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS SignUpRoll (id INTEGER PRIMARY KEY, rollno TEXT, password TEXT)")
//  statement.Exec()


   

	//   statement1, _ = database1.Prepare("DELETE FROM people WHERE id=?")
	// //   delete from userinfo where uid=?
    //  statement1.Exec(3)
    // rows, _ := database1.Query("SELECT id, rollno, coins FROM people")
    // var id int
    // var rollno string
    // var coins string
    // for rows.Next() {
    //     rows.Scan(&id, &rollno, &coins)
    //     fmt.Println(strconv.Itoa(id) + ": " + rollno + " " + coins)
	// }
 


	
	http.HandleFunc("/login", Login)
	http.HandleFunc("/secretpage", Home)
	http.HandleFunc("/refresh", Refresh)
	http.HandleFunc("/signup", SignUp)
	http.HandleFunc("/coins", AwardCoins)
	http.HandleFunc("/transaction", DoTransaction)
	http.HandleFunc("/wallet", Showcoins)
	http.HandleFunc("/history", History)
	http.HandleFunc("/makeadmin", MakeAdmin)
	log.Fatal(http.ListenAndServe(":8080", nil))

	
}