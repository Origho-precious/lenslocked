package controllers

import (
	"fmt"
	appcontext "github/Origho-precious/lenslocked/context"
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

func (g Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		UserId int
		Title  string
	}

	data.UserId = int(appcontext.User(r.Context()).Id)
	data.Title = r.FormValue("title")

	gallery, err := g.GalleryService.Create(data.Title, data.UserId)

	if err != nil {
		g.Templates.New.Execute(w, r, data, err)
		return
	}
	// This page doesn't exist, but we will want to redirect here eventually.
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}
