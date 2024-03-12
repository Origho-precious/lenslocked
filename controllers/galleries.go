package controllers

import (
	"fmt"
	appcontext "github/Origho-precious/lenslocked/context"
	"github/Origho-precious/lenslocked/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Templates struct {
	New  Template
	Edit Template
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

	data.UserId = appcontext.User(r.Context()).Id
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

func (g Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		// 404 error - page isn't found.
		http.Error(w, "Invalid ID", http.StatusNotFound)
		return
	}

	gallery, err := g.GalleryService.ByID(id)
	if err != nil {
		if err == models.ErrNotFound {
			http.Error(w, "Gallery not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	user := appcontext.User(r.Context())
	if gallery.UserId != user.Id {
		http.Error(w, "You are not authorized to edit this gallery", http.StatusForbidden)
		return
	}

	var data struct {
		Id    int
		Title string
	}

	data.Id = gallery.ID
	data.Title = gallery.Title

	g.Templates.Edit.Execute(w, r, data)
}
