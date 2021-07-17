package main

import (
"encoding/json"
"fmt"
"net/http"
"time"
"strconv"
// "context"
 "log"

"github.com/dgrijalva/jwt-go"
_ "github.com/mattn/go-sqlite3"
"database/sql"
// "golang.org/x/crypto/bcrypt"
)


var jwtKey = []byte("secret_key")

type Credentials struct {
Username string `json:"rollno"`
Password string `json:"password"`
}



type User struct{
Username string `json:"rollno"`
Password string `json:"password"`
}

type RollNo struct {
RequestedRollNo string `json:"rollno"`
}



func SignUp(response http.ResponseWriter, request *http.Request) {
response.Header().Set("Content-Type","application/json")
var user User
json.NewDecoder(request.Body).Decode(&user)

database, err := sql.Open("sqlite3", "./data.db")
Checkerr(err)
statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS userdata (id INTEGER PRIMARY KEY, rollno TEXT, password TEXT, coins INTEGER, rewards INTEGER)")
Checkerr(err)
statement.Exec()


// statement1, _ := database.Prepare("INSERT INTO userdata (rollno, password) VALUES (?, ?)")
// statement1.Exec("190100", "chin@1234")



user.Password = getHash([]byte(user.Password))

rows, err := database.Query("SELECT id, rollno, password FROM userdata")
Checkerr(err)
var id int
var rollno string
var password string
var flag = 1

fmt.Println(user.Username)



for rows.Next() {
rows.Scan(&id, &rollno, &password)
if user.Username == rollno {
flag = 0
break;
}
}

if flag == 1 {
var signupuser string = user.Username
var signuppassword string = user.Password

statement, _ := database.Prepare("INSERT INTO userdata (rollno, password, coins, rewards) VALUES (?, ?, ?, ?)")
statement.Exec(signupuser, signuppassword, 0, 0)
} else {
response.Write([]byte(fmt.Sprintf("You have already signed up")))
}


rows1, _ := database.Query("SELECT id, rollno, password, coins, rewards FROM userdata")
var id1 int
var rollno1 string
var password1 string
var coins1 int
var rewards1 int
for rows1.Next() {
rows1.Scan(&id1, &rollno1, &password1, &coins1, &rewards1)
fmt.Println(strconv.Itoa(id1) + ": " + rollno1 + " " + password1 + " " + strconv.Itoa(coins1) + " " +
strconv.Itoa(rewards1))
}
}





func Login(w http.ResponseWriter, r *http.Request) {
var credentials Credentials
database, _ := sql.Open("sqlite3", "./data.db")
err := json.NewDecoder(r.Body).Decode(&credentials)
if err != nil {
w.WriteHeader(http.StatusBadRequest)
return
}


rows, _ := database.Query("SELECT id, rollno, password FROM userdata")
var id int
var rollno string
var password string

m := make(map[string]string)

for rows.Next() {
rows.Scan(&id, &rollno, &password)
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
Name: "token",
Value: tokenString,
Expires: expirationTime,
})

}

func Home(w http.ResponseWriter, r *http.Request) {
CheckjwtToken(w,r)
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
// w.WriteHeader(http.StatusBadRequest)
// return
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
Name: "refresh_token",
Value: tokenString,
Expires: expirationTime,
})

}













func Showcoins(response http.ResponseWriter, request *http.Request){
response.Header().Set("Content-Type","application/json")
var rollno RollNo
json.NewDecoder(request.Body).Decode(&rollno)
database, _ := sql.Open("sqlite3", "./data.db")

row, _ := database.Query("SELECT id, rollno, coins FROM userdata")

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


func Checkerr(err error){
	if err != nil {
		log.Fatal(err)
	}
}