package handlers

import (
	"database/sql"
	"net/http"
	"TaskTracker/internal/database"
	"golang.org/x/crypto/bcrypt"
	"html/template"
)

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	var errorMessage string

	
	db := database.GetDB()

	if r.Method != http.MethodPost {
		tmpl, err := template.ParseFiles("templates/register.html")
		if err != nil {
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	npm := r.FormValue("npm")
	username := r.FormValue("username")
	password := r.FormValue("password")

	if npm == "" || username == "" || password == "" {
		http.Error(w, "Semua kolom harus diisi", http.StatusBadRequest)
		return
	}

	var existingNPM string
	err = db.QueryRow("SELECT npm FROM mahasiswa WHERE npm = ?", npm).Scan(&existingNPM)
	if err == nil {
		errorMessage = "NPM sudah terdaftar."
	} else if err != sql.ErrNoRows {
		http.Error(w, "Error checking NPM", http.StatusInternalServerError)
		return
	}

	if errorMessage != "" {
		tmpl, err := template.ParseFiles("templates/register.html")
		if err != nil {
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, struct {
			ErrorMessage string
		}{ErrorMessage: errorMessage})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO mahasiswa (npm, username, password) VALUES (?, ?, ?)", npm, username, string(hashedPassword))
	if err != nil {
		http.Error(w, "Failed to register", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}
