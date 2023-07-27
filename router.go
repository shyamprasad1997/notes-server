package main

import (
	"notes-server/db"
	"notes-server/loggers"
	"notes-server/middlewares"
	"notes-server/repositories"
	"sync"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

type IChiRouter interface {
	InitRouter() *chi.Mux
}

type router struct{}

func (router *router) InitRouter() *chi.Mux {
	notesController := ServiceContainer().InjectNotesController()
	loginController := ServiceContainer().InjectLoginController()

	r := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	logger := loggers.NewLogger()
	r.Route("/v1/api", func(r chi.Router) {
		r.Use(cors.Handler)
		r.Group(func(r chi.Router) {
			r.Use(middlewares.RequestID)
			r.Use(middleware.Recoverer)
			r.Use(middleware.Logger)
			r.Use(cors.Handler)
			r.Post("/signup", loginController.SignUp)
			r.Post("/login", loginController.Login)
			r.Route("/", func(r chi.Router) {
				r.Use(middlewares.TokenValidation(repositories.NewLoginRepository(db.NewDB(), logger), logger))
				r.Post("/notes", notesController.GetNotes) // need to make is post to send token in body
				r.Post("/note", notesController.AddNote)
				r.Delete("/note", notesController.DeleteNote)
			})
		})
	})
	return r
}

var (
	m          *router
	routerOnce sync.Once
)

func ChiRouter() IChiRouter {
	if m == nil {
		routerOnce.Do(func() {
			m = &router{}
		})
	}
	return m
}
