package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/daltbunker/soul_climbers/db"
	"github.com/daltbunker/soul_climbers/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed loading .env")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not found in env")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not found in env")
	}
	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		log.Fatal("SESSION_KEY not found in env")
	}
	handlers.Init(sessionKey)

	err = db.InitDatabase(dbURL)
	if err != nil {
		log.Fatal("Failed to init database")
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	handlers.InitPages()
	registerRoutes(router)

	srv := &http.Server{
		Handler: router,
		Addr: ":" + port,
	}

	fmt.Printf("Server running on port %v", port)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func registerRoutes(r *chi.Mux) {
	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// public
	r.Group(func(r chi.Router) {
		// page routes
		r.Get("/", handlers.HandleGetHome)
		r.Get("/login", handlers.HandleGetLogin)
		// r.Get("/login/reset/email", handlers.HandleGetResetEmail) // This is part of the password reset process, not to reset email
		// r.Get("/login/reset/password", handlers.HandleGetResetPassword)
		r.Get("/signup", handlers.HandleGetSignup)

		// data routes
		r.Get("/v1/healthz", handlers.HandleGetHealth)
		// r.Post("/v1/signup", handlers.HandleNewUser)
		r.Post("/v1/login", handlers.HandleLoginUser)
		// r.Post("/v1/login/reset/email", handlers.HandleEmailResetLink)
		// r.Post("/v1/login/reset/password", handlers.HandlePasswordReset)
		r.Get("/blog/{id}", handlers.HandleGetBlog)
	})

	// protected
	r.Group(func(r chi.Router) {
		// r.Use(sessionMiddleware) TODO: commented out for testing only
		r.Get("/admin", handlers.HandleGetAdmin)
		r.Get("/admin/blog", handlers.HandleGetBlogForm)
		r.Get("/admin/blog/{id}", handlers.HandleGetBlogForm)
		r.Get("/admin/blog/preview/{id}", handlers.HandleGetBlogPreview)
		r.Post("/admin/blog/preview/{id}", handlers.HandleUpdateBlogPreview) // HTML forms only allow GET and POST
		r.Post("/admin/blog/preview", handlers.HandleNewBlogPreview)
		r.Post("/admin/blog/{id}", handlers.HandlePublishBlog) // HTML forms only allow GET and POST
		r.Get("/v1/blog/{id}/{imgName}", handlers.GetBlogImg)
		r.Delete("/v1/blog/{id}/thumbnail", handlers.DeleteBlogImg)
	})
}

func sessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := handlers.GetSession(r)
		if err != nil {
			handlers.HandleServerError(w, err)
			return
		}

		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			handlers.HandleUnautorized(w, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}
