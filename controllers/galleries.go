package controllers

import (
	"github/Origho-precious/lenslocked/models"
	"net/http"
)

type Templates struct {
	New Template
}

type Galleries struct {
	Templates      Templates
	GalleryService *models.GalleryService
}

func (g Galleries) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
	}

	data.Title = r.FormValue("title")

	g.Templates.New.Execute(w, r, data)
}
