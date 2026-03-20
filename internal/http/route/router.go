package route

import (
	handler "test/internal/http"

	"github.com/go-chi/chi/v5"
)

func WithObjectHandlers(r chi.Router, handler *handler.Object) {
	r.Route("/task", func(r chi.Router) {
		r.Use(handler.AuthMiddleware)                       // Применяем middleware для всех маршрутов /task
		r.Post("/", handler.PostTaskHandlerWithFilters)     // POST /task
		r.Get("/status/{taskID}", handler.GetStatusHandler) // GET /task/status/{taskID}
		r.Get("/result/{taskID}", handler.GetResultHandler) // GET /task/result/{taskID}
	})
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", handler.RegisterHandler) // POST /auth/register
		r.Post("/login", handler.LoginHandler)       // POST /auth/login
	})
	r.Route("/commit", func(r chi.Router) {
		r.Post("/", handler.CommitHandler) // POST /auth/register
	})
}
