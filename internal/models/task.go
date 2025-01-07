package models
	type Task struct {
		ID           int    `json:"id"`
		Matkul       string `json:"matkul"`
		Completed    bool   `json:"completed"`
		Deadline     string `json:"deadline"`
		MahasiswaNPM int    `json:"mahasiswa_npm"`
	}
