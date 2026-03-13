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
	languages []int
	bio string
	contract int
}

type Language struct {
	id int
	name string
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

	dbLanguages := getLanguages(db)
	var languagesIds []int

	// неоптимизировано + ужасно, сейчас не хочу править и не знаю как
	for _, value := range languages{
		for _, language := range dbLanguages{
			if language.name == value {
				languagesIds = append(languagesIds, language.id)
				break
			}
		}
	}

	// Maybe unnescessary check, cause contract is allways must be on :)
	var contract int
	if request.FormValue("contract") == "on"{
		contract = 1
	} else {
		contract = 0
	}

	person := Person{
		fullName, phone,  email, gender,
		birthDate, languagesIds, bio, contract,
	}

	return person
}


func dataBseConnection(){
	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DBUSER")
	cfg.Passwd = os.Getenv("DBPASS")
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "web3"

	db, _ = sql.Open("mysql", cfg.FormatDSN())

    pingErr := db.Ping()
    if pingErr != nil {
        log.Fatal(pingErr)
    }

	fmt.Println("Bd connected!")
}


func getLanguages(db *sql.DB) []Language {
    rows, _ := db.Query("SELECT id, name FROM programming_languages")

	var languages []Language

    for rows.Next() {
		var language Language
		rows.Scan(&language.id, &language.name)
		languages = append(languages, language)
    }

	return languages
}


func savePerson(output http.ResponseWriter, request *http.Request){
	if request.Method != "Post" {
		http.Redirect(output, request, "/", http.StatusSeeOther)
	}
	person := parseFormRequest(request)

	// Querry to add new person
	result, _ := db.Exec(
		"INSERT INTO person (full_name, phone, email, birth_date, gender, biography) VALUES (?,?,?,?,?,?);",
		person.fullName, person.phone, person.email, person.birthDate,
		person.gender, person.bio,
	)

	personID, err := result.LastInsertId()
	if err != nil{
        log.Fatal(err)
	}


	// Querry to add favorite languages for person!
	for _, langID := range person.languages {
		db.Exec(`
			INSERT INTO person_language (person_id, language_id) 
			VALUES (?, ?)`,
			personID, langID,
		)
	}
}


func main(){
	dataBseConnection()
	

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/submit", savePerson)
	http.ListenAndServe(":8080", nil)
}
