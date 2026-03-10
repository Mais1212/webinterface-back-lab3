package main

import (
	"fmt"
	"net/http"
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

func savePerson(output http.ResponseWriter, request *http.Request){
	if request.Method != "Post" {
		http.Redirect(output, request, "/", http.StatusSeeOther)
	}
	person := parseFormRequest(request)
	
	// here must be function to save perosn in mysql XDD

	fmt.Println(person)
}

func main(){
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/submit", savePerson)
	http.ListenAndServe(":8080", nil)
}
