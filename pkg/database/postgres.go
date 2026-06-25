package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Config содержит параметры подключения к PostgreSQL
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewPostgresDB создает новое подключение к PostgreSQL
func NewPostgresDB(cfg Config) (*sql.DB, error) {
	// TODO: Реализовать подключение к PostgreSQL
	// Шаги:
	// 1. Сформировать строку подключения (DSN) из параметров конфигурации
	// 2. Открыть соединение с БД используя sql.Open("postgres", dsn)
	// 3. Проверить соединение методом Ping()
	// 4. Настроить пул соединений (SetMaxOpenConns, SetMaxIdleConns)
	// 5. Вернуть подключение или ошибку
	//return nil, fmt.Errorf("not implemented")
	// Формируем DSN (строку подключения)

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	// Настраиваем пул соединений
	//db.SetMaxOpenConns(cfg.MaxOpenConns)
	//db.SetMaxIdleConns(cfg.MaxIdleConns)

	return db, nil
}

// Migrate выполняет миграции базы данных
func Migrate(db *sql.DB) error {
	// TODO: Реализовать применение миграций
	// Шаги:
	// 1. Создать таблицу users если не существует
	// 2. Создать таблицу posts если не существует
	// 3. Создать таблицу comments если не существует
	// 4. Создать необходимые индексы
	// 5. Вернуть ошибку если что-то пошло не так
	//
	// Структура таблиц:
	// - users: id, username, email, password_hash, created_at, updated_at
	// - posts: id, title, content, author_id, created_at, updated_at
	// - comments: id, content, post_id, author_id, created_at, updated_at

	queries := []string{
		// `CREATE TABLE IF NOT EXISTS users (...)`,
		`CREATE TABLE IF NOT EXISTS users (
											id SERIAL PRIMARY KEY,
											username VARCHAR(50) UNIQUE NOT NULL,
											email VARCHAR(255) UNIQUE NOT NULL,
											password VARCHAR(255) NOT NULL,
											created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
											updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP);`,

		// `CREATE TABLE IF NOT EXISTS posts (...)`,
		`CREATE TABLE IF NOT EXISTS posts (
											id SERIAL PRIMARY KEY,
											title VARCHAR(200) NOT NULL,
											content TEXT NOT NULL,
											author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
											created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
											updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,

		// `CREATE TABLE IF NOT EXISTS comments (...)`,
		`CREATE TABLE IF NOT EXISTS comments (
												id SERIAL PRIMARY KEY,
												content TEXT NOT NULL,
												post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
												author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
												created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,

		// `CREATE INDEX IF NOT EXISTS ...`, индекс на user_id в таблице posts
		`CREATE INDEX IF NOT EXISTS idx_posts_author_id ON posts(author_id);`,
		`CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);`,
		`CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at);`,
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Выполняем все запросы
	for _, query := range queries {
		if _, err := tx.Exec(query); err != nil {
			tx.Rollback()
			return fmt.Errorf("error executing query: %s, error: %w", query, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
	// TODO: Выполнить каждый запрос в транзакции
	//_ = queries // Удалить после реализации
	//return fmt.Errorf("not implemented")
}

// CheckConnection проверяет соединение с базой данных
func CheckConnection(db *sql.DB) error {
	// TODO: Реализовать проверку соединения
	// Использовать db.Ping() для проверки
	//return fmt.Errorf("not implemented")
	// Пингуем базу данных, чтобы проверить соединение
	if err := db.Ping(); err != nil {
		return fmt.Errorf("нет соединения с базой данных: %w", err)
	}
	return nil
}

// GetDSN формирует строку подключения к PostgreSQL
func GetDSN(cfg Config) string {
	// TODO: Сформировать DSN строку
	// Формат: "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s"
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
}

// Close закрывает соединение с базой данных
func Close(db *sql.DB) error {
	// TODO: Корректно закрыть соединение
	// return fmt.Errorf("not implemented")
	if err := db.Close(); err != nil {
		return fmt.Errorf("не удалось закрыть соединение: %w", err)
	}
	return nil
}

// TestConnection выполняет тестовый запрос к БД (опциональное задание)
func TestConnection(db *sql.DB) error {
	// TODO: Выполнить простой запрос для проверки работы БД
	// Например: SELECT 1
	//return fmt.Errorf("not implemented")
	var result int
	err := db.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		return fmt.Errorf("тестовая команда не выполнена: %w", err)
	}
	return nil
}
