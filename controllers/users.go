package controllers

import "net/http"

type UserTemplate struct {
	New Template
}

type User struct {
	Template UserTemplate
}

func (u User) New(w http.ResponseWriter, r *http.Request) {
	u.Template.New.Execute(w, nil)
}

func (u User) ReturnJson(w http.ResponseWriter, r *http.Request) {
	u.Template.New.Execute(w, nil)
}
