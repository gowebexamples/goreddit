package web

import (
	"html/template"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/gowebexamples/goreddit"
)

type PostHandler struct {
	store    goreddit.Store
	sessions *scs.SessionManager
}

func (h *PostHandler) Create() http.HandlerFunc {
	type data struct {
		SessionData

		CSRF   template.HTML
		Thread goreddit.Thread
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/post_create.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		t, err := h.store.Thread(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data{
			SessionData: GetSessionData(h.sessions, r.Context()),
			CSRF:        csrf.TemplateField(r),
			Thread:      t,
		})
	}
}

func (h *PostHandler) Show() http.HandlerFunc {
	type data struct {
		SessionData
		CSRF     template.HTML
		Thread   goreddit.Thread
		Post     goreddit.Post
		Comments []goreddit.Comment
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/post.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		postIDStr := chi.URLParam(r, "postID")
		threadIDStr := chi.URLParam(r, "threadID")

		postID, err := uuid.Parse(postIDStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		threadID, err := uuid.Parse(threadIDStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		p, err := h.store.Post(postID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cc, err := h.store.CommentsByPost(p.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		t, err := h.store.Thread(threadID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data{
			SessionData: GetSessionData(h.sessions, r.Context()),
			CSRF:        csrf.TemplateField(r),
			Thread:      t,
			Post:        p,
			Comments:    cc,
		})
	}
}

func (h *PostHandler) Store() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.FormValue("title")
		content := r.FormValue("content")

		idStr := chi.URLParam(r, "id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		t, err := h.store.Thread(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		p := &goreddit.Post{
			ID:       uuid.New(),
			ThreadID: t.ID,
			Title:    title,
			Content:  content,
		}
		if err := h.store.CreatePost(p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		h.sessions.Put(r.Context(), "flash", "Your post has been created.")

		http.Redirect(w, r, "/threads/"+t.ID.String()+"/"+p.ID.String(), http.StatusFound)
	}
}

func (h *PostHandler) Vote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "postID")

		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		p, err := h.store.Post(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		dir := r.URL.Query().Get("dir")
		if dir == "up" {
			p.Votes++
		} else if dir == "down" {
			p.Votes--
		}

		if err := h.store.UpdatePost(&p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, r.Referer(), http.StatusFound)
	}
}
