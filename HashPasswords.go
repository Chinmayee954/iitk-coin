package main

import (
	// "encoding/json"
	 "fmt"
	//  "net/http"
	//   "time"
	//    "strconv"
	//   "context"
	//   "log"

	// "github.com/dgrijalva/jwt-go"
     _ "github.com/mattn/go-sqlite3"
	//  "database/sql"
	    "golang.org/x/crypto/bcrypt"
)

func getHash(pwd []byte) string {        
    hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)          
    if err != nil {
      fmt.Println(err)
    }
    return string(hash)
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