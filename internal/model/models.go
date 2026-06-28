package model

import (
	"errors"
	"time"
)

// ErrUserAlreadyExists - ошибка, возникающая при попытке зарегистрировать пользователя, который уже существует.
var ErrUserAlreadyExists = errors.New("user already exists")

// User представляет модель пользователя в системе
type User struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"` // Хешированный пароль, не отдаем в JSON
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserResponse - структура для ответа с данными пользователя (без пароля)
// Поля: ID, Username, Email, CreatedAt
type UserResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// User.ToResponse() UserResponse - преобразует User в UserResponse
// ToResponse преобразует User в UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}

// Post представляет модель поста в блоге
type Post struct {
	ID        int       `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	AuthorID  int       `json:"author_id" db:"author_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CanBeEditedBy проверяет, может ли пользователь редактировать пост
func (p *Post) CanBeEditedBy(userID int) bool {
	return p.AuthorID == userID
}

// CanBeDeletedBy проверяет, может ли пользователь удалить пост если он его
func (p *Post) CanBeDeletedBy(userID int) bool {
	return p.AuthorID == userID
}

// Comment представляет модель комментария к посту
type Comment struct {
	ID        int       `json:"id" db:"id"`
	Content   string    `json:"content" db:"content"`
	PostID    int       `json:"post_id" db:"post_id"`
	AuthorID  int       `json:"author_id" db:"author_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CanBeEditedBy проверяет, может ли пользователь редактировать комментарий
// Post.CanBeEditedBy(userID int) bool - проверяет, может ли пользователь редактировать пост
// HINT: Пользователь может редактировать/удалять только свои  комментарии
func (c *Comment) CanBeEditedBy(userID int) bool {
	return c.AuthorID == userID
}

// CanBeDeletedBy проверяет, может ли пользователь удалить комментарий
// Post.CanBeDeletedBy(userID int) bool - проверяет, может ли пользователь удалить пост
func (c *Comment) CanBeDeletedBy(userID int) bool {
	return c.AuthorID == userID
}

// Response
//
// PostResponse - структура для ответа с данными поста
// Поля: ID, Title, Content, Author (UserResponse), CreatedAt, UpdatedAt
type PostResponse struct {
	ID        int          `json:"id"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	Author    UserResponse `json:"author"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// CommentResponse - структура для ответа с данными комментария
// Поля: ID, Content, PostID, Author (UserResponse), CreatedAt, UpdatedAt
type CommentResponse struct {
	ID        int          `json:"id"`
	Content   string       `json:"content"`
	PostID    int          `json:"post_id"`
	Author    UserResponse `json:"author"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// TokenResponse - структура для ответа с JWT токеном
// Поля: Token (string), ExpiresAt (time.Time), User (UserResponse)
type TokenResponse struct {
	Token     string       `json:"token"`
	ExpiresAt time.Time    `json:"expires_at"`
	User      UserResponse `json:"user"`
}

// Request
//
// UserCreateRequest представляет запрос на создание пользователя
type UserCreateRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// UserLoginRequest представляет запрос на вход пользователя
type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// PostCreateRequest представляет запрос на создание поста
type PostCreateRequest struct {
	Title   string `json:"title" validate:"required,min=1,max=200"`
	Content string `json:"content" validate:"required,min=1"`
}

// PostUpdateRequest представляет запрос на обновление поста
type PostUpdateRequest struct {
	Title   string `json:"title" validate:"required,min=1,max=200"`
	Content string `json:"content" validate:"required,min=1"`
}

// CommentCreateRequest представляет запрос на создание комментария
type CommentCreateRequest struct {
	AuthorID int    `json:"author_id" db:"author_id"`
	PostID   int    `json:"post_id" db:"author_id"`
	Content  string `json:"content" validate:"required,min=1,max=1000"`
}

// CommentUpdateRequest представляет запрос на обновление комментария
type CommentUpdateRequest struct {
	CommentID int    `json:"commentId"`
	Content   string `json:"content" validate:"required,min=1"`
}
