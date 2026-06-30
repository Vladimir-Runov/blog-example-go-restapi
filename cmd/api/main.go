package main

import (
	"blog-example-go-restapi/internal/handler"
	"blog-example-go-restapi/internal/repository"
	"blog-example-go-restapi/internal/service"
	"blog-example-go-restapi/pkg/auth"
	"blog-example-go-restapi/pkg/database"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем конфигурацию из .env файла
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}

	// TODO: Загрузить конфигурацию из переменных окружения
	cfg := loadConfig()

	// TODO: Подключиться к базе данных
	// - Создать database.Config из параметров конфигурации
	dbConfig := database.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
	}

	log.Printf("%v", dbConfig)
	//connStr := "host=localhost port=5433 user=bloguser password=blogpassword dbname=blogdb sslmode=disable"

	// - Вызвать database.NewPostgresDB
	// - Обработать ошибки подключения
	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	// - Не забыть defer db.Close()
	defer db.Close()

	// TODO: Выполнить миграции базы данных
	// - Вызвать database.Migrate(db)
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Could not migrate the database: %v", err)
	}

	// TODO: Инициализировать JWT менеджер
	// - Создать jwtManager через auth.NewJWTManager
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiryHours)

	// TODO: Создать слои приложения
	// 1. Репозитории (передать db)
	userRepo := repository.NewUserRepo(db)

	// 2. Сервисы (передать репозитории и jwtManager)
	userService := service.NewUserService(userRepo, jwtManager)

	// 3. Хендлеры (передать сервисы)
	userHandler := handler.NewAuthHandler(userService)

	// 4. Middleware (передать необходимые зависимости)

	// Настраиваем маршруты
	router := chi.NewRouter()
	// TODO: Настроить middleware

	// - Добавить глобальные middleware (logging, recovery, CORS)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.AllowContentType("application/json")) // CORS middleware

	// TODO: Настроить маршруты
	// Публичные эндпоинты:
	router.Post("/api/register", userHandler.Register) // - POST /api/register
	router.Post("/api/login", userHandler.Login)       // - POST /api/login
	// - GET /api/posts
	// - GET /api/posts/{id}
	// - GET /api/posts/{id}/comments
	//
	// Защищенные эндпоинты (требуют JWT):
	// - POST /api/posts
	// - PUT /api/posts/{id}
	// - DELETE /api/posts/{id}
	// - POST /api/posts/{id}/comments

	// Health check эндпоинт
	router.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"blog-api"}`))
	})

	// TODO: Запустить HTTP сервер
	// - Сформировать адрес из конфигурации
	addr := fmt.Sprintf("%s:%d", cfg.DBHost, cfg.DBPort)

	// - Вывести информацию о запуске
	log.Printf("Starting server on %s...", addr)
	// - Запустить сервер и обработать ошибки
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}

// Config представляет конфигурацию приложения
type Config struct {
	// Server
	ServerHost string
	ServerPort int

	// Database
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// JWT
	JWTSecret      string
	JWTExpiryHours int

	// Cache
	CacheTTLMinutes int
}

// loadConfig загружает конфигурацию из переменных окружения
func loadConfig() *Config {
	// TODO: Реализовать загрузку всех параметров конфигурации
	// Использовать вспомогательные функции getEnv и getEnvAsInt
	// Установить разумные значения по умолчанию
	return &Config{
		ServerHost:      getEnv("SERVER_HOST", "localhost"),
		ServerPort:      getEnvAsInt("SERVER_PORT", 8080),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnvAsInt("DB_PORT", 5432),
		DBUser:          getEnv("DB_USER", "user"),
		DBPassword:      getEnv("DB_PASSWORD", "password"),
		DBName:          getEnv("DB_NAME", "dbname"),
		DBSSLMode:       getEnv("DB_SSLMODE", "disable"),
		JWTSecret:       getEnv("JWT_SECRET", "supersecretkey"),
		JWTExpiryHours:  getEnvAsInt("JWT_EXPIRY_HOURS", 24),
		CacheTTLMinutes: getEnvAsInt("CACHE_TTL_MINUTES", 60),
	}

}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt получает значение переменной окружения как int или возвращает значение по умолчанию
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
