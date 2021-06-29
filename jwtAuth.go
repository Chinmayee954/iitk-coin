package main

import (
	// "encoding/json"
	 "fmt"
	 "net/http"
	//   "time"
	//    "strconv"
	//   "context"
	//   "log"

	 "github.com/dgrijalva/jwt-go"
     _ "github.com/mattn/go-sqlite3"
	//  "database/sql"
	//    "golang.org/x/crypto/bcrypt"
)


func CheckjwtToken(w http.ResponseWriter, r *http.Request){
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