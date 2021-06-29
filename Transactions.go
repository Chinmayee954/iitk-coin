
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
	Coins string `json:"coins"`
 }

 type Claims struct {
	Username string `json:"rollno"`
	jwt.StandardClaims
}





func DoTransaction(response http.ResponseWriter, request *http.Request){

		cookie, err := request.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			response.WriteHeader(http.StatusUnauthorized)
			return
		}
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenStr := cookie.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			response.WriteHeader(http.StatusUnauthorized)
			return
		}
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	if !tkn.Valid {
		response.WriteHeader(http.StatusUnauthorized)
		return
	}

	response.Write([]byte(fmt.Sprintf("Hello, %s", claims.Username)))


		response.Header().Set("Content-Type","application/json")
	  var transaction Transaction
	  json.NewDecoder(request.Body).Decode(&transaction)
	   database, err := sql.Open("sqlite3", "./data.db")	
	   
	   if err != nil {
		log.Fatal(err)
	    }

	      coins, _ := strconv.Atoi(transaction.Coins)


		  ctx := context.Background()
	    tx, err := database.BeginTx(ctx, nil)

	if err != nil {
		log.Fatal(err)
	}

       var rollno1 string
	    var rollno2 string
	   var id1 string
	   var id2 string
	   var coins1 string
	   var coins2 string

	    var tcoins1 int
	    var tcoins2 int

	     row1, _ := tx.Query("SELECT id, rollno, coins FROM people")

		  for row1.Next() {
           err =  row1.Scan(&id1, &rollno1, &coins1) 
		    if err != nil {
		    tx.Rollback()
		    return
     	    }
		 if claims.Username == rollno1 {
			 fmt.Println(rollno1)
			   tcoins1,_ = strconv.Atoi(coins1)
			
		 }
    }

		 row2, _ := tx.Query("SELECT id, rollno, coins FROM people")

		   for row2.Next() {
        err =  row2.Scan(&id2, &rollno2, &coins2) 

		 if err != nil {
		    tx.Rollback()
		    return
     	    }
		 if transaction.RollNo == rollno2 {
			 fmt. Println(rollno1)
			   tcoins2,_ = strconv.Atoi(coins2)
			
		 }
    }

	        // fmt.Println(coins1)
			// fmt.Println(coins2)
			tax := (coins*2)/100
		
		     tcoins1 = tcoins1-coins - tax
		     tcoins2 = tcoins2+coins - tax
			// fmt.Println(tcoins1)
			// fmt.Println(tcoins2)

			

			if tcoins1 < 0 {
				 response.Write([]byte(fmt.Sprintf("%s does not have enough coins", claims.Username )))
				  tx.Rollback()
				return
			} else {
			// stmtupdate1, _ := database.Prepare("UPDATE people set coins=? WHERE rollno=?")
			_,err = tx.ExecContext(ctx, "UPDATE people set coins=? WHERE rollno=?", strconv.Itoa(tcoins1), claims.Username)
			    if err != nil {
		         tx.Rollback()
		         return
	                }
            // stmtupdate1.Exec(strconv.Itoa(tcoins1), transaction.RollNo1)

			// stmtupdate2, _ := database.Prepare("update people set coins=? where rollno=?")
			_,err = tx.ExecContext(ctx, "UPDATE people set coins=? WHERE rollno=?", strconv.Itoa(tcoins2), transaction.RollNo)
			      if err != nil {
		         tx.Rollback()
		         return
	                }
            // stmtupdate2.Exec(strconv.Itoa(tcoins2), transaction.RollNo2)
			}


		updatedrow, _ := database.Query("SELECT id, rollno, coins FROM people")
	
	var id int
    var rollno string
    var updatedcoins string

	   for updatedrow.Next() {
         updatedrow.Scan(&id, &rollno, &updatedcoins)
	 	 fmt.Println(strconv.Itoa(id) + ": " + rollno + " " + updatedcoins)

	 	}

		 err = tx.Commit()
		//  recordtime := time.Now();
		fmt.Println(transaction.RollNo)
			fmt.Println(claims.Username)

		  AddHistory(transaction.RollNo, claims.Username, transaction.Coins)


		//  History()
	    if err != nil {
		log.Fatal(err)
	}

		 response.Write([]byte(fmt.Sprintf("transaction success")))
		//  response.Write([]byte(fmt.Sprintf("%s has %s coins", transaction.RollNo2, strconv.Itoa(tcoins2))))

	}

