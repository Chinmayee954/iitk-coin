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

type RemoveUser struct{
	Rollno string `json:"rollno"`
}

type FreezeUser struct{
	Rollno string `json:"rollno"`
}

type UserRewards struct {
	Rollno string `json:"rollno"`
	Rewards int `json:"rewards"`
}


func Admin(res http.ResponseWriter, request *http.Request){
	Username := CheckjwtToken(res, request)
	res.Header().Set("Content-Type","application/json")
 	// var adminrollno AdminRollNo
 	// json.NewDecoder(request.Body).Decode(&adminrollno)
	// database , _ := sql.Open("sqlite3", "./admins.db")
	if Username != "false" && IsAdmin(Username) == "true" {
		res.Write([]byte(fmt.Sprintf("Welcome Admin")))
	} else {
		res.Write([]byte(fmt.Sprintf("You are not an Admin")))
		return
	}
}

func MakeAdmin(response http.ResponseWriter, request *http.Request) { 
	Username := CheckjwtToken(response, request)
	response.Header().Set("Content-Type","application/json")
 	var adminrollno AdminRollNo
 	json.NewDecoder(request.Body).Decode(&adminrollno)
	database , _ := sql.Open("sqlite3", "./admins.db")
	if Username != "false" && IsAdmin(Username) == "true" {
		// fmt.Println(IsAdmin(Username))
		fmt.Println("You are in")
        statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS adminrecords (id INTEGER PRIMARY KEY, rollno TEXT, coins INTEGER)")
        statement.Exec()
		coins := ShowAdminCoins(adminrollno.Rollno)
		fmt.Println(adminrollno.Rollno)
		fmt.Println(coins)
		stmt, _ := database.Prepare("INSERT INTO adminrecords (rollno, coins) VALUES (?, ?)")
		stmt.Exec(adminrollno.Rollno, 0)
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
      row, _ := database.Query("SELECT id, rollno, coins FROM userdata")
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

func DeleteUser(response http.ResponseWriter, request *http.Request) {
	Username := CheckjwtToken(response, request)
	 response.Header().Set("Content-Type","application/json")
 	  var removeuser RemoveUser
 	  json.NewDecoder(request.Body).Decode(&removeuser)
	   database , _ := sql.Open("sqlite3", "./data.db")

	   if  Username != "false" && IsAdmin(Username) == "true" {
		statement, _ := database.Prepare("DELETE FROM userdata WHERE rollno = ?")
		statement.Exec(removeuser.Rollno)
	   }

}

func DeleteAdmin(response http.ResponseWriter, request *http.Request) {
	Username := CheckjwtToken(response, request)
	 response.Header().Set("Content-Type","application/json")
 	  var adminrollno AdminRollNo 
 	  json.NewDecoder(request.Body).Decode(&adminrollno)
	   database , _ := sql.Open("sqlite3", "./admins.db")

	   if  Username != "false" && IsAdmin(Username) == "true" {
		statement, _ := database.Prepare("DELETE FROM adminrecords WHERE rollno = ?")
		statement.Exec(adminrollno.Rollno)
	   }

}

func freezeuser(response http.ResponseWriter, request *http.Request) {
	Username := CheckjwtToken(response, request)
	 response.Header().Set("Content-Type","application/json")
 	  var freezerollno FreezeUser 
 	  json.NewDecoder(request.Body).Decode(&freezerollno)
	//    database , _ := sql.Open("sqlite3", "./admins.db")

	   if  Username != "false" && IsAdmin(Username) == "true" {
		database, _ := sql.Open("sqlite3", "./blacklist.db")
		statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS freezerecords (id INTEGER PRIMARY KEY, rollno TEXT,coins INT)")
		statement.Exec()
		coins := ShowUserCoins(freezerollno.Rollno)
		statement, _ = database.Prepare("INSERT INTO freezerecords (rollno, coins) VALUES (?, ?)")
		statement.Exec(freezerollno, coins)
	   }
}

func IfFrozen(rollno string) string {
	database, _ := sql.Open("sqlite3", "./blacklist.db")
    row, _ := database.Query("SELECT id, rollno, coins FROM freezerecords")
	var id int	
	var frozenrollno string
    var coins int
	var flag = 0

	   for row.Next() {
          row.Scan(&id, &frozenrollno, &coins)
		   if rollno == frozenrollno {
             flag = 1
		   }	  
		}
		if flag == 1 {
			return "true"
		} else {
			return "false"
		}  
}

func GiveRewards(response http.ResponseWriter, request *http.Request){
	Username := CheckjwtToken(response, request)
	response.Header().Set("Content-Type","application/json")
	  var userrewards UserRewards
	  json.NewDecoder(request.Body).Decode(&userrewards) 
	  
	  if Username != "false" && IsAdmin(Username) == "true" {
	   database, _ := sql.Open("sqlite3", "./data.db")
	//    fmt.Println(userrewards.RollNo, userrewards.InsertedCoins)
	//    stmtupdate, _ := database.Prepare("UPDATE userdata set coins=? WHERE rollno = ?")
    //       stmtupdate.Exec(100, coins.RollNo)		    
    intialrow, _ := database.Query("SELECT id, rollno, rewards FROM userdata")
	var id int
    var rollno string
    var intialrewards int
	var flag = 0

	   for intialrow.Next() {
          intialrow.Scan(&id, &rollno, &intialrewards)
		  fmt.Println(rollno)
		//  if coins.RollNo == rollno {	 
		// 	 fmt.Println("found")
        //      flag = 1
		//  }
		}

		if flag == 1{
			  stmtupdate, _ := database.Prepare("UPDATE userdata set rewards=? WHERE rollno = ?")
          stmtupdate.Exec(userrewards.Rewards, userrewards.Rollno)
		  response.Write([]byte(fmt.Sprintf("%s has %s rewards", userrewards.Rollno, strconv.Itoa(userrewards.Rewards))))
		} else {
			response.Write([]byte(fmt.Sprintf("NO DATA AVAILABLE")))
		}
	  }

}	


    // updatedrow, _ := database.Query("SELECT id, rollno, coins FROM userdata")
	// var id1 int
    // var rollno1 string
    // var coins1 int

	//    for updatedrow.Next() {
    //      updatedrow.Scan(&id1, &rollno1, &coins1)
	// 	 response.Write([]byte(fmt.Sprintf(strconv.Itoa(id1) + ": " + rollno1 + " " + strconv.Itoa(coins1))))
	// }


	func ShowUserCoins(rollno string) int {
		database, _ := sql.Open("sqlite3", "./data.db")
		row, _ := database.Query("SELECT id, rollno, coins FROM userdata")
		var requestedrollno string
		var coins int
		var id int
		var usercoins int
		   for row.Next() {
			  row.Scan(&id, &requestedrollno, &coins)
			 if rollno == requestedrollno {
				 usercoins = coins
			 }
			}
			return usercoins
}