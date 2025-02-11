package main

import (
	"log"
	"net/http"
	"os"

	"github.com/daltbunker/soul_climbers/db"
	"github.com/daltbunker/soul_climbers/handlers"
	"github.com/daltbunker/soul_climbers/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed loading .env: %v", err)
	}
	utils.ValidateEnv()
	db.InitDatabase(os.Getenv("DB_URL"))
	handlers.InitPages()
	handlers.InitSession(os.Getenv("SESSION_KEY"))

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	registerRoutes(router)

	port := os.Getenv("PORT")

	srv := &http.Server{
		Handler: router,
		Addr: ":" + port,
	}

	log.Printf("Server running on port %v", port)

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
		r.Get("/blog/{id}", handlers.HandleGetBlog)

		// data routes
		r.Get("/v1/healthz", handlers.HandleGetHealth)
		r.Post("/v1/signout", handlers.HandleSignOut)
		r.Post("/v1/signup", handlers.HandleNewUser)
		r.Post("/v1/login", handlers.HandleLoginUser)
		// r.Post("/v1/login/reset/email", handlers.HandleEmailResetLink)
		// r.Post("/v1/login/reset/password", handlers.HandlePasswordReset)
		r.Get("/v1/blog/{id}/{imgName}", handlers.GetBlogImg)
	})

	// protected (paths with "/admin" are restricted to admin role)
	r.Group(func(r chi.Router) {
		r.Use(handlers.SessionMiddleware)

		// page routes
		r.Get("/account", handlers.HandleGetAccount)
		r.Get("/placement-test", handlers.HandleGetPlacementTest)
		r.Get("/admin", handlers.HandleGetAdmin)
		r.Get("/admin/blog", handlers.HandleGetBlogForm)
		r.Get("/admin/blog/{id}", handlers.HandleGetBlogForm)
		r.Get("/admin/blog/preview/{id}", handlers.HandleGetBlogPreview)
		r.Post("/admin/blog/preview/{id}", handlers.HandleUpdateBlogPreview) // HTML forms only allow GET and POST
		r.Post("/admin/blog/preview", handlers.HandleNewBlogPreview)
		r.Post("/admin/blog/{id}", handlers.HandlePublishBlog) // HTML forms only allow GET and POST

		// data routes
		r.Get("/v1/admin/blog/{id}/{imgName}", handlers.GetBlogImg)
		r.Delete("/v1/admin/blog/{id}/thumbnail", handlers.DeleteBlogImg)
		r.Post("/v1/placement-test", handlers.HandleSubmitPlacementTest)
		r.Delete("/v1/admin/blog/{id}", handlers.DeleteBlog)
	})
}
