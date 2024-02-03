package controllers

import (
	"github/Origho-precious/lenslocked/views"
	"net/http"
)

func Statichandler(tpl views.Template, data interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, data)
	}
}
