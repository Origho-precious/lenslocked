package controllers

import (
	"fmt"
	"net/http"
)

type UserTemplate struct {
	New Template
}

type User struct {
	Template UserTemplate
}

func (u User) New(w http.ResponseWriter, r *http.Request) {
	u.Template.New.Execute(w, nil)
}

func (u User) Create(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	fmt.Println(email, password)
	fmt.Fprint(w, "Response")
}
