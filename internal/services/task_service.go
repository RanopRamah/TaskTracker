package services

import(
	"fmt"
	"database/sql"
	"log"
	"TaskTracker/internal/database"
	"TaskTracker/internal/models"
)

func GetTasksFromDB(mahasiswaNPM int) ([]models.Task, error) {
		db := database.GetDB()
		query := "SELECT id, matkul, completed, deadline FROM tasks WHERE mahasiswa_npm = ?"
		rows, err := db.Query(query, mahasiswaNPM)
		if err != nil {
			return nil, fmt.Errorf("error querying database: %v", err)
		}
		defer rows.Close()

		var tasks []models.Task
		for rows.Next() {
			var task models.Task
			var deadline sql.NullString
			err := rows.Scan(&task.ID, &task.Matkul, &task.Completed, &deadline)
			if err != nil {
				return nil, fmt.Errorf("error scanning row: %v", err)
			}

			if deadline.Valid {
				task.Deadline = deadline.String
			} else {
				task.Deadline = ""
			}

			tasks = append(tasks, task)
		}

		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("error iterating over rows: %v", err)
		}

		return tasks, nil
	}

	func InsertTaskToDB(task models.Task, mahasiswaNPM int) error {
		db := database.GetDB()
		
		query := "INSERT INTO tasks (matkul, completed, deadline, mahasiswa_npm) VALUES (?, ?, ?, ?)"
		
		_, err := db.Exec(query, task.Matkul, task.Completed, task.Deadline, mahasiswaNPM)
		if err != nil {
			return fmt.Errorf("error inserting task into database: %v", err)
		}

		log.Printf("Task successfully added to the database: %+v\n", task)
		return nil
	}

		func DeleteTaskFromDB(taskID string) error {
		db := database.GetDB()
		query := "DELETE FROM tasks WHERE id = ?"
		_, err := db.Exec(query, taskID)
		return err
	}

