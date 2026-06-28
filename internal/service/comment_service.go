package service

import (
	"blog-api/internal/model"
	"blog-api/internal/repository"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrCommentNotFound = errors.New("comment not found")
	ErrPostNotExists   = errors.New("post does not exist")
)

type CommentService struct {
	commentRepo repository.CommentRepository
	postRepo    repository.PostRepository
	userRepo    repository.UserRepository
}

func NewCommentService(
	commentRepo repository.CommentRepository,
	postRepo repository.PostRepository,
	userRepo repository.UserRepository,
) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
		userRepo:    userRepo,
	}
}

func (s *CommentService) Create(ctx context.Context, userID int, req *model.CommentCreateRequest) (*model.Comment, error) {
	// TODO: Создать новый комментарий
	// Шаги:
	// 1. Валидация данных (content не пустой и <= 1000 символов)
	// 2. Проверить что пост существует
	// 3. Создать модель комментария с userID как автором
	// 4. Сохранить через репозиторий
	// 5. Опционально: обогатить ответ информацией об авторе
	// 6. Вернуть созданный комментарий

	// 1. Валидация данных
	if req.Content == "" || len(req.Content) > 1000 {
		return nil, fmt.Errorf("invalid content: %w", errors.New("content must be non-empty and less than or equal to 1000 characters"))
	}

	// 2. Проверить, что пост существует
	exists, err := s.postRepo.Exists(ctx, req.PostID)
	if err != nil {
		return nil, err // Возвращаем ошибку, если не удалось проверить существование поста
	}
	if !exists {
		return nil, errors.New("post not found")
	}

	//	post, err := s.postRepo.GetByID(ctx, req.PostID)
	//	if err != nil {
	//		if errors.Is(err, repository.ErrPostNotFound) {
	//			return nil, ErrPostNotExists
	//		}
	//		return nil, fmt.Errorf("failed to check post existence: %w", err)
	//	}

	// 3. Создание модели комментария
	comment := &model.Comment{
		AuthorID:  userID,
		PostID:    req.PostID,
		Content:   req.Content,
		CreatedAt: time.Now(), // добавляем время создания
	}

	// 4. Сохранение через репозиторий
	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("failed to save comment: %w", err)
	}

	// 5. Опционально: обогащение ответа информацией об авторе
	author, err := s.userRepo.GetByID(ctx, userID)
	if err == nil {
		comment.AuthorID = author.ID
	}

	// 6. Возврат созданного комментария
	return comment, nil
	//return nil, fmt.Errorf("not implemented")
}

func (s *CommentService) GetByID(ctx context.Context, id int) (*model.Comment, error) {
	// TODO: Получить комментарий по ID
	// Шаги:
	// 1. Получить комментарий через репозиторий
	// 2. Опционально: добавить информацию об авторе
	// 3. Вернуть результат или ErrCommentNotFound

	// Step 1: Get the comment from the repository
	comment, err := s.commentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error retrieving comment: %w", err)
	}

	// Step 2: Check if the comment exists
	if comment == nil {
		return nil, ErrCommentNotFound // Return a specific error if the comment is not found
	}

	// Step 3: Optionally add author information (if required)
	author, err := s.userRepo.FindByID(ctx, comment.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving author: %w", err)
	}

	// Assuming Comment has an Author field to hold the author's information
	comment.Author = author

	// Step 4: Return the result
	return comment, nil
	// return nil, fmt.Errorf("not implemented")
}

const (
	defaultLimit = 20
	maxLimit     = 100
)

func (s *CommentService) GetByPost(ctx context.Context, postID int, limit, offset int) ([]*model.Comment, int, error) {
	// TODO: Получить комментарии к посту с пагинацией
	// Шаги:
	// 1. Валидировать параметры пагинации (limit по умолчанию 20, максимум 100)
	// 2. Опционально: проверить существование поста
	// 3. Получить комментарии через репозиторий
	// 4. Получить общее количество для пагинации
	// 5. Опционально: обогатить данные информацией об авторах
	// 6. Вернуть комментарии и общее количество
	// 1. Валидировать параметры пагинации
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	// 2. Опционально: проверить существование поста
	// Если у вас есть метод проверки существования поста, вы можете его вызвать здесь
	// if !s.postExists(ctx, postID) {
	//     return nil, 0, errors.New("post not found")
	// }

	// 3. Получить комментарии через репозиторий
	comments, err := s.commentRepo.GetByPostID(ctx, postID, limit, offset)
	if err != nil {
		return nil, 0, err // Возвращаем ошибку, если не удалось получить комментарии
	}

	// 4. Получить общее количество для пагинации
	totalCount, err := s.commentRepo.GetCountByPostID(ctx, postID)
	if err != nil {
		return nil, 0, err // Возвращаем ошибку, если не удалось получить общее количество
	}

	// 5. Опционально: обогатить данные информацией об авторах
	// Например, можно добавить информацию о пользователях в комментарии

	// 6. Вернуть комментарии и общее количество
	return comments, totalCount, nil //return nil, 0, fmt.Errorf("not implemented")
}

func (s *CommentService) Update(ctx context.Context, id int, userID int, req *model.CommentUpdateRequest) (*model.Comment, error) {
	// TODO: Обновить комментарий
	// Шаги:
	// 1. Найти существующий комментарий
	// 2. Проверить что userID является автором (иначе ErrForbidden)
	// 3. Валидировать новый content
	// 4. Обновить content и временную метку
	// 5. Сохранить через репозиторий
	// 6. Опционально: добавить информацию об авторе
	// 7. Вернуть обновленный комментарий
	// 1. Найти существующий комментарий
	comment, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err // Возвращаем ошибку, если комментарий не найден
	}

	// 2. Проверить, что userID является автором
	if comment.AuthorID != userID {
		return nil, errors.New("forbidden: user is not the author of the comment")
	}

	// 3. Валидировать новый content
	if req.Content == "" {
		return nil, errors.New("content cannot be empty")
	}

	// 4. Обновить content и временную метку
	comment.Content = req.Content
	comment.UpdatedAt = time.Now()

	// 5. Сохранить через репозиторий
	updatedComment, err := s.repo.Update(ctx, comment)
	if err != nil {
		return nil, err // Возвращаем ошибку, если обновление не удалось
	}

	// 6. Опционально: добавить информацию об авторе (если требуется)
	// Например, можно добавить данные о пользователе в комментарий

	// 7. Вернуть обновленный комментарий
	return updatedComment, nil //return nil, fmt.Errorf("not implemented")
}

func (s *CommentService) Delete(ctx context.Context, id int, userID int) error {
	// TODO: Удалить комментарий
	// Шаги:
	// 1. Найти комментарий и проверить существование
	// 2. Проверить что userID является автором
	// 3. Удалить через репозиторий
	// 4. Вернуть соответствующую ошибку при неудаче
	// 1. Найти комментарий и проверить существование
	comment, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err // Возвращаем ошибку, если комментарий не найден
	}

	// 2. Проверить, что userID является автором
	if comment.AuthorID != userID {
		return errors.New("user is not the author of the comment")
	}

	// 3. Удалить через репозиторий
	err = s.repo.Delete(ctx, id)
	if err != nil {
		return err // Возвращаем ошибку, если удаление не удалось
	}

	// 4. Возвращаем nil, если удаление прошло успешно
	return nil //return fmt.Errorf("not implemented")
}

func (s *CommentService) GetByAuthor(ctx context.Context, authorID int, limit, offset int) ([]*model.Comment, int, error) {
	// TODO: Получить комментарии конкретного автора
	// Шаги:
	// 1. Валидировать параметры пагинации
	// 2. Получить комментарии автора через репозиторий
	// 3. Получить общее количество комментариев автора
	// 4. Опционально: добавить информацию об авторе
	// 5. Вернуть результат с общим количеством
	// 1. Валидировать параметры пагинации
	if limit <= 0 {
		return nil, 0, errors.New("limit must be greater than zero")
	}
	if offset < 0 {
		return nil, 0, errors.New("offset cannot be negative")
	}

	// 2. Получить комментарии автора через репозиторий
	comments, err := s.repo.GetCommentsByAuthor(ctx, authorID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// 3. Получить общее количество комментариев автора
	totalCount, err := s.repo.GetTotalCommentsByAuthor(ctx, authorID)
	if err != nil {
		return nil, 0, err
	}

	// 4. Опционально: добавить информацию об авторе (если необходимо)

	// 5. Вернуть результат с общим количеством
	return comments, totalCount, nil //return nil, 0, fmt.Errorf("not implemented")
}

// validateCommentCreateRequest проверяет корректность данных для создания комментария
func validateCommentCreateRequest(req *model.CommentCreateRequest) error {
	// TODO: Реализовать валидацию content и PostID
	// Проверка на пустое содержимое
	if strings.TrimSpace(req.Content) == "" {
		return errors.New("content cannot be empty")
	}

	// Проверка длины содержимого
	const maxContentLength = 500 // максимальная длина комментария
	if len(req.Content) > maxContentLength {
		return errors.New("content exceeds maximum length")
	}

	// Проверка на наличие PostID
	if req.PostID <= 0 {
		return errors.New("invalid PostID")
	}

	return nil
}

// validateCommentUpdateRequest проверяет корректность данных для обновления комментария
func validateCommentUpdateRequest(req *model.CommentUpdateRequest) error {
	// TODO: Реализовать валидацию content
	// Проверка на пустое содержимое
	if strings.TrimSpace(req.Content) == "" {
		return errors.New("content cannot be empty")
	}

	// Проверка длины содержимого
	const maxContentLength = 500 // максимальная длина комментария
	if len(req.Content) > maxContentLength {
		return errors.New("content exceeds maximum length")
	}

	// Здесь можно добавить дополнительные проверки, если необходимо

	return nil
}
