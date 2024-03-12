package controllers

import (
	"errors"
	"fmt"
	appcontext "github/Origho-precious/lenslocked/context"
	"github/Origho-precious/lenslocked/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Templates struct {
	New   Template
	Edit  Template
	Show  Template
	Index Template
}

type Galleries struct {
	Templates      Templates
	GalleryService *models.GalleryService
}

type galleryOpt func(http.ResponseWriter, *http.Request, *models.Gallery) error

var dogImages = []string{
	"https://images.unsplash.com/photo-1591160690555-5debfba289f0?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3wxNDA1Mzh8MHwxfHJhbmRvbXx8fHx8fHx8fDE3MTAyNTMyMjB8&ixlib=rb-4.0.3&q=80&w=1080",
	"https://images.unsplash.com/photo-1589941013453-ec89f33b5e95?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3wxNDA1Mzh8MHwxfHJhbmRvbXx8fHx8fHx8fDE3MTAyNTMyMjB8&ixlib=rb-4.0.3&q=80&w=1080",
	"https://images.unsplash.com/photo-1585908286456-991b5d0e53f4?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3wxNDA1Mzh8MHwxfHJhbmRvbXx8fHx8fHx8fDE3MTAyNTMyMjB8&ixlib=rb-4.0.3&q=80&w=1080",
	"https://images.unsplash.com/photo-1618173745201-8e3bf8978acc?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3wxNDA1Mzh8MHwxfHJhbmRvbXx8fHx8fHx8fDE3MTAyNTMyMjB8&ixlib=rb-4.0.3&q=80&w=1080",
	"https://images.unsplash.com/photo-1633722715463-d30f4f325e24?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3wxNDA1Mzh8MHwxfHJhbmRvbXx8fHx8fHx8fDE3MTAyNTMyMjB8&ixlib=rb-4.0.3&q=80&w=1080",
	"https://images.unsplash.com/photo-1568640347023-a616a30bc3bd?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3wxNDA1Mzh8MHwxfHJhbmRvbXx8fHx8fHx8fDE3MTAyNTMyMjB8&ixlib=rb-4.0.3&q=80&w=1080",
	"https://images.unsplash.com/photo-1534551767192-78b8dd45b51b?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3wxNDA1Mzh8MHwxfHJhbmRvbXx8fHx8fHx8fDE3MTAyNTMyMjB8&ixlib=rb-4.0.3&q=80&w=1080",
	"https://images.unsplash.com/photo-1540411003967-af56b79be677?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3wxNDA1Mzh8MHwxfHJhbmRvbXx8fHx8fHx8fDE3MTAyNTMyMjB8&ixlib=rb-4.0.3&q=80&w=1080",
	"https://images.unsplash.com/photo-1612940960267-4549a58fb257?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3wxNDA1Mzh8MHwxfHJhbmRvbXx8fHx8fHx8fDE3MTAyNTMyMjB8&ixlib=rb-4.0.3&q=80&w=1080",
	"https://images.unsplash.com/photo-1600077029182-92ac8906f9a3?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3wxNDA1Mzh8MHwxfHJhbmRvbXx8fHx8fHx8fDE3MTAyNTMyMjB8&ixlib=rb-4.0.3&q=80&w=1080",
	"https://images.unsplash.com/flagged/photo-1550973078-10a2d124c99c?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3wxNDA1Mzh8MHwxfHJhbmRvbXx8fHx8fHx8fDE3MTAyNTMyMjB8&ixlib=rb-4.0.3&q=80&w=1080",
	"https://images.unsplash.com/photo-1574760112346-8443c3773437?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3wxNDA1Mzh8MHwxfHJhbmRvbXx8fHx8fHx8fDE3MTAyNTMyMjB8&ixlib=rb-4.0.3&q=80&w=1080",
	"https://images.unsplash.com/photo-1598875184988-5e67b1a874b8?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3wxNDA1Mzh8MHwxfHJhbmRvbXx8fHx8fHx8fDE3MTAyNTMyMjB8&ixlib=rb-4.0.3&q=80&w=1080",
	"https://images.unsplash.com/photo-1553776590-89774e24b34a?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3wxNDA1Mzh8MHwxfHJhbmRvbXx8fHx8fHx8fDE3MTAyNTMyMjB8&ixlib=rb-4.0.3&q=80&w=1080",
	"https://images.unsplash.com/photo-1557495235-340eb888a9fb?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3wxNDA1Mzh8MHwxfHJhbmRvbXx8fHx8fHx8fDE3MTAyNTMyMjB8&ixlib=rb-4.0.3&q=80&w=1080",
	"https://images.unsplash.com/photo-1591946614720-90a587da4a36?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3wxNDA1Mzh8MHwxfHJhbmRvbXx8fHx8fHx8fDE3MTAyNTMyMjB8&ixlib=rb-4.0.3&q=80&w=1080",
	"https://images.unsplash.com/photo-1600077106724-946750eeaf3c?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3wxNDA1Mzh8MHwxfHJhbmRvbXx8fHx8fHx8fDE3MTAyNTMyMjB8&ixlib=rb-4.0.3&q=80&w=1080",
	"https://images.unsplash.com/photo-1520580413066-ac45756bdc71?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3wxNDA1Mzh8MHwxfHJhbmRvbXx8fHx8fHx8fDE3MTAyNTMyMjB8&ixlib=rb-4.0.3&q=80&w=1080",
	"https://images.unsplash.com/photo-1502673530728-f79b4cab31b1?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3wxNDA1Mzh8MHwxfHJhbmRvbXx8fHx8fHx8fDE3MTAyNTMyMjB8&ixlib=rb-4.0.3&q=80&w=1080",
	"https://images.unsplash.com/photo-1514984879728-be0aff75a6e8?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3wxNDA1Mzh8MHwxfHJhbmRvbXx8fHx8fHx8fDE3MTAyNTMyMjB8&ixlib=rb-4.0.3&q=80&w=1080",
}

func userMustOwnGallery(
	w http.ResponseWriter, r *http.Request, gallery *models.Gallery,
) error {
	user := appcontext.User(r.Context())
	if user.ID != gallery.UserID {
		http.Error(w, "You are not authorized to edit this gallery", http.StatusForbidden)
		return fmt.Errorf("user does not have access to this gallery")
	}

	return nil
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
		UserID int
		Title  string
	}

	data.UserID = appcontext.User(r.Context()).ID
	data.Title = r.FormValue("title")

	gallery, err := g.GalleryService.Create(data.Title, data.UserID)

	if err != nil {
		g.Templates.New.Execute(w, r, data, err)
		return
	}
	// This page doesn't exist, but we will want to redirect here eventually.
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (g Galleries) galleryByID(
	w http.ResponseWriter, r *http.Request, opts ...galleryOpt,
) (*models.Gallery, error) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusNotFound)
		return nil, err
	}

	gallery, err := g.GalleryService.ByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Gallery not found", http.StatusNotFound)
			return nil, err
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return nil, err
	}

	for _, opt := range opts {
		err = opt(w, r, gallery)
		if err != nil {
			return nil, err
		}
	}

	return gallery, nil
}

func (g Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}

	var data struct {
		ID    int
		Title string
	}

	data.ID = gallery.ID
	data.Title = gallery.Title

	g.Templates.Edit.Execute(w, r, data)
}

func (g Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}

	title := r.FormValue("title")
	gallery.Title = title

	err = g.GalleryService.Update(gallery)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (g Galleries) Index(w http.ResponseWriter, r *http.Request) {
	type Gallery struct {
		ID    int
		Title string
	}

	var data struct {
		Galleries []Gallery
	}

	user := appcontext.User(r.Context())

	galleries, err := g.GalleryService.ByUserID(user.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	for _, gallery := range galleries {
		data.Galleries = append(data.Galleries, Gallery{
			ID:    gallery.ID,
			Title: gallery.Title,
		})
	}

	g.Templates.Index.Execute(w, r, data)
}

func (g Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	var data struct {
		ID     int
		Title  string
		Images []string
	}

	data.ID = gallery.ID
	data.Title = gallery.Title

	data.Images = append(data.Images, dogImages...)

	g.Templates.Show.Execute(w, r, data)
}
