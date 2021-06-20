package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	  "strconv"
	  "context"
	  "log"

	"github.com/dgrijalva/jwt-go"
     _ "github.com/mattn/go-sqlite3"
	 "database/sql"
	   "golang.org/x/crypto/bcrypt"
)


var jwtKey = []byte("secret_key")

// var users = map[string]string{
// 	"user1": "password1",
// 	"user2": "password2",
// }

type Credentials struct {
	Username string `json:"rollno"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"rollno"`
	jwt.StandardClaims
}

type User struct{
    Username string `json:"rollno"`
    Password string `json:"password"`
}

type Coins struct{
    RollNo string `json:"rollno"`
    InsertedCoins string `json:"coins"`
}

type Transaction struct{
	RollNo1 string `json:"rollno1"`
	RollNo2 string `json:"rollno2"`
	Coins string `json:"coins"`
 }

 type RollNo struct {
	 RequestedRollNo string `json:"rollno"`
 }

func getHash(pwd []byte) string {        
    hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)          
    if err != nil {
      fmt.Println(err)
    }
    return string(hash)
}


// func HashPassword(password string) (string) {
//     bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
//     return string(bytes)
// }


// func IfUserExists(dataUsername string, signupUsername string) bool{
//     if(dataUsername == signupUsername){
// 		Println("user already exists")
// 		return false
// 	}

// }



func SignUp(response http.ResponseWriter, request *http.Request) {
      response.Header().Set("Content-Type","application/json")
	  var user User
	  json.NewDecoder(request.Body).Decode(&user)

    database, _ := sql.Open("sqlite3", "./signup_rollno.db")

	  user.Password = getHash([]byte(user.Password))

	    rows1, _ := database.Query("SELECT id, rollno, password FROM SignUpRoll")
	    var id1 int
    var rollno1 string
    var password1 string

	   for rows1.Next() {
         rows1.Scan(&id1, &rollno1, &password1)
		 if(user.Username == rollno1){
			 fmt.Println("user already exists")
			 return
		 }
        // fmt.Println(strconv.Itoa(id) + ": " + username + " " + password)
    }

	  
	//   fmt.Println(user.Password)
	  
	//   stmt := "SELECT id FROM signupUsers WHERE username = ?"
	//   row := database.QueryRow(stmt)

	//   var uID string 
	//   err := row.Scan(&uID)

	//   if err != sql.ErrNoRows {
	// 	  fmt.Println("user already exists", err);
	// 	  return;
	//   }

	var signupuser string = user.Username
	var signuppassword string = user.Password

	// fmt.Println(signupuser, signuppassword)

	

	    statement, _ := database.Prepare("INSERT INTO SignUpRoll (rollno, password) VALUES (?, ?)")
	    statement.Exec(signupuser, signuppassword)

	    rows, _ := database.Query("SELECT id, rollno, password FROM SignUpRoll")
	    var id int
    var rollno string
    var password string
    for rows.Next() {
        rows.Scan(&id, &rollno, &password)
        fmt.Println(strconv.Itoa(id) + ": " + rollno + " " + password)
    }

	//   response.Write([]byte(fmt.Sprintf("congrats, %s", user.Username)))
	  

}


func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		//log.Println(err)
		return false
	}

	return true
}


func Login(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials
     database, _ := sql.Open("sqlite3", "./signup_rollno.db")
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	
	    rows, _ := database.Query("SELECT id, rollno, password FROM SignUpRoll")
	    var id int
    var rollno string
    var password string

	m := make(map[string]string)

	 for rows.Next() {
        rows.Scan(&id, &rollno, &password)
        // fmt.Println(strconv.Itoa(id) + ": " + rollno + " " + password)
        //   expectedPassword, ok := password

	    //  if !ok || expectedPassword != credentials.Password {
		// w.WriteHeader(http.StatusUnauthorized)
		// return
         m[rollno] = password
	}
    
	
	 expectedPassword, ok := m[credentials.Username]

	if !ok || !comparePasswords(expectedPassword, []byte(credentials.Password) ) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(time.Minute * 5)

	claims := &Claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w,
		&http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})

}

func Home(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
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
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Write([]byte(fmt.Sprintf("Hello, %s", claims.Username)))

}

func Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
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
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	expirationTime := time.Now().Add(time.Minute * 5)

	claims.ExpiresAt = expirationTime.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w,
		&http.Cookie{
			Name:    "refresh_token",
			Value:   tokenString,
			Expires: expirationTime,
		})

}


func AwardCoins(response http.ResponseWriter, request *http.Request)  {
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
        
	 	 fmt.Println(strconv.Itoa(id) + ": " + rollno + " " + updatedcoins)

	 	}

	}

	func DoTransaction(response http.ResponseWriter, request *http.Request){
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


	//    var rollno1 string  = transaction.RollNo1
	//    var rollno2 string  = transaction.RollNo2

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
		 if transaction.RollNo1 == rollno1 {
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
		 if transaction.RollNo2 == rollno2 {
			 fmt. Println(rollno1)
			   tcoins2,_ = strconv.Atoi(coins2)
			
		 }
    }

	        // fmt.Println(coins1)
			// fmt.Println(coins2)
		
		     tcoins1 = tcoins1-coins
		     tcoins2 = tcoins2+coins

			// fmt.Println(tcoins1)
			// fmt.Println(tcoins2)

			fmt.Println(transaction.RollNo1)

			if tcoins1 < 0 {
				 response.Write([]byte(fmt.Sprintf("%s does not have enough coins", transaction.RollNo1 )))
				  tx.Rollback()
				return
			} else {
			// stmtupdate1, _ := database.Prepare("UPDATE people set coins=? WHERE rollno=?")
			_,err = tx.ExecContext(ctx, "UPDATE people set coins=? WHERE rollno=?", strconv.Itoa(tcoins1), transaction.RollNo1)
			    if err != nil {
		         tx.Rollback()
		         return
	                }
            // stmtupdate1.Exec(strconv.Itoa(tcoins1), transaction.RollNo1)

			// stmtupdate2, _ := database.Prepare("update people set coins=? where rollno=?")
			_,err = tx.ExecContext(ctx, "UPDATE people set coins=? WHERE rollno=?", strconv.Itoa(tcoins2), transaction.RollNo2)
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
	if err != nil {
		log.Fatal(err)
	}

		 response.Write([]byte(fmt.Sprintf("transaction success")))
		//  response.Write([]byte(fmt.Sprintf("%s has %s coins", transaction.RollNo2, strconv.Itoa(tcoins2))))

	}

		func Showcoins(response http.ResponseWriter, request *http.Request){
                  response.Header().Set("Content-Type","application/json")
	  var rollno RollNo
	  json.NewDecoder(request.Body).Decode(&rollno)
	   database, _ := sql.Open("sqlite3", "./data.db")
          
         row, _ := database.Query("SELECT id, rollno, coins FROM people")

	var requestedrollno string
    var coins string
    var id int

	   for row.Next() {
          row.Scan(&id, &requestedrollno, &coins)
		 if rollno.RequestedRollNo == requestedrollno {
			 response.Write([]byte(fmt.Sprintf("coins : %s", coins)))
			 return
		 }
		}

	}