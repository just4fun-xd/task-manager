package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/just4fun-xd/task-manager/internal/api"
	"github.com/just4fun-xd/task-manager/internal/task"
)

func main() {
	dsn := "postgres://shalygin@127.0.0.1:5432/tasks?sslmode=disable"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Database is unreachable: %v", err)
	}
	log.Println("Connected to PostgreSQL successfully!")

	repo := task.NewPostgresRepository(db)
	service := task.NewService(repo)
	handler := api.NewHandler(service)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/tasks", func(r chi.Router) {
		r.Post("/", handler.CreateTask)
		r.Get("/", handler.GetAllTasks)
		r.Get("/{id}", handler.GetTask)
		r.Put("/{id}", handler.UpdateTask)
		r.Delete("/{id}", handler.DeleteTask)
	})
	log.Printf("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
