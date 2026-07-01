package handler_test

import (
	"context"
	"log"
	"os"
	"strconv"
	"testing"

	"blog-example-go-restapi/internal/model"
	"blog-example-go-restapi/internal/repository"
	"blog-example-go-restapi/internal/service"
	"blog-example-go-restapi/pkg/auth"
	"blog-example-go-restapi/pkg/database"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_Register(t *testing.T) {
	// Загрузка конфигурации из .env файла тестового контура
	if err := godotenv.Load("../../.env_test"); err != nil {
		log.Printf("Warning: .env test file not found")
	}

	cfg := loadConfig()
	// Настройка конфигурации базы данных
	dbConfig := database.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
	}

	// Подключение к базе данных
	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		log.Fatalf("Could not connect to the test-database: %v", err)
	}
	defer db.Close()

	// Миграция базы данных
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Could not migrate the database: %v", err)
	}

	// Создание JWT менеджера и пользовательского репозитория
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiryHours)
	userRepo := repository.NewUserRepo(db)
	userService := service.NewUserService(userRepo, jwtManager)

	// Создание обработчика аутентификации
	///authHandler := handler.NewAuthHandler(userService)

	// Подготовка тестовых данных
	req := &model.UserCreateRequest{
		Username: "testuser",
		Email:    "testuser@example.com",
		Password: "password123",
	}

	// Вызов метода Register
	tokenResponse, err := userService.Register(context.Background(), req)

	// Проверка результатов
	assert.NoError(t, err)
	assert.NotNil(t, tokenResponse)
	assert.NotEmpty(t, tokenResponse.Token)

	// Дополнительно: проверить, что пользователь был успешно добавлен в базу данных
	user, err := userRepo.GetByEmail(context.Background(), req.Email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

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
		DBPort:          getEnvAsInt("DB_PORT", 5433),
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
