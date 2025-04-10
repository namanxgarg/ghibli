package upload

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/namanxgarg/ghibli-backend/pkg/auth"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10 MB max
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "file too big or invalid form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "could not read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file type (extension)
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		http.Error(w, "only .jpg, .jpeg and .png allowed", http.StatusUnsupportedMediaType)
		return
	}

	// Get user ID from context
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Save the file locally
	filename := fmt.Sprintf("%s_%s%s", userID, uuid.New().String(), ext)
	outPath := filepath.Join("uploads", filename)

	outFile, err := os.Create(outPath)
	if err != nil {
		http.Error(w, "could not save file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "error saving file", http.StatusInternalServerError)
		return
	}

	SaveUpload(Upload{
		Filename:   filename,
		UserID:     userID,
		UploadedAt: time.Now(),
	})
	
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Image uploaded as %s", filename)
}

func ListUserUploadsHandler(w http.ResponseWriter, r *http.Request) {
    userID, ok := auth.GetUserIDFromContext(r.Context())
    if !ok {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }

    userUploads := GetUploadsByUser(userID)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(userUploads)
}

