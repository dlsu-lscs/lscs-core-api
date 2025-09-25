package server

import (
	"net/http"
	"time"

	"github.com/dlsu-lscs/lscs-core-api/internal/handlers"
	"github.com/dlsu-lscs/lscs-core-api/internal/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(httprate.LimitByIP(100, time.Minute))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Public routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("it works"))
	})
	r.Post("/auth/google/callback", handlers.GoogleLoginHandler)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middlewares.JwtAuthentication)

		r.Get("/members", handlers.GetAllMembersHandler)
		r.Get("/committees", handlers.GetAllCommitteesHandler)
		r.Post("/member", handlers.GetMemberInfo)
		r.Post("/member-id", handlers.GetMemberInfoById)
		r.Post("/check-email", handlers.CheckEmailHandler)
		r.Post("/check-id", handlers.CheckIDIfMember)
	})

	return r
}

