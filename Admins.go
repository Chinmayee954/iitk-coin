package main

import (
    "encoding/json"
	"fmt"
	"net/http"
	// "time"
	  "strconv"
	//   "context"
	//   "log"

	// "github.com/dgrijalva/jwt-go"
     _ "github.com/mattn/go-sqlite3"
	 "database/sql"
	//    "golang.org/x/crypto/bcrypt"
)


type AdminRollNo struct{
    Rollno string `json:"rollno"`
}

func MakeAdmin(response http.ResponseWriter, request *http.Request) {

   
	Username := CheckjwtToken(response, request)
	


	 response.Header().Set("Content-Type","application/json")
 	  var adminrollno AdminRollNo
 	  json.NewDecoder(request.Body).Decode(&adminrollno)
	   database , _ := sql.Open("sqlite3", "./admins.db")



	   if IsAdmin(Username) == "true" {

		   fmt.Println(IsAdmin(Username))
		fmt.Println("You are in")
		
        statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS adminrecords (id INTEGER PRIMARY KEY, rollno TEXT, coins INTEGER)")
        statement.Exec()

		
		  coins := ShowAdminCoins(adminrollno.Rollno)
		fmt.Println(adminrollno.Rollno)
		  fmt.Println(coins)
		 stmt, _ := database.Prepare("INSERT INTO adminrecords (rollno, coins) VALUES (?, ?)")
		   stmt.Exec(adminrollno.Rollno, coins)
	   } else {
	  	 response.Write([]byte(fmt.Sprintf("You do not have rights to access this page")))
      }


	
        row1, _ := database.Query("SELECT id, rollno, coins FROM adminrecords")
     
	var id int	
	var rollno string
    var coins1 int


	   for row1.Next() {
          row1.Scan(&id, &rollno, &coins1)
             fmt.Println(strconv.Itoa(id) + ": " + rollno + " " + strconv.Itoa(coins1))
		}
}


func IsAdmin(rollno string) string {
     	database , _ := sql.Open("sqlite3", "./admins.db")
        row, _ := database.Query("SELECT id, rollno, coins FROM adminrecords")

		var id int	
	var adminrollno string
    var coins int
	var flag = 0

	   for row.Next() {
          row.Scan(&id, &adminrollno, &coins)
		   if rollno == adminrollno {
             flag = 1
		   }
		  
		}

		if flag == 1 {
			return "true"
		} else {
			return "false"
		}
     
}

func ShowAdminCoins(rollno string) int{

	  database, _ := sql.Open("sqlite3", "./data.db")
      row, _ := database.Query("SELECT id, rollno, coins FROM people")

	var requestedrollno string
    var coins int
    var id int

	var returncoins int

	   for row.Next() {
          row.Scan(&id, &requestedrollno, &coins)
		 if rollno == requestedrollno {
			 returncoins = coins
		 }
		}
  
		return returncoins

}