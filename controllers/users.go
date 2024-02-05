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
	var data struct {
		Email string
	}

	data.Email = r.FormValue("email")

	fmt.Println(data)

	u.Template.New.Execute(w, data)
}

func (u User) Create(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	fmt.Println(email, password)
	fmt.Fprint(w, "Response")
}
