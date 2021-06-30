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
    InsertedCoins int `json:"coins"`
}


func AwardCoins(response http.ResponseWriter, request *http.Request)  {

   CheckjwtToken(response, request)

	response.Header().Set("Content-Type","application/json")
	  var coins Coins
	  json.NewDecoder(request.Body).Decode(&coins)
      
	   database, _ := sql.Open("sqlite3", "./data.db")

	   

	   fmt.Println(coins.RollNo, coins.InsertedCoins)

	//    stmtupdate, _ := database.Prepare("UPDATE userdata set coins=? WHERE rollno = ?")
    //       stmtupdate.Exec(100, coins.RollNo)
		    
	      intialrow, _ := database.Query("SELECT id, rollno, coins FROM userdata")

	  var id int
     var rollno string
     var intialcoins int
	 var flag = 0

	   for intialrow.Next() {
          intialrow.Scan(&id, &rollno, &intialcoins)
		  fmt.Println(rollno)
	  
		 if coins.RollNo == rollno {
			 
			 fmt.Println("found")
             flag = 1
		 }
		}

		if flag == 1{
			  stmtupdate, _ := database.Prepare("UPDATE userdata set coins=? WHERE rollno = ?")
          stmtupdate.Exec(coins.InsertedCoins, coins.RollNo)
		} else {
			response.Write([]byte(fmt.Sprintf("NO DATA AVAILABLE")))
		}


	     updatedrow, _ := database.Query("SELECT id, rollno, coins FROM userdata")
	
	var id1 int
    var rollno1 string
    var coins1 int

	   for updatedrow.Next() {
         updatedrow.Scan(&id1, &rollno1, &coins1)
	 	 fmt.Println(strconv.Itoa(id1) + ": " + rollno1 + " " + strconv.Itoa(coins1))

	 	}

	}
