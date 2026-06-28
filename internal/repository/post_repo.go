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
	ErrPostNotFound = errors.New("post not found")
)

// PostRepo представляет репозиторий для работы с постами
type PostRepo struct {
	db *sql.DB
}

// NewPostRepo создает новый репозиторий постов
func NewPostRepo(db *sql.DB) *PostRepo {
	return &PostRepo{db: db}
}

// Create создает новый пост
func (r *PostRepo) Create(ctx context.Context, post *model.Post) error {
	// TODO: Реализовать создание поста
	// 1. Подготовить SQL запрос INSERT INTO posts...
	// 2. Установить created_at и updated_at = time.Now()
	// 3. Выполнить запрос и получить ID созданной записи
	// 4. Установить ID в структуру post
	//
	// HINT: Используйте QueryRowContext с RETURNING id
	// Пример запроса:
	// INSERT INTO posts (title, content, author_id, created_at, updated_at)
	// VALUES ($1, $2, $3, $4, $5) RETURNING id

	// TODO: Выполнить запрос и обработать результат
	// Устанавливаем время создания и обновления
	now := time.Now()
	post.CreatedAt = now
	post.UpdatedAt = now

	// Подготавливаем SQL запрос
	query := `
        INSERT INTO posts (title, content, author_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `

	// Выполняем запрос и получаем ID созданной записи  err := r.db.QueryRowContext(ctx, query, ...).Scan(&post.ID)
	err := r.db.QueryRowContext(ctx, query, post.Title, post.Content, post.AuthorID, post.CreatedAt, post.UpdatedAt).Scan(&post.ID)
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	return nil
}

// GetByID получает пост по ID
func (r *PostRepo) GetByID(ctx context.Context, id int) (*model.Post, error) {
	// TODO: Реализовать получение поста по ID
	// 1. Подготовить SQL запрос SELECT ... FROM posts WHERE id = $1
	// 2. Выполнить запрос
	// 3. Просканировать результат в структуру Post
	// 4. Обработать случай sql.ErrNoRows -> вернуть ErrPostNotFound

	query := `
        SELECT id, title, content, author_id, created_at, updated_at
        FROM posts
        WHERE id = $1
    `

	var post model.Post
	// Выполняем запрос и просканируем результат в структуру Post
	err := r.db.QueryRowContext(ctx, query, id).Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrPostNotFound // Возвращаем ошибку, если пост не найден
		}
		return nil, fmt.Errorf("failed to get post: %w", err) // Обработка других ошибок
	}

	return &post, nil // Возвращаем найденный пост
}

// GetAll получает все посты с пагинацией
func (r *PostRepo) GetAll(ctx context.Context, limit, offset int) ([]*model.Post, error) {
	// TODO: Реализовать получение всех постов с пагинацией
	// 1. Подготовить SQL запрос с ORDER BY created_at DESC
	// 2. Добавить LIMIT и OFFSET для пагинации
	// 3. Выполнить запрос и получить rows
	// 4. Итерировать по rows и собрать массив постов
	// 5. Не забудьте закрыть rows (defer rows.Close())
	//
	// HINT: Используйте QueryContext для получения множества записей
	// Пример запроса:
	// SELECT id, title, content, author_id, created_at, updated_at
	// FROM posts
	// ORDER BY created_at DESC
	// LIMIT $1 OFFSET $2

	query := `
		SELECT id, title, content, author_id, created_at, updated_at
		FROM posts
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	// TODO: Выполнить запрос
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		// todo:
	}
	defer rows.Close()

	// TODO: Итерировать по результатам
	var posts []*model.Post
	for rows.Next() {
		var post model.Post

		err := rows.Scan(&post.ID, &post.Title)
		if err != nil {
			// todo:
		}
		posts = append(posts, &post)
	}

	//_ = query // Удалите эту строку после реализации
	return nil, fmt.Errorf("not implemented")
}

// GetTotalCount получает общее количество постов
func (r *PostRepo) GetTotalCount(ctx context.Context) (int, error) {
	// TODO: Реализовать подсчет общего количества постов
	// HINT: Используйте SELECT COUNT(*) FROM posts

	query := `SELECT COUNT(*) FROM posts`

	var count int
	// TODO: Выполнить запрос и получить количество
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get total count of posts: %w", err) // Обработка ошибки
	}

	return count, nil // Возвращаем общее количество постов
}

// Update обновляет пост
func (r *PostRepo) Update(ctx context.Context, post *model.Post) error {
	// TODO: Реализовать обновление поста
	// 1. Подготовить SQL запрос UPDATE posts SET ... WHERE id = $X
	// 2. Обновить только title, content и updated_at
	// 3. Выполнить запрос с помощью ExecContext
	// 4. Проверить RowsAffected - если 0, вернуть ErrPostNotFound
	//
	// HINT:
	// UPDATE posts
	// SET title = $1, content = $2, updated_at = $3
	// WHERE id = $4

	query := `
		UPDATE posts
		SET title = $1, content = $2, updated_at = $3
		WHERE id = $4
	`

	post.UpdatedAt = time.Now()

	// Выполняем запрос с помощью ExecContext
	result, err := r.db.ExecContext(ctx, query, post.Title, post.Content, post.UpdatedAt, post.ID)
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err) // Обработка ошибки
	}

	// Проверяем количество затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err) // Обработка ошибки
	}

	// Если ни одна строка не была затронута, возвращаем ошибку
	if rowsAffected == 0 {
		return ErrPostNotFound // Возвращаем ошибку, если пост не найден
	}

	return nil // Возвращаем nil, если обновление прошло успешно
}

// Delete удаляет пост
func (r *PostRepo) Delete(ctx context.Context, id int) error {
	// TODO: Реализовать удаление поста
	// 1. Подготовить SQL запрос DELETE FROM posts WHERE id = $1
	// 2. Выполнить запрос с помощью ExecContext
	// 3. Проверить RowsAffected - если 0, вернуть ErrPostNotFound

	query := `DELETE FROM posts WHERE id = $1`

	// Выполняем запрос с помощью ExecContext
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err) // Обработка ошибки
	}

	// Проверяем количество затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err) // Обработка ошибки
	}

	// Если ни одна строка не была затронута, возвращаем ошибку
	if rowsAffected == 0 {
		return ErrPostNotFound // Возвращаем ошибку, если пост не найден
	}

	return nil // Возвращаем nil, если удаление прошло успешно
}

// Exists проверяет существование поста
func (r *PostRepo) Exists(ctx context.Context, id int) (bool, error) {
	// TODO: Реализовать проверку существования поста
	// HINT: SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1)

	query := `SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1)`

	var exists bool
	// Выполняем запрос с помощью QueryRowContext
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check post existence: %w", err) // Обработка ошибки
	}

	return exists, nil // Возвращаем результат проверки
}

// GetByAuthorID получает посты определенного автора
func (r *PostRepo) GetByAuthorID(ctx context.Context, authorID int, limit, offset int) ([]*model.Post, error) {
	// TODO: (Опционально) Реализовать получение постов автора
	// Аналогично GetAll, но с дополнительным условием WHERE author_id = $X

	// Prepare the SQL query to fetch posts by author ID with pagination
	query := `
        SELECT id, title, content, author_id, created_at, updated_at
        FROM posts
        WHERE author_id = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `

	// Execute the query
	rows, err := r.db.QueryContext(ctx, query, authorID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying posts: %w", err)
	}
	defer rows.Close()

	// Slice to hold the retrieved posts
	var posts []*model.Post

	// Iterate over the result set
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning post: %w", err)
		}
		posts = append(posts, &post)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return posts, nil //return nil, fmt.Errorf("not implemented")
}
