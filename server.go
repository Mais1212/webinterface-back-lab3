package main

import (
	"fmt"
	"net/http"
)

type Person struct {
	fullName string
	phone string
	email string
	gender int8
	birthDate string
	languages []string
	bio string
	contract bool
}

func indexHandler(output http.ResponseWriter, request *http.Request){
	http.ServeFile(output, request, "index.html")
}

func parseForm(output http.ResponseWriter, request *http.Request){
	if request.Method != "Post" {
		http.Redirect(output, request, "/", http.StatusSeeOther)
	}

		fullName := request.FormValue("fullName") // first only

		phone := request.FormValue("phone")
		email := request.FormValue("email")
		gender := request.FormValue("gender")
		birthDate := request.FormValue("birthDate")
		languages := request.Form["languages"] // all!
		bio := request.FormValue("bio")
		contract := request.FormValue("contract")

		fmt.Println(
			fullName,"\n",phone, "\n", email, "\n", gender, "\n", birthDate, "\n", 
			languages, "\n",  bio, "\n", contract)
}

func main(){
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/submit", parseForm)
	http.ListenAndServe(":8080", nil)
}
