package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "<h1>Hello, World!</h1>")
}

func main(){
	http.HandleFunc("/", handlerFunc)
	fmt.Println("Starting the server on PORT 3001...")
	http.ListenAndServe(":3001", nil)
}