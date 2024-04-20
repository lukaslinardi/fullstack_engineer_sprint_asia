package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type TaskRequest struct {
	TaskName string `json:"task_name"`
	ParentId *int   `json:"parent_id"`
}

type TaskResponse struct {
	ID         int       `json:"id" db:"id"`
	TaskName   string    `json:"task_name" db:"task_name"`
	TaskStatus int       `json:"task_status" db:"task_status"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

func main() {
	//connect to database
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//create the table if it doesn't exist
	// _, err = db.Exec("CREATE TABLE IF NOT EXISTS  (id SERIAL PRIMARY KEY, name TEXT, email TEXT)") if err != nil {
	// 	log.Fatal(err)
	// }

	//create router
	router := mux.NewRouter()
	router.HandleFunc("/task", getTasks(db)).Methods("GET")
	//	router.HandleFunc("/users/{id}", getTask(db)).Methods("GET")
	router.HandleFunc("/task", createTask(db)).Methods("POST")
	// router.HandleFunc("/users/{id}", updateTask(db)).Methods("PUT")
	// router.HandleFunc("/users/{id}", deleteTask(db)).Methods("DELETE")

	//start server
	log.Fatal(http.ListenAndServe(":8000", jsonContentTypeMiddleware(router)))
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// get all task
func getTasks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, task_name, task_status, created_at, updated_at FROM task")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		tasks := []TaskResponse{}
		for rows.Next() {
			var t TaskResponse
			if err := rows.Scan(&t.ID, &t.TaskName, &t.TaskStatus, &t.CreatedAt, &t.UpdatedAt); err != nil {
				log.Fatal(err)
			}
			tasks = append(tasks, t)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(tasks)
	}
}

// create user
func createTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t TaskRequest
		json.NewDecoder(r.Body).Decode(&t)

		err := db.QueryRow("INSERT INTO task (task_name, task_status, parent_id, created_at) VALUES ($1, $2, $3) RETURNING id",
			t.TaskName,
			1,
			time.Now().UTC())

		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(t)
	}
}

// get user by id
// func getTask(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
//         vars := mux.Vars(r)
// 		id := vars["id"]
//
// 		var u Task
// 		err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&u.ID, &u.TaskName, &u.TaskStatus)
// 		if err != nil {
// 			w.WriteHeader(http.StatusNotFound)
// 			return
// 		}
//
// 		json.NewEncoder(w).Encode(u)
// 	}
// }

//
// // update user
// func updateTask(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var u Task
// 		json.NewDecoder(r.Body).Decode(&u)
//
// 		vars := mux.Vars(r)
// 		id := vars["id"]
//
// 		_, err := db.Exec("UPDATE users SET name = $1, email = $2 WHERE id = $3", u.Name, u.Email, id)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		json.NewEncoder(w).Encode(u)
// 	}
// }
//
// // delete user
// func deleteTask(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		vars := mux.Vars(r)
// 		id := vars["id"]
//
// 		var u Task
// 		err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Email)
// 		if err != nil {
// 			w.WriteHeader(http.StatusNotFound)
// 			return
// 		} else {
// 			_, err := db.Exec("DELETE FROM users WHERE id = $1", id)
// 			if err != nil {
// 				//todo : fix error handling
// 				w.WriteHeader(http.StatusNotFound)
// 				return
// 			}
//
// 			json.NewEncoder(w).Encode("Task deleted")
// 		}
// 	}
// }
