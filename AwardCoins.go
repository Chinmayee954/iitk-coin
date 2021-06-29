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

type Coins struct{
    RollNo string `json:"rollno"`
    InsertedCoins string `json:"coins"`
}


func AwardCoins(response http.ResponseWriter, request *http.Request)  {

   CheckjwtToken(response, request)

	response.Header().Set("Content-Type","application/json")
	  var coins Coins
	  json.NewDecoder(request.Body).Decode(&coins)
      
	   database, _ := sql.Open("sqlite3", "./data.db")

	      intialrow, _ := database.Query("SELECT id, rollno, coins FROM people")


	   fmt.Println(coins.RollNo, coins.InsertedCoins)
	  

	     var initialid int
    var initialrollno string
    var initialcoins string
	 var flag int

	   for intialrow.Next() {
          intialrow.Scan(&initialid, &initialrollno, &initialcoins)
		 if coins.RollNo == initialrollno {
			 
			 flag = 1
		    
		//   return
		 }
		}

		if flag == 1 {
			fmt.Println("found")
               stmtupdate, _ := database.Prepare("UPDATE people set coins=? WHERE rollno = (?)")
          stmtupdate.Exec(coins.InsertedCoins, coins.RollNo)
		} else {
			fmt.Println("notfound")
                stmtinsert, _ := database.Prepare("INSERT INTO people (rollno, coins) VALUES (?, ?)")
	             stmtinsert.Exec(coins.RollNo, coins.InsertedCoins)
	
		} 

	  
       

		

	     updatedrow, _ := database.Query("SELECT id, rollno, coins FROM people")
	
	var id int
    var rollno string
    var updatedcoins string

	   for updatedrow.Next() {
         updatedrow.Scan(&id, &rollno, &updatedcoins)
	 	 fmt.Println(strconv.Itoa(id) + ": " + rollno + " " + updatedcoins)

	 	}

	}
