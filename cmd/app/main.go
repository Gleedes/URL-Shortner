package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Gleedes/URL-Shortner/internal/config"
	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}
	fmt.Println("Server starting on port: ", cfg.ServerPort)

	r := chi.NewRouter()

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "pong"}`))
	})

	db, err := sql.Open("sqlite", cfg.DatabasePath)
	if err != nil {
		log.Fatal("не удалось открыть БД", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Отсутствует соединение с БД")
	}
	fmt.Println("Database connect")
	defer db.Close()

	data, err := os.ReadFile("./migrations/001_create_urls_table.up.sql")
	if err != nil {
		panic(err)
	}

	result, err := db.Exec(string(data))
	if err != nil {
		panic(err)
	}

	fmt.Println(result.RowsAffected())

	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		log.Fatal("Сервер упал: ", err)
	}
}
