package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"html/template"
	"net/url"

	"github.com/go-sql-driver/mysql"
)

type Person struct {
	fullName  string
	phone     string
	email     string
	gender    string
	birthDate string
	languages []int
	bio       string
	contract  int
}

type Language struct {
	id   int
	name string
}

var db *sql.DB

func validatePerson(p Person) bool {
	// fullName: required, max 255 chars, кириллица/латиница + пробелы/дефисы
	if len(p.fullName) == 0 || len(p.fullName) > 255 {
		return false
	}
	nameRe := regexp.MustCompile(`^[ЁёА-Яа-яA-Za-z\s\-]+$`)
	if !nameRe.MatchString(p.fullName) {
		return false
	}

	// phone: required (HTML *), российский формат +7(999)123-45-67 или 89991234567
	if p.phone == "" || len(p.phone) < 10 {
		return false
	}
	phoneRe := regexp.MustCompile(`^(?:\+?7|8|7)?[\s\-\(\)]?[0-9]{3}[\s\-\)]?[0-9]{3}[\s\-]?[0-9]{2}[\s\-]?[0-9]{2}$`)
	if !phoneRe.MatchString(p.phone) {
		return false
	}

	// email: required (HTML *), стандарт
	if p.email == "" {
		return false
	}
	emailRe := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRe.MatchString(p.email) {
		return false
	}

	// birthDate: required (HTML *), YYYY-MM-DD от input type="date"
	if p.birthDate == "" {
		return false
	}
	dateRe := regexp.MustCompile(`^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01])$`)
	if !dateRe.MatchString(p.birthDate) {
		return false
	}

	// gender: male/female (из radio-buttons HTML)
	if p.gender != "male" && p.gender != "female" {
		return false
	}

	// languages: любой набор (пустой OK)
	for _, id := range p.languages {
		if id <= 0 {
			return false
		}
	}

	// contract: required checkbox (1 или 0)
	if p.contract != 0 && p.contract != 1 {
		return false
	}

	// bio: optional textarea
	return true
}


func indexHandler(output http.ResponseWriter, request *http.Request) {
    data := struct {
        Error string
    }{}
    errStr := request.URL.Query().Get("error")
    if errStr != "" {
        data.Error = errStr
    }

    tmpl := template.Must(template.ParseFiles("static/index.html"))
    tmpl.Execute(output, data)
}

func parseFormRequest(request *http.Request) Person {
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
	for _, value := range languages {
		for _, language := range dbLanguages {
			if language.name == value {
				languagesIds = append(languagesIds, language.id)
				break
			}
		}
	}

	// Maybe unnescessary check, cause contract is allways must be on :)
	var contract int
	if request.FormValue("contract") == "on" {
		contract = 1
	} else {
		contract = 0
	}

	person := Person{
		fullName, phone, email, gender,
		birthDate, languagesIds, bio, contract,
	}

	return person
}

func dataBseConnection() {
	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DBUSER")
	cfg.Passwd = os.Getenv("DBPASS")
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "web3"

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Bd connected!")
}

func getLanguages(db *sql.DB) []Language {

	rows, err := db.Query("SELECT id, name FROM programming_languages")

	if err != nil {
		log.Printf("Error querying languages: %v", err)
		return nil // или return []Language{}
	}
	defer rows.Close()

	var languages []Language

	for rows.Next() {
		var language Language
		rows.Scan(&language.id, &language.name)
		languages = append(languages, language)
	}

	return languages
}

func savePerson(output http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		fmt.Println(" err post!!")
		http.Redirect(output, request, "/", http.StatusSeeOther)
		return
	}

	person := parseFormRequest(request)

	if !validatePerson(person) {
			v := url.Values{}
			v.Add("error", "Некорректные данные!")
			http.Redirect(output, request, "/?"+v.Encode(), http.StatusSeeOther)
			fmt.Println(" errrorrr!!")
			return
	}

	// Querry to add new person
	result, err := db.Exec(
		"INSERT INTO person (full_name, phone, email, birth_date, gender, biography) VALUES (?,?,?,?,?,?);",
		person.fullName, person.phone, person.email, person.birthDate,
		person.gender, person.bio,
	)
	if err != nil {
		log.Printf("Error inserting person: %v", err)
		return
	}

	personID, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	// Querry to add favorite languages for person!
	for _, langID := range person.languages {
		_, err := db.Exec(`
			INSERT INTO person_language (person_id, language_id) 
			VALUES (?, ?)`,
			personID, langID,
		)

		if err != nil {
			log.Printf("Error insterting language %d for person %d: %v", langID,
				personID, err,
			)
		}
	}

	fmt.Println("ready!!!")
}

func main() {
	dataBseConnection()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/submit", savePerson)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.ListenAndServe(":8080", nil)
}
