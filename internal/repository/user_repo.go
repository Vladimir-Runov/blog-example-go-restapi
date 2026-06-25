package repository

import (
	"blog-api/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
)

// UserRepo представляет репозиторий для работы с пользователями
type UserRepo struct {
	db *sql.DB
}

// NewUserRepo создает новый репозиторий пользователей
func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

// Create создает нового пользователя
func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	// TODO: Реализовать создание пользователя
	// 1. Подготовить SQL запрос INSERT INTO users...
	// 2. Установить created_at и updated_at = time.Now()
	// 3. Выполнить запрос и получить ID созданной записи
	// 4. Установить ID в структуру user
	//
	// HINT: Используйте QueryRowContext с RETURNING id для получения ID
	// Пример запроса:
	// INSERT INTO users (username, email, password, created_at, updated_at)
	// VALUES ($1, $2, $3, $4, $5) RETURNING id

	query := `
		INSERT INTO users (username, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// TODO: Выполнить запрос и обработать результат
	//return fmt.Errorf("not implemented")

	// Execute the query and get the ID of the newly created user
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no rows were returned")
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID получает пользователя по ID
func (r *UserRepo) GetByID(ctx context.Context, id int) (*model.User, error) {
	// TODO: Реализовать получение пользователя по ID
	// 1. Подготовить SQL запрос SELECT ... FROM users WHERE id = $1
	// 2. Выполнить запрос
	// 3. Просканировать результат в структуру User
	// 4. Обработать случай, когда пользователь не найден (sql.ErrNoRows)
	//
	// HINT: Используйте QueryRowContext и Scan

	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user model.User
	// TODO: Выполнить запрос и просканировать результат
	// Не забудьте обработать sql.ErrNoRows и вернуть ErrUserNotFound

	//	_ = query // Удалите эту строку после реализации
	//	return nil, fmt.Errorf("not implemented")

	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err) // Обработка случая, когда пользователь не найден
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil // Возвращаем найденного пользователя
}

// GetByEmail получает пользователя по email
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	// TODO: Реализовать получение пользователя по email
	// Аналогично GetByID, но поиск по полю email
	//return nil, fmt.Errorf("not implemented")
	query := `
        SELECT id, username, email, password, created_at, updated_at
        FROM users
        WHERE email = $1
    `

	var user model.User

	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err) // Обработка случая, когда пользователь не найден
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil // Возвращаем найденного пользователя
}

// GetByUsername получает пользователя по username
func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	// TODO: Реализовать получение пользователя по username
	// Аналогично GetByID, но поиск по полю username
	//return nil, fmt.Errorf("not implemented")
	query := `
        SELECT id, username, email, password, created_at, updated_at
        FROM users
        WHERE username = $1
    `

	var user model.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err) // Обработка случая, когда пользователь не найден
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil // Возвращаем найденного пользователя
}

// ExistsByEmail проверяет существование пользователя по email
func (r *UserRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	// TODO: Реализовать проверку существования пользователя
	// HINT: Используйте SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	// TODO: Выполнить запрос и просканировать результат в переменную exists

	//_ = query // Удалите эту строку после реализации
	//return false, fmt.Errorf("not implemented")

	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check existence of user: %w", err)
	}

	return exists, nil // Возвращаем результат проверки
}

// ExistsByUsername проверяет существование пользователя по username
func (r *UserRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	// TODO: Реализовать проверку существования пользователя по username
	// Аналогично ExistsByEmail

	//return false, fmt.Errorf("not implemented")
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check existence of user: %w", err)
	}

	return exists, nil // Возвращаем результат проверки
}

// Update обновляет данные пользователя
func (r *UserRepo) Update(ctx context.Context, user *model.User) error {
	// TODO: (Опционально) Реализовать обновление пользователя
	// 1. Подготовить SQL запрос UPDATE users SET ... WHERE id = $X
	// 2. Обновить updated_at = time.Now()
	// 3. Выполнить запрос
	// 4. Проверить, что запись была обновлена (RowsAffected)

	//return fmt.Errorf("not implemented")
	// Подготовить SQL запрос для обновления данных пользователя

	query := `
        UPDATE users 
        SET username = $1, email = $2, updated_at = $3 
        WHERE id = $4
    `

	// Установим текущее время для updated_at
	updatedAt := time.Now()

	// Выполнить запрос
	res, err := r.db.ExecContext(ctx, query, user.Username, user.Email, updatedAt, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Проверить, что запись была обновлена
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows were updated")
	}

	return nil // Успешное обновление
}

// Delete удаляет пользователя
func (r *UserRepo) Delete(ctx context.Context, id int) error {
	// TODO: (Опционально) Реализовать удаление пользователя
	// 1. Подготовить SQL запрос DELETE FROM users WHERE id = $1
	// 2. Выполнить запрос
	// 3. Проверить, что запись была удалена (RowsAffected)

	//return fmt.Errorf("not implemented")
	// Подготовить SQL запрос для удаления пользователя
	query := `DELETE FROM users WHERE id = $1`

	// Выполнить запрос
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Проверить, что запись была удалена
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows were deleted")
	}

	return nil // Успешное удаление
}
