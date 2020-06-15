package web

import (
	"html/template"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/csrf"
	"github.com/gowebexamples/goreddit"
)

func NewHandler(store goreddit.Store, sessions *scs.SessionManager, csrfKey []byte) *Handler {
	h := &Handler{
		Mux:      chi.NewMux(),
		store:    store,
		sessions: sessions,
	}

	threads := ThreadHandler{store: store, sessions: sessions}
	posts := PostHandler{store: store, sessions: sessions}
	comments := CommentHandler{store: store, sessions: sessions}

	h.Use(middleware.Logger)
	h.Use(csrf.Protect(csrfKey, csrf.Secure(false)))
	h.Use(sessions.LoadAndSave)

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

	store    goreddit.Store
	sessions *scs.SessionManager
}

func (h *Handler) Home() http.HandlerFunc {
	type data struct {
		SessionData

		Posts []goreddit.Post
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/home.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		pp, err := h.store.Posts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data{
			SessionData: GetSessionData(h.sessions, r.Context()),
			Posts:       pp,
		})
	}
}
