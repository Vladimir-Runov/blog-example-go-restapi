package handler

// https://github.com/Vladimir-Runov/blog-example-go-restapi

import (
	"blog-example-go-restapi/internal/model"
	"blog-example-go-restapi/internal/service"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

var (
	ErrPostNotFound = errors.New("post not found")
	ErrForbidden    = errors.New("forbidden")
)

// Create обрабатывает создание нового поста
// POST /api/posts
// Требует аутентификации
func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: Реализовать создание поста
	// Шаги:
	// 1. Проверить метод запроса (должен быть POST)
	// 2. Получить userID из контекста (установлен middleware)
	// 3. Декодировать JSON тело в PostCreateRequest
	// 4. Создать пост через postService.Create
	// 5. Вернуть созданный пост как JSON (201 Created)
	//http.Error(w, "Not implemented", http.StatusNotImplemented)

	// 1. Проверка метода запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Получение userID из контекста
	userIDstr, ok := getUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 3. Декодирование JSON тела в структуру PostCreateRequest
	var req model.PostCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 4. Создание поста через postService.Create
	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	post, err := h.postService.Create(r.Context(), userID, &req)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	// 5. Отправка созданного поста как JSON с статусом 201
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(post); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GetByID возвращает пост по ID
// GET /api/posts/{id}
// Не требует аутентификации
func (h *PostHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// TODO: Реализовать получение поста по ID
	// Шаги:
	// 1. Проверить метод запроса (должен быть GET)
	// 2. Извлечь ID из URL пути
	// 3. Получить пост через postService.GetByID
	// 4. Обработать ошибки (ErrPostNotFound -> 404)
	// 5. Вернуть пост как JSON (200 OK)
	// http.Error(w, "Not implemented", http.StatusNotImplemented)

	// 1. Проверка метода запроса
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// : /api/posts/{id}
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	idStr := pathParts[len(pathParts)-1]

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	//	vars := mux.Vars(r)
	//	idStr, ok := vars["id"]
	//	if !ok {
	//		http.Error(w, "Bad Request: missing id", http.StatusBadRequest)
	//		return
	//	}
	// преобразование строки в число
	//id, err := strconv.Atoi(idStr)
	//if err != nil || id <= 0 {
	//	http.Error(w, "Bad Request: invalid id", http.StatusBadRequest)
	//	return
	//}

	// 3. Получить пост через postService.GetByID
	post, err := h.postService.GetByID(r.Context(), id)
	if err != nil {
		// 4. Обработка ошибок
		if err == ErrPostNotFound {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 5. Вернуть пост как JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(post); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// GetAll возвращает список постов с пагинацией
// GET /api/posts?limit=10&offset=0
// Не требует аутентификации
func (h *PostHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// TODO: Реализовать получение списка постов
	// Шаги:
	// 1. Проверить метод запроса (должен быть GET)
	// 2. Извлечь параметры пагинации из query string
	// 3. Получить посты через postService.GetAll
	// 4. Создать ответ с метаданными пагинации
	// 5. Вернуть список постов как JSON (200 OK)
	// http.Error(w, "Not implemented", http.StatusNotImplemented)

	// 1. Проверяем метод
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Извлекаем параметры из query string
	query := r.URL.Query()

	limitStr := query.Get("limit")
	offsetStr := query.Get("offset")

	// Значения по умолчанию
	limit := 10
	offset := 0

	// Парсим limit
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	// Парсим offset
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// 3. Получаем посты через сервис
	posts, total, err := h.postService.GetAll(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	if posts == nil {
	}

	// 4. Создаем ответ с метаданными пагинации
	response := struct {
		Items      []model.Post `json:"items"`
		TotalCount int          `json:"total_count"`
		Limit      int          `json:"limit"`
		Offset     int          `json:"offset"`
	}{
		//		Items:      posts, todo:
		TotalCount: total,
		Limit:      limit,
		Offset:     offset,
	}

	// 5. Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// Update обновляет пост
// PUT /api/posts/{id}
// Требует аутентификации, может обновить только автор
func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	// TODO: Реализовать обновление поста
	// Шаги:
	// 1. Проверить метод запроса (должен быть PUT)
	// 2. Получить userID из контекста
	// 3. Извлечь ID поста из URL
	// 4. Декодировать JSON тело в PostUpdateRequest
	// 5. Обновить через postService.Update
	// 6. Обработать ошибки (404 для не найден, 403 для чужого поста)
	// 7. Вернуть обновленный пост как JSON (200 OK)
	// http.Error(w, "Not implemented", http.StatusNotImplemented)

	// 1. Проверяем метод
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Получаем userID из контекста (только авторы могут обновлять)
	userIDstr, ok := getUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	// 3. Извлекаем ID поста из URL
	//vars := mux.Vars(r)
	//idStr, exists := vars["id"]
	//if !exists {
	//	http.Error(w, "Bad Request: missing post ID", http.StatusBadRequest)
	//	return
	//}
	//postID, err := strconv.Atoi(idStr)
	//if err != nil || postID <= 0 {
	postIDStr := r.URL.Path[len("/api/posts/"):] // Извлечение ID поста из URL
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}

	// 4. Декодируем JSON тело в PostUpdateRequest
	var req model.PostUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// 5. Обновляем через postService.Update
	updatedPost, err := h.postService.Update(r.Context(), postID, userID, &req)
	if err != nil {
		if errors.Is(err, service.ErrPostNotFound) {
			http.Error(w, "Post Not Found", http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 6. Возвращаем обновленный пост как JSON (200 OK)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedPost)
}

// Delete удаляет пост
// DELETE /api/posts/{id}
// Требует аутентификации, может удалить только автор
func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// TODO: Реализовать удаление поста
	// Шаги:
	// 1. Проверить метод запроса (должен быть DELETE)
	// 2. Получить userID из контекста
	// 3. Извлечь ID поста из URL
	// 4. Удалить через postService.Delete
	// 5. Обработать ошибки (404 для не найден, 403 для чужого поста)
	// 6. Вернуть 204 No Content при успехе

	// http.Error(w, "Not implemented", http.StatusNotImplemented)
	// 1. Проверка метода
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Получение userID из контекста
	userIDstr, ok := getUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	// 3. Извлечение ID поста из URL
	//vars := mux.Vars(r)
	//idStr, exists := vars["id"]
	postIDStr := r.URL.Path[len("/api/posts/"):] // Извлечение ID поста из URL
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}
	//	if !exists {
	//		http.Error(w, "Bad Request: missing post ID", http.StatusBadRequest)
	//		return
	//	}
	//	postID, err := strconv.Atoi(idStr)
	//	if err != nil || postID <= 0 {
	//		http.Error(w, "Invalid post ID", http.StatusBadRequest)
	//		return
	//	}

	// 4. Удаление поста через сервис
	err = h.postService.Delete(r.Context(), postID, userID)
	if err != nil {
		// 5. Обработка ошибок
		if err == ErrPostNotFound {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		if err == ErrForbidden {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 6. Успешное удаление возвращает 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

// GetByAuthor возвращает посты конкретного автора
// GET /api/posts/author/{authorID}?limit=10&offset=0
// Не требует аутентификации
func (h *PostHandler) GetByAuthor(w http.ResponseWriter, r *http.Request) {
	// TODO: Реализовать получение постов автора
	// Шаги:
	// 1. Проверить метод запроса (должен быть GET)
	// 2. Извлечь ID автора из URL
	// 3. Извлечь параметры пагинации из query string
	// 4. Получить посты через postService.GetByAuthor
	// 5. Создать ответ с метаданными и списком постов
	// 6. Вернуть как JSON (200 OK)

	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// extractIDFromPath извлекает ID из пути URL
func extractIDFromPath(path, prefix string) string {
	// TODO: Реализовать извлечение ID из пути
	// Пример: path = "/api/posts/123", prefix = "/api/posts/"
	// Должен вернуть "123"
	// return ""
	// Убираем префикс из начала пути
	if len(path) <= len(prefix) || path[:len(prefix)] != prefix {
		return ""
	}
	idPart := path[len(prefix):]

	// Можно дополнительно очистить строку (например, убрать слеши)
	idPart = strings.Trim(idPart, "/")
	return idPart
}

// getUserIDFromContext извлекает ID пользователя из контекста
func getUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value("userID").(string) // Предполагается, что userID хранится как строка в контексте
	return userID, ok
}
