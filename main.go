package main

import (
	"log"
	"net/http"
	"TaskTracker/internal/database" // Import package database
	"TaskTracker/internal/handlers" // Import package handlers
)

func main() {

	
	// Inisialisasi koneksi database
	_, err := database.InitDB()
	if err != nil {
		log.Fatal("Database initialization failed: ", err)
	}

	// Menyiapkan route
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", handlers.ServeHome)
	http.HandleFunc("/register", handlers.HandleRegister)
	http.HandleFunc("/login", handlers.HandleLogin)
	http.HandleFunc("/logout", handlers.HandleLogout)
	http.HandleFunc("/add-task", handlers.HandleAddTask)
	http.HandleFunc("/update-status", handlers.HandleUpdateTaskStatus)
	http.HandleFunc("/delete-task", handlers.HandleDeleteTask)
	http.HandleFunc("/update-task", handlers.HandleUpdateTask)

	// Mulai server pada port yang ditentukan
	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
