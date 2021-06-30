
package main

import (
	 "encoding/json"
	 "fmt"
	 "net/http"
	//   "time"
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

       var rollno1 string
	    var rollno2 string
	   var id1 int
	   var id2 int
	   var coins1 int
	   var coins2 int


	     row1, _ := tx.Query("SELECT id, rollno, coins FROM userdata")

		  for row1.Next() {
           err =  row1.Scan(&id1, &rollno1, &coins1) 
		    if err != nil {
		    tx.Rollback()
		    return
     	    }
		 if Username == rollno1 {
			 fmt.Println(rollno1)
			   coins1 = coins1
			
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
			 fmt. Println(rollno1)
			   coins2 = coins2
			
		 }
    }


			tax := (coins*2)/100
		
		     coins1 = coins1-coins - tax
		     coins2 = coins2+coins - tax
			
			

			if coins1 < 0 {
				 response.Write([]byte(fmt.Sprintf("%s does not have enough coins", Username )))
				  tx.Rollback()
				return
			} else {
			// stmtupdate1, _ := database.Prepare("UPDATE people set coins=? WHERE rollno=?")
			_,err = tx.ExecContext(ctx, "UPDATE userdata set coins=? WHERE rollno=?", coins1, Username)
			    if err != nil {
		         tx.Rollback()
		         return
	                }
            // stmtupdate1.Exec(strconv.Itoa(tcoins1), transaction.RollNo1)

			// stmtupdate2, _ := database.Prepare("update people set coins=? where rollno=?")
			_,err = tx.ExecContext(ctx, "UPDATE userdata set coins=? WHERE rollno=?", coins2, transaction.RollNo)
			      if err != nil {
		         tx.Rollback()
		         return
	                }
            // stmtupdate2.Exec(strconv.Itoa(tcoins2), transaction.RollNo2)
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
		//  response.Write([]byte(fmt.Sprintf("%s has %s coins", transaction.RollNo2, strconv.Itoa(tcoins2))))

	}

