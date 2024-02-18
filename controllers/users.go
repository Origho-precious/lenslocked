package controllers

import (
	"fmt"
	"github/Origho-precious/lenslocked/models"
	"net/http"
)

type UserTemplate struct {
	New Template
}

type Users struct {
	Template    UserTemplate
	UserService *models.UserService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}

	data.Email = r.FormValue("email")

	u.Template.New.Execute(w, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := u.UserService.Create(email, password)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	fmt.Println(email, password)
	fmt.Fprintf(w, "User created: %+v", user)
}
