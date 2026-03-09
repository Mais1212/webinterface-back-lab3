package main

import (
	"fmt"
	"net/http"
)

func myHandler(output http.ResponseWriter, r *http.Request){
	title := r.URL.Path[len("/"):]
	fmt.Fprintf(output, "<h6>%s</h6>", title)
}

func main(){
	http.HandleFunc("/", myHandler)
	http.ListenAndServe(":8080", nil)
}
