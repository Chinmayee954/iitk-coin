
package main

import (
	 "encoding/json"
	 "fmt"
	 "net/http"
	//"time"
	   "strconv"
	   "context"
	   "log"

	 "github.com/dgrijalva/jwt-go"
     _ "github.com/mattn/go-sqlite3"
	 "database/sql"
	//    "golang.org/x/crypto/bcrypt"
)




type Transaction struct{
	RollNo string `json:"rollno"`
	Coins int `json:"coins"`
 }

 type Claims struct {
	Username string `json:"rollno"`
	jwt.StandardClaims
}

type RedeemReq struct {
	Item string `json:"item"`
}




func DoTransaction(response http.ResponseWriter, request *http.Request){
    Username := CheckjwtToken(response, request)
    response.Header().Set("Content-Type","application/json")
	var transaction Transaction
	json.NewDecoder(request.Body).Decode(&transaction)
	database, err := sql.Open("sqlite3", "./data.db")	   
	if err != nil {
	log.Fatal(err)
	}
	coins := transaction.Coins
	ctx := context.Background()
	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	if Username != "false" {
    if IfFrozen(Username) == "true" {
		tx.Rollback()
		return
	}
	if IfFrozen(transaction.RollNo) == "true" {
		tx.Rollback()
		return
	}
       var rollno1 string
	   var rollno2 string
	   var id1 int
	   var id2 int
	   var coins1 int
	   var coins2 int
	   var tcoins1 int
	   var tcoins2 int
    row1, _ := tx.Query("SELECT id, rollno, coins FROM userdata")
	       for row1.Next() {
           err =  row1.Scan(&id1, &rollno1, &coins1) 
		    if err != nil {
		    tx.Rollback()
		    return
     	    }
            if Username == rollno1 {
			 fmt.Println(rollno1)
			   tcoins1 = coins1
		 }
    }
	row2, _ := tx.Query("SELECT id, rollno, coins FROM userdata")
		for row2.Next() {
         err =  row2.Scan(&id2, &rollno2, &coins2) 
		 if err != nil {
		    tx.Rollback()
		    return
     	    }
		 if transaction.RollNo == rollno2 {
			 fmt.Println(rollno1)
			   tcoins2 = coins2
		 }
    }
	       	    tax := (coins*2)/100
		    tcoins1 = tcoins1-coins - tax
		    tcoins2 = tcoins2+coins - tax
			if tcoins1 < 0 {
				 response.Write([]byte(fmt.Sprintf("%s does not have enough coins", Username )))
				  tx.Rollback()
				return
			} else {
			_,err = tx.ExecContext(ctx, "UPDATE userdata set coins=? WHERE rollno=?", tcoins1, Username)
			    if err != nil {
		         tx.Rollback()
		         return
	                }
			_,err = tx.ExecContext(ctx, "UPDATE userdata set coins=? WHERE rollno=?", tcoins2, transaction.RollNo)
			      if err != nil {
		         tx.Rollback()
		         return
				  }
			}
		} else {
			response.WriteHeader(http.StatusUnauthorized)
		}

	updatedrow, _ := database.Query("SELECT id, rollno, coins FROM userdata")
    var id int
    var rollno string
    var updatedcoins int
	   for updatedrow.Next() {
         updatedrow.Scan(&id, &rollno, &updatedcoins)
	 	 fmt.Println(strconv.Itoa(id) + ": " + rollno + " " + strconv.Itoa(updatedcoins))
	 	}
		 err = tx.Commit()
		//  recordtime := time.Now();
		fmt.Println(transaction.RollNo)
			fmt.Println(Username)
		  AddHistory(Username, transaction.RollNo, transaction.Coins)
		//  History()
	    if err != nil {
		log.Fatal(err)
	}
         	 response.Write([]byte(fmt.Sprintf("transaction success")))
	
	}




func Redeem(response http.ResponseWriter, request *http.Request) {
		Username := CheckjwtToken(response, request)
		response.Header().Set("Content-Type","application/json")
		  var redeemreq RedeemReq
		  json.NewDecoder(request.Body).Decode(&redeemreq)    
		   database, _ := sql.Open("sqlite3", "./store.db")

		if Username != "false" {
        
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS redeemreqs (id INTEGER PRIMARY KEY, rollno TEXT, coins INTEGER, reqitem TEXT, price INTEGER, status TEXT)")
    statement.Exec()
    coins := ShowUserCoins(Username)
	price := ShowItemPrice(redeemreq.Item)
	fmt.Println(price)
	fmt.Println(coins)
	if price < coins {
      stmt, _ := database.Prepare("INSERT INTO redeemreqs (rollno, coins, reqitem, price, status) VALUES (?, ?, ?, ?, ?)")
      stmt.Exec(Username, coins, redeemreq.Item, price, "requested")
	} else {
		response.Write([]byte(fmt.Sprintf("You do not have enough coins")))
	}
    row, _ := database.Query("SELECT rollno, reqitem, status FROM redeemreqs")
    // var id int
    var rollno string
	var reqitem string
    // var coins1 int
	// var price1 int
	var status string
	   for row.Next() {
         row.Scan(&rollno, &reqitem, &status)
	 	 fmt.Println(rollno + " " + reqitem + " " + status)
	 	}

	}
}

	func ShowItemPrice(reqitem string) int {
		 database, _ := sql.Open("sqlite3", "./store.db")
		 mp       := make(map[string]int)
		 var id int
		 var item string
		 var price int
     row,_ := database.Query("SELECT id, item, price FROM storeitems")
     for row.Next() {
     row.Scan(&id, &item, &price)
     mp[item] = price
     }

	 return mp[reqitem]
}