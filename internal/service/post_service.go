package service

import (
	"blog-example-go-restapi/internal/model"
	"blog-example-go-restapi/internal/repository"
	"blog-example-go-restapi/pkg/auth"
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrPostNotFound = errors.New("post not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

type PostService struct {
	postRepo   repository.PostRepository
	userRepo   repository.UserRepository
	jwtManager *auth.JWTManager
}

func NewPostService(postRepo repository.PostRepository, userRepo repository.UserRepository, jwtManager *auth.JWTManager) *PostService {
	return &PostService{
		postRepo:   postRepo,
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (s *PostService) Create(ctx context.Context, userID int, req *model.PostCreateRequest) (*model.Post, error) {
	// TODO: Создать новый пост
	// Шаги:
	// 1. Валидация данных (title не пустой и <= 200 символов, content не пустой)
	// 2. Создать модель поста с данными из запроса и userID
	// 3. Сохранить через репозиторий
	// 4. Вернуть созданный пост

	if req.Title == "" || len(req.Title) > 200 {
		return nil, fmt.Errorf("title must be non-empty and less than or equal to 200 characters")
	}
	if req.Content == "" {
		return nil, fmt.Errorf("content must be non-empty")
	}

	post := &model.Post{
		Title:     req.Title,
		Content:   req.Content,
		AuthorID:  userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Step 3: Save through repository
	err := s.postRepo.Create(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("failed to save post: %w", err)
	}

	return post, nil
	//return nil, fmt.Errorf("not implemented")
}

func (s *PostService) GetByID(ctx context.Context, id int) (*model.Post, error) {
	// TODO: Получить пост по ID
	// Шаги:
	// 1. Получить пост через репозиторий
	// 2. Опционально: загрузить информацию об авторе
	// 3. Вернуть пост

	post, err := s.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get post by ID: %w", err)
	}

	if post.AuthorID != 0 {
		author, err := s.userRepo.GetByID(ctx, post.AuthorID)
		if err == nil {
			post.AuthorID = author.ID
		}
	}

	return post, nil
	//return nil, fmt.Errorf("not implemented")
}

func (s *PostService) GetAll(ctx context.Context, limit, offset int) ([]*model.Post, int, error) {
	// TODO: Получить список постов с пагинацией
	// Шаги:
	// 1. Валидировать и нормализовать параметры пагинации (limit по умолчанию 10, максимум 100)
	// 2. Получить посты через репозиторий
	// 3. Получить общее количество для пагинации
	// 4. Опционально: обогатить данные информацией об авторах
	// 5. Вернуть посты и общее количество

	if limit <= 0 {
		limit = 10 // default limit
	} else if limit > 100 {
		limit = 100 // maximum limit
	}

	if offset < 0 {
		offset = 0 // ensure offset is not negative
	}

	// Step 2: Get posts from the repository
	posts, err := s.postRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get posts: %w", err)
	}

	// Step 3: Get total count for pagination
	totalCount, err := s.postRepo.GetTotalCount(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total post count: %w", err)
	}

	for _, post := range posts {
		if post.AuthorID != 0 {
			author, err := s.userRepo.GetByID(ctx, post.AuthorID)
			if err == nil {
				post.AuthorID = author.ID
			}
		}
	}

	return posts, totalCount, nil
	//return nil, 0, fmt.Errorf("not implemented")
}

func (s *PostService) Update(ctx context.Context, id int, userID int, req *model.PostUpdateRequest) (*model.Post, error) {
	// TODO: Обновить пост
	// Шаги:
	// 1. Получить существующий пост
	// 2. Проверить что userID является автором (иначе ErrForbidden)
	// 3. Валидировать новые данные (если предоставлены)
	// 4. Обновить только измененные поля
	// 5. Сохранить через репозиторий
	// 6. Вернуть обновленный пост

	post, err := s.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	if post == nil {
		return nil, fmt.Errorf("post not found")
	}

	// Step 2: Check that userID is the author
	if post.AuthorID != userID {
		return nil, fmt.Errorf("forbidden: you are not the author of this post")
	}

	// Step 3: Validate new data (if provided)
	if err := validatePostUpdateRequest(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Step 4: Update only changed fields
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}

	err = s.postRepo.Update(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	return post, nil
	//return nil, fmt.Errorf("not implemented")
}

func (s *PostService) Delete(ctx context.Context, id int, userID int) error {
	// TODO: Удалить пост
	// Шаги:
	// 1. Найти пост и проверить существование
	// 2. Проверить что userID является автором
	// 3. Удалить через репозиторий
	// 4. Вернуть соответствующую ошибку при неудаче

	// Step 1: Find post and check existence
	post, err := s.postRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get post: %w", err)
	}
	if post == nil {
		return fmt.Errorf("post not found")
	}

	// Step 2: Check that userID is the author
	if post.AuthorID != userID {
		return fmt.Errorf("forbidden: you are not the author of this post")
	}

	// Step 3: Delete through repository
	if err := s.postRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	// Step 4: Return nil if successful
	return nil
	//return fmt.Errorf("not implemented")
}

func (s *PostService) GetByAuthor(ctx context.Context, authorID int, limit, offset int) ([]*model.Post, int, error) {
	// TODO: Получить посты конкретного автора
	// Шаги:
	// 1. Валидировать параметры пагинации
	// 2. Получить посты автора через репозиторий
	// 3. Получить общее количество постов автора
	// 4. Опционально: добавить информацию об авторе к постам
	// 5. Вернуть результат с общим количеством

	if limit <= 0 {
		return nil, 0, fmt.Errorf("invalid limit: must be greater than 0")
	}
	if offset < 0 {
		return nil, 0, fmt.Errorf("invalid offset: must be non-negative")
	}

	posts, err := s.postRepo.GetByAuthorID(ctx, authorID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get posts by author: %w", err)
	}

	totalPosts, err := s.postRepo.CountByAuthorID(ctx, authorID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posts by author: %w", err)
	}

	//	for _, post := range posts {
	//		post.AuthorID =
	//	}

	return posts, totalPosts, nil
	//return nil, 0, fmt.Errorf("not implemented")
}

// validatePostCreateRequest проверяет корректность данных для создания поста
func validatePostCreateRequest(req *model.PostCreateRequest) error {
	// TODO: Реализовать валидацию title и content

	if req.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}

	// Check if content is empty
	if req.Content == "" {
		return fmt.Errorf("content cannot be empty")
	}

	const maxTitleLength = 100
	const maxContentLength = 5000

	if len(req.Title) > maxTitleLength {
		return fmt.Errorf("title cannot exceed %d characters", maxTitleLength)
	}

	if len(req.Content) > maxContentLength {
		return fmt.Errorf("content cannot exceed %d characters", maxContentLength)
	}

	return nil
}

// validatePostUpdateRequest проверяет корректность данных для обновления поста
func validatePostUpdateRequest(req *model.PostUpdateRequest) error {
	// TODO: Реализовать валидацию опциональных полей
	// Validate Title if provided

	if req.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}

	const maxTitleLength = 100
	if len(req.Title) > maxTitleLength {
		return fmt.Errorf("title cannot exceed %d characters", maxTitleLength)
	}

	// Validate Content if provided

	if req.Content == "" {
		return fmt.Errorf("content cannot be empty")
	}

	const maxContentLength = 5000
	if len(req.Content) > maxContentLength {
		return fmt.Errorf("content cannot exceed %d characters", maxContentLength)
	}

	return nil
}
