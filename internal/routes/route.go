package routes

import (
	handler "github.com/SaikatDeb12/storeX/internal/handlers"
	"github.com/SaikatDeb12/storeX/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func SetUpRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Route("/v1", func(v1 chi.Router) {
		v1.Get("/health", handler.CheckHealth)
		v1.Route("/auth", func(r chi.Router) {
			r.Post("/login", handler.Login)
			r.Post("/Register", handler.Register)
			r.Group(func(r chi.Router) {
				r.Use(middleware.Authenticate)
				r.Post("/logout", handler.Logout)
			})
		})

		v1.Group(func(r chi.Router) {
			r.Use(middleware.Authenticate)
			r.Route("/users", func(r chi.Router) {
				r.Get("/", handler.GetAllUsers)
				r.Get("/{id}", handler.GetUserInfoByID)
				r.Group(func(r chi.Router) {
					r.Use(middleware.CheckUserRole)
					// r.Delete("/{id}", handler)
				})
			})
			r.Route("/assets", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Use(middleware.CheckUserRole)
					r.Get("/", handler.FetchAssets)
					r.Post("/", handler.CreateAsset)
					r.Put("/update/{id}", handler.UpdateAsset)
					r.Put("/assign", handler.AssignedAssets)
					r.Put("/service/{id}", handler.SentToService)
				})
			})
		})
	})

	return router
}
