package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	  "strconv"

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