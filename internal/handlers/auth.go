package handlers

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"TaskTracker/internal/database"
)

// handleLogout menangani logout pengguna
func HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "mahasiswa_npm",
		Value:  "",
		Path:   "/",
		MaxAge: -1, 
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// handleLogin menangani login pengguna
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			RenderLoginPage(w, "Error parsing form")
			return
		}

		mahasiswaNPM := r.FormValue("npm")
		password := r.FormValue("password")

		if mahasiswaNPM == "" || password == "" {
			RenderLoginPage(w, "Both NPM and Password are required")
			return
		}

		db := database.GetDB()
		var storedPassword string
		err = db.QueryRow("SELECT password FROM mahasiswa WHERE npm = ?", mahasiswaNPM).Scan(&storedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				RenderLoginPage(w, "Invalid NPM or Password")
			} else {
				RenderLoginPage(w, "Error querying database")
			}
			return
		}

		// Membandingkan password yang dimasukkan dengan yang ada di database
		err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
		if err != nil {
			RenderLoginPage(w, "Invalid NPM or Password")
			return
		}

		cookie := &http.Cookie{
			Name:     "mahasiswa_npm",
			Value:    mahasiswaNPM,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   3600, 
		}
		http.SetCookie(w, cookie)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	RenderLoginPage(w, "")
}

// renderLoginPage menampilkan halaman login
func RenderLoginPage(w http.ResponseWriter, errorMessage string) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	data := struct {
		ErrorMessage string
	}{ErrorMessage: errorMessage}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
