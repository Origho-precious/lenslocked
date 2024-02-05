package controllers

import (
	"github/Origho-precious/lenslocked/views"
	"net/http"
)

type User struct {
	Template struct {
		New views.Template
	}
}

func (u User) Create(w http.ResponseWriter, r *http.Request) {
	u.Template.New.Execute(w, nil)
}
