package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/just4fun-xd/task-manager/internal/api"
	"github.com/just4fun-xd/task-manager/internal/config"
	"github.com/just4fun-xd/task-manager/internal/task"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Config failed: %v", err)
	}
	auth := cfg.DBUser
	if cfg.DBPassword != "" {
		auth = fmt.Sprintf("%s:%s", cfg.DBUser, cfg.DBPassword)
	}
	dsn := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", auth, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()
	for i := 1; i <= 5; i++ {
		err = db.Ping()
		if err == nil {
			log.Println("Успешное подключение к базе данных!")
			break
		}
		log.Printf("Попытка %d: база данных недоступна, ждем %d сек...", i, 2*i)
		time.Sleep(time.Duration(2*i) * time.Second)
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
	log.Printf("Server starting on :%s...", cfg.ServerPort)
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Ошибка при остановке сервера: %v", err)
		return
	}
	log.Println("Сервер удачно остановлен, закрываем базу данных...")

	if err := db.Close(); err != nil {
		log.Printf("Ошибка при закрытии БД: %v", err)
	}
	log.Println("Приложение успешно завершило работу")
}
