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
	dbURL := os.Getenv("DB_URL_DVL") // _PROD | _DVL
	if dbURL == "" {
		log.Fatal("DB_URL not found in env")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not found in env")
	}

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

	registerPublicRoutes(router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	fmt.Printf("Server running on port %v", port)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func registerPublicRoutes(r *chi.Mux) {
	publicRouter := chi.NewRouter()

	fs := http.FileServer(http.Dir("./static"))
	publicRouter.Handle("/static/*", http.StripPrefix("/static/", fs))

	// routes 
	publicRouter.Get("/", handlers.HandleGetHome)
	publicRouter.Get("/login", handlers.HandleGetLogin)
	publicRouter.Get("/signup", handlers.HandleGetSignup)

	// data requests 
	publicRouter.Get("/v1/healthz", handlers.HandleGetHealth)
	publicRouter.Post("/v1/signup", handlers.HandleNewUser)
	publicRouter.Post("/v1/login", handlers.HandleLoginUser)
	publicRouter.Post("/v1/validate/email", handlers.HandleEmailValidate)
	publicRouter.Post("/v1/validate/username", handlers.HandleUsernameValidate)
	publicRouter.Post("/v1/validate/password", handlers.HandlePasswordValidate)

	r.Mount("/", publicRouter)
}
