package web

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gowebexamples/goreddit"
)

func NewHandler(store goreddit.Store) *Handler {
	h := &Handler{
		Mux:   chi.NewMux(),
		store: store,
	}

	threads := ThreadHandler{store: store}
	posts := PostHandler{store: store}
	comments := CommentHandler{store: store}

	h.Use(middleware.Logger)

	h.Get("/", h.Home())
	h.Route("/threads", func(r chi.Router) {
		r.Get("/", threads.List())
		r.Get("/new", threads.Create())
		r.Post("/", threads.Store())
		r.Get("/{id}", threads.Show())
		r.Post("/{id}/delete", threads.Delete())
		r.Get("/{id}/new", posts.Create())
		r.Post("/{id}", posts.Store())
		r.Get("/{threadID}/{postID}", posts.Show())
		r.Get("/{threadID}/{postID}/vote", posts.Vote())
		r.Post("/{threadID}/{postID}", comments.Store())
	})
	h.Get("/comments/{id}/vote", comments.Vote())

	return h
}

type Handler struct {
	*chi.Mux

	store goreddit.Store
}

func (h *Handler) Home() http.HandlerFunc {
	type data struct {
		Posts []goreddit.Post
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/home.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		pp, err := h.store.Posts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data{Posts: pp})
	}
}
