package main

import (
    "fmt"
    "net/http"
	"encoding/json" 
    "github.com/go-chi/chi/v5"
    "github.com/namanxgarg/ghibli-backend/internal/user"
	"github.com/namanxgarg/ghibli-backend/pkg/auth"
	"github.com/namanxgarg/ghibli-backend/internal/upload"
	"path/filepath"


)

func main() {
    r := chi.NewRouter()

    r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprintln(w, "ok")
    })

    r.Post("/signup", user.SignupHandler)
	r.Post("/login", user.LoginHandler)

	r.Group(func(r chi.Router) {
		r.Use(auth.JWTMiddleware)
	
		r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
			userID, ok := auth.GetUserIDFromContext(r.Context())
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
	
			user, ok := user.FindUserByID(userID)
			if !ok {
				http.Error(w, "user not found", http.StatusNotFound)
				return
			}
	
			json.NewEncoder(w).Encode(user)
		})
		r.Post("/upload", upload.UploadHandler)
		r.Get("/my-uploads", upload.ListUserUploadsHandler)

	})

	//we want it to be public, so outside jwt group
	r.Get("/images/{filename}", func(w http.ResponseWriter, r *http.Request) {
		filename := chi.URLParam(r, "filename")
		http.ServeFile(w, r, filepath.Join("uploads", filename))
	})	


    fmt.Println("🚀 Server running on http://localhost:8080")
    err := http.ListenAndServe(":8080", r)
    if err != nil {
        fmt.Println("❌ Failed to start server:", err)
    }
}
