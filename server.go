package main

import (
	"os"
	"fmt"
	"net/http"
    "database/sql"
	"log"

    "github.com/go-sql-driver/mysql"
)

type Person struct {
	fullName string
	phone string
	email string
	gender string
	birthDate string
	languages []string
	bio string
	contract bool
}

var db *sql.DB

func indexHandler(output http.ResponseWriter, request *http.Request){
	http.ServeFile(output, request, "index.html")
}

func parseFormRequest(request *http.Request) Person{
	fullName := request.FormValue("fullName") // first only fiels
	phone := request.FormValue("phone")
	email := request.FormValue("email")
	gender := request.FormValue("gender")
	birthDate := request.FormValue("birthDate")
	languages := request.Form["languages"] // all fields!
	bio := request.FormValue("bio")

	// Maybe unnescessary check, cause contract is allways must be on :)
	var contract bool
	if request.FormValue("contract") == "on"{
		contract = true
	} else {
		contract = false
	}

	person := Person{
		fullName, phone,  email, gender,
		birthDate, languages, bio, contract,
	}

	return person
}

func dataBseConnection(){
	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DBUSER")
	cfg.Passwd = os.Getenv("DBPASS")
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "web3_users"

	db, _ = sql.Open("mysql", cfg.FormatDSN())

    pingErr := db.Ping()
    if pingErr != nil {
        log.Fatal(pingErr)
    }

	fmt.Println("Bd connected!")
}

func savePerson(output http.ResponseWriter, request *http.Request){
	if request.Method != "Post" {
		http.Redirect(output, request, "/", http.StatusSeeOther)
	}
	person := parseFormRequest(request)
	
	// here must be function to save perosn in mysql XDD

	fmt.Println(person)
}

func main(){
	dataBseConnection()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/submit", savePerson)
	http.ListenAndServe(":8080", nil)
}
