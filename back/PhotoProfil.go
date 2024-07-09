package back

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func uploadProfilePictureHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "Authentication error: user ID cookie not found", http.StatusUnauthorized)
		return
	}
	userID := cookie.Value

	r.Body = http.MaxBytesReader(w, r.Body, 20<<20) // Limit to 20 MB
	if err := r.ParseMultipartForm(20 << 20); err != nil {
		http.Error(w, "The uploaded file is too big. Please choose a file that is less than 20MB in size.", http.StatusBadRequest)
		return
	}

	var currentImagePath sql.NullString
	err = db.QueryRow("SELECT profile_image FROM users WHERE id = ?", userID).Scan(&currentImagePath)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Failed to retrieve current profile image", http.StatusInternalServerError)
		return
	}

	file, header, err := r.FormFile("profile")
	if err != nil {
		http.Error(w, "Could not get uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	uuidWithHyphen := uuid.New()
	fileExt := filepath.Ext(header.Filename)
	newFilename := uuidWithHyphen.String() + fileExt

	dir := "uploads"
	filePath := filepath.Join(dir, newFilename)

	out, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Unable to create the file for writing.", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	filePath = strings.ReplaceAll(filePath, "\\", "/")

	if currentImagePath.Valid && currentImagePath.String != "" {
		oldFilePath := strings.ReplaceAll(currentImagePath.String, "\\", "/")
		os.Remove(oldFilePath)
	}

	_, err = db.Exec("UPDATE users SET profile_image = ? WHERE id = ?", filePath, userID)
	if err != nil {
		http.Error(w, "Failed to update user profile in database", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success":  true,
		"filePath": "/" + filePath,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
