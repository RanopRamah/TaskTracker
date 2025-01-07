package handlers

	import(

	"net/http"
	"html/template"
	"strconv"
	"log"
	"encoding/json"
	"time"
	"TaskTracker/internal/database"
	"TaskTracker/internal/models"
	"TaskTracker/internal/services"
	)

func HandleAddTask(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("mahasiswa_npm")
		if err != nil || cookie == nil {
			log.Println("Cookie 'mahasiswa_npm' tidak ditemukan atau error:", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		log.Println("Cookie 'mahasiswa_npm' ditemukan dengan nilai:", cookie.Value)

		mahasiswaNPM, err := strconv.Atoi(cookie.Value)
		if err != nil {
			http.Error(w, "Invalid NPM", http.StatusBadRequest)
			log.Println("Error parsing mahasiswa_npm:", err) // Log error if parsing fails
			return
		}

		log.Println("Parsed mahasiswa_npm:", mahasiswaNPM) // Log the parsed NPM

		// If method is POST, we add the task
		if r.Method == http.MethodPost {
			log.Println("Content-Type:", r.Header.Get("Content-Type"))

			var task models.Task
			err := json.NewDecoder(r.Body).Decode(&task)
			if err != nil {
				http.Error(w, "Invalid JSON data", http.StatusBadRequest)
				log.Println("Error decoding JSON:", err)
				return
			}

			log.Printf("Received Task Data: %+v\n", task)

			if task.Matkul == "" || task.Deadline == "" {
				http.Error(w, "Matkul or deadline cannot be empty", http.StatusBadRequest)
				return
			}

			layout := "2006-01-02T15:04:05"

			if len(task.Deadline) == 10 {
				task.Deadline += "T00:00:00"
			}
			parsedDeadline, err := time.Parse(layout, task.Deadline)
			if err != nil {
				http.Error(w, "Invalid deadline format", http.StatusBadRequest)
				log.Println("Error parsing deadline:", err)
				return
			}

			task.Deadline = parsedDeadline.Format("2006-01-02 15:04:05")

			err = services.InsertTaskToDB(task, mahasiswaNPM)
			if err != nil {
				http.Error(w, "Error inserting task into database", http.StatusInternalServerError)
				log.Println("Error inserting task into DB:", err)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "success"})
			return
		}

		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "error loading template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	}


	
	
	func HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
		taskID := r.URL.Query().Get("id")

		if taskID == "" {
			http.Error(w, "Missing task ID", http.StatusBadRequest)
			return
		}

		err := services.DeleteTaskFromDB(taskID)
		if err != nil {
			http.Error(w, "Failed to delete task", http.StatusInternalServerError)
			return
		}

		response := map[string]string{"status": "success"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}

func HandleUpdateTask(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	if task.Matkul == "" || task.Deadline == "" {
		http.Error(w, "New task description and deadline are required", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE tasks SET matkul = ?, deadline = ? WHERE id = ?", task.Matkul, task.Deadline, id)
	if err != nil {
		http.Error(w, "Error updating task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"status":  "success",
		"message": "Task updated successfully",
	}
	
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}




	func ServeHome(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
    cookie, err := r.Cookie("mahasiswa_npm")
    if err != nil || cookie == nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    mahasiswaNPM, err := strconv.Atoi(cookie.Value)
    if err != nil {
        http.Error(w, "Invalid NPM", http.StatusBadRequest)
        return
    }

    var mahasiswaName string
    err = db.QueryRow("SELECT username FROM mahasiswa WHERE npm = ?", mahasiswaNPM).Scan(&mahasiswaName)
    if err != nil {
        http.Error(w, "Error fetching student data", http.StatusInternalServerError)
        return
    }

    tasks, err := services.GetTasksFromDB(mahasiswaNPM)
    if err != nil {
        http.Error(w, "Error fetching tasks", http.StatusInternalServerError)
        return
    }

    tmpl, err := template.ParseFiles("templates/index.html")
    if err != nil {
        http.Error(w, "Error loading template", http.StatusInternalServerError)
        return
    }

   
    err = tmpl.Execute(w, struct {
        Username string
        Tasks    []models.Task
    }{
        Username: mahasiswaName,
        Tasks:    tasks,
    })
    if err != nil {
        http.Error(w, "Error rendering template", http.StatusInternalServerError)
        return
    }
}