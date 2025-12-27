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
	os.Exit(run())
}

func run() int {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Ошибка загрузки конфигурации: %v", err)
		return 1
	}

	auth := cfg.DBUser
	if cfg.DBPassword != "" {
		auth = fmt.Sprintf("%s:%s", cfg.DBUser, cfg.DBPassword)
	}
	dsn := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", auth, cfg.DBHost, cfg.DBPort, cfg.DBName)
	fmt.Println(dsn)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Printf("Не удалось подключиться к базе данных: %v", err)
		return 1
	}
	defer db.Close()

	for i := 1; i <= 5; i++ {
		err = db.Ping()
		if err == nil {
			log.Println("Успешное подключение к базе данных")
			break
		}
		log.Printf("Попытка %d: база данных недоступна, ожидание %d сек...", i, 2*i)
		time.Sleep(time.Duration(2*i) * time.Second)
	}

	if err != nil {
		log.Printf("Не удалось установить соединение с базой данных после 5 попыток: %v", err)
		return 1
	}
	log.Println("Подключение к PostgreSQL выполнено успешно")

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

	log.Printf("Запуск сервера на порту :%s...", cfg.ServerPort)
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}
	errChan := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case <-quit:
		log.Println("Получен сигнал завершения")
	case err := <-errChan:
		log.Printf("Сервер завершился с ошибкой: %v", err)
		return 1
	}

	log.Println("Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Ошибка при остановке сервера: %v", err)
		return 1
	}

	log.Println("Сервер успешно остановлен")
	log.Println("Приложение успешно завершило работу")
	return 0
}
