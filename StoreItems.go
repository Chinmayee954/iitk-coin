package main




import (
	//  "encoding/json"
	 "fmt"
	 "net/http"
	//   "time"
	   "strconv"
	   "encoding/json"
	//   "context"
	   "log"

	//  "github.com/dgrijalva/jwt-go"
     _ "github.com/mattn/go-sqlite3"
	 "database/sql"
	//    "golang.org/x/crypto/bcrypt"
)

type StoreItem struct {
	Item string `json:"item"`
	Price int 	`json:"price"`
}

type Approval struct {
	Id int `json:"id"`
}

func AddStoreItem(response http.ResponseWriter, request *http.Request) {
	Username := CheckjwtToken(response, request)
    response.Header().Set("Content-Type","application/json")
	var storeitem StoreItem
	json.NewDecoder(request.Body).Decode(&storeitem)
	database, err := sql.Open("sqlite3", "./store.db")
	if err != nil {
	log.Fatal(err)
	}	
	if Username != "false" && IsAdmin(Username) == "true" {
    statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS storeitems (id INTEGER PRIMARY KEY, item TEXT, price INTEGER)")
    statement.Exec()
    stmt, _ := database.Prepare("INSERT INTO storeitems (item, price) VALUES (?, ?)")
    stmt.Exec(storeitem.Item, storeitem.Price)

	row, _ := database.Query("SELECT id, item, price FROM storeitems")
    var id int
    var item string
    var price int
	for row.Next() {
         row.Scan(&id, &item, &price)
	 	 fmt.Println(strconv.Itoa(id) + ": " + item + " " + strconv.Itoa(price))
	}
 }
}

func ApproveRedeem(response http.ResponseWriter, request *http.Request){
	Username := CheckjwtToken(response, request)
	response.Header().Set("Content-Type","application/json")
	var approval Approval
	json.NewDecoder(request.Body).Decode(&approval)
	database, err := sql.Open("sqlite3", "./store.db")
	Checkerr(err)
	database1, err := sql.Open("sqlite3", "./data.db")
	Checkerr(err)
    if Username != "false" && IsAdmin(Username) == "true" {
    fmt.Println(IsAdmin(Username))
	row, _ := database.Query("SELECT id, rollno, coins, reqitem, price, status FROM redeemreqs")
	
    var id int
	var rollno string
	var coins int
	var reqitem string
    var price int
    var status string
	var returncoins int
	var usercoins int
	for row.Next() {
         row.Scan(&id, &rollno, &coins, &reqitem, &price, &status)
		 fmt.Println(id, approval.Id)
			   if id == approval.Id {
				   usercoins = ShowUserCoins(rollno)
				fmt.Println(approval.Id)
				returncoins = usercoins-price
				fmt.Println(returncoins)
                stmt1, _ := database1.Prepare("UPDATE userdata set coins=? WHERE rollno = ?")
				stmt1.Exec(returncoins, rollno)
			   }           
	}
			stmt2, err := database.Prepare("DELETE from redeemreqs WHERE id=?")
	 Checkerr(err)
	stmt2.Exec(approval.Id)
	response.Write([]byte(fmt.Sprintf("Request has been approved!")))
}

row, _ := database.Query("SELECT id, rollno, coins, reqitem, price, status FROM redeemreqs")
    var id1 int
	var rollno1 string
	var coins1 int
	var reqitem1 string
    var price1 int
    var status1 string
	for row.Next() {
         row.Scan(&id1, &rollno1, &coins1, &reqitem1, &price1, &status1)
	 	 fmt.Println(strconv.Itoa(id1) + ": " + rollno1 + " " + strconv.Itoa(coins1) + 
		 " " + strconv.Itoa(price1) + reqitem1 + " " + status1)
	}

row1, _ := database1.Query("SELECT id, rollno, coins FROM userdata")
 var id2 int
	var rollno2 string
	var coins2 int	
for row1.Next() {
         row1.Scan(&id2, &rollno2, &coins2)
	 	 fmt.Println(strconv.Itoa(id2) + ": " + rollno2 + " " + strconv.Itoa(coins2))
	}
}


