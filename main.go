package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/lukaslinardi/fullstack_engineer_sprint_asia/types"
	"github.com/rs/cors"
)

func main() {
	// connect to database
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},  // Include PUT method
		AllowedHeaders:   []string{"Content-Type", "Authorization"}, // Add any custom headers you need
		AllowCredentials: true,
		Debug:            true,
	})


	router := mux.NewRouter()
	router.HandleFunc("/tasks", getTasks(db)).Methods("GET")
	router.HandleFunc("/sub-tasks/{id}", createSubTask(db)).Methods("POST")
	router.HandleFunc("/tasks", createTask(db)).Methods("POST")
	router.HandleFunc("/tasks/{id}", updateTask(db)).Methods("PUT")
	router.HandleFunc("/update-tasks/{id}", updateTaskStatus(db)).Methods("PUT")
	router.HandleFunc("/tasks/{id}", deleteTask(db)).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", c.Handler(jsonContentTypeMiddleware(router))))
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func getTasks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`SELECT
          t1.id,
          t1.task_name,
          t1.task_status,
          t1.deadline_task,
          t1.parent_id,
          tmp.completion_percentage * 100,
          t1.created_at,
          t1.updated_at
        FROM task t1
        LEFT JOIN LATERAL (
          SELECT
            (
              SELECT count(t3.*) FROM task t3
              WHERE t3.task_status = 2 AND t3.parent_id = t1.id
            )::float / NULLIF(count(t2.*), 0) as completion_percentage
          FROM task t2
          WHERE t2.parent_id = t1.id
        ) tmp ON true`)
		if err != nil {
			log.Fatal(err)
		}

		defer rows.Close()

		tasks := []types.TaskResponse{}
		for rows.Next() {
			var t types.TaskResponse
			if err := rows.Scan(&t.ID, &t.TaskName, &t.TaskStatus, &t.DeadlineTask, &t.ParentId, &t.Percentage, &t.CreatedAt, &t.UpdatedAt); err != nil {
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

func createTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t types.TaskRequest
		json.NewDecoder(r.Body).Decode(&t)

		_, err := db.Exec("INSERT INTO task (task_name, task_status, deadline_task, parent_id, created_at) VALUES ($1, $2, $3, $4, $5)",
			t.TaskName,
			1,
			t.DeadlineTask,
			t.ParentId,
			time.Now().UTC())
		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func updateTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t types.TaskUpdate
		json.NewDecoder(r.Body).Decode(&t)
		vars := mux.Vars(r)
		id := vars["id"]

		_, err := db.Exec("UPDATE task SET task_name = $1, deadline_task = $2 WHERE id = $3", t.TaskName, t.DeadlineTask, id)
		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func updateTaskStatus(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var parentId *int

		err := db.QueryRow("SELECT parent_id from task where id = $1", id).Scan(&parentId)
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("UPDATE task SET task_status = 2 WHERE id = $1", id)
		if err != nil {
			log.Fatal(err)
		}

		rows, err := db.Query("SELECT t.task_status from task t where t.parent_id = $1", parentId)
		if err != nil {
			log.Fatal(err)
		}

		isAllSubTaskCompleted := true

		tasks := []int{}
		for rows.Next() {
			var status int
			if err := rows.Scan(&status); err != nil {
				log.Fatal(err)
			}
			tasks = append(tasks, status)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		for _, task := range tasks {
			if task != 2 {
				isAllSubTaskCompleted = false
				break
			}
		}

		if isAllSubTaskCompleted {
			_, err := db.Exec("UPDATE task SET task_status = 2 WHERE id = $1", parentId)
			if err != nil {
				log.Fatal(err)
			}

		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func deleteTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		_, err := db.Exec("DELETE FROM task WHERE id = $1", id)
		if err != nil {
			log.Fatal(err)
			//			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		//		json.NewEncoder(w).Encode("Task deleted")
	}
}

func createSubTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t types.TaskRequest
		json.NewDecoder(r.Body).Decode(&t)
		_, err := db.Exec("INSERT INTO task (task_name, task_status, deadline_task, parent_id, created_at) VALUES ($1, $2, $3, $4, $5)",
			t.TaskName,
			1,
			t.DeadlineTask,
			t.ParentId,
			time.Now().UTC())
		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
