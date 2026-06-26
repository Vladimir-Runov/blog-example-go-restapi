package handler

import (
	"blog-api/internal/model"
	"blog-api/internal/service"
	"context"
	"encoding/json"
	"net/http"
)

// UserCreateRequest представляет тело запроса на регистрацию
type UserCreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// UserLoginRequest представляет тело запроса на вход
type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthHandler struct {
	userService service.UserService
}

func NewAuthHandler(userService *service.UserService) *AuthHandler {
	return &AuthHandler{
		userService: *userService,
	}
}

// Register обрабатывает запрос на регистрацию нового пользователя
// POST /api/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// TODO: Реализовать обработку регистрации
	// Шаги:
	// 1. Проверить метод запроса (должен быть POST)
	// 2. Декодировать JSON тело в UserCreateRequest
	// 3. Вызвать userService.Register
	// 4. Обработать ошибки (ErrUserAlreadyExists -> 409 Conflict)
	// 5. Вернуть JSON ответ с токеном (201 Created)

	//http.Error(w, "Not implemented", http.StatusNotImplemented)
	// 1. Проверить метод запроса (должен быть POST)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Декодировать JSON тело в UserCreateRequest
	var req UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// 3. Вызвать userService.Register
	token, err := h.userService.Register(req.Username, req.Password, req.Email)

	// 4. Обработать ошибки
	if err != nil {
		if err == service.ErrUserAlreadyExists {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// 5. Вернуть JSON ответ с токеном (201 Created)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"token": token})

	//	response := map[string]string{"token": token}
	//	if err := json.NewEncoder(w).Encode(response); err != nil {
	//		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	//	}

}

// Login обрабатывает запрос на вход пользователя
// POST /api/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// TODO: Реализовать обработку входа
	// Шаги:
	// 1. Проверить метод запроса (должен быть POST)
	// 2. Декодировать JSON тело в UserLoginRequest
	// 3. Вызвать userService.Login
	// 4. Обработать ошибки (ErrInvalidCredentials -> 401 Unauthorized)
	// 5. Вернуть JSON ответ с токеном (200 OK)

	//http.Error(w, "Not implemented", http.StatusNotImplemented)

	// 1. Проверить метод запроса (должен быть POST)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Декодировать JSON тело в UserLoginRequest
	var req UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// 3. Вызвать userService.Login
	token, err := h.userService.Login(req.Username, req.Password)

	// 4. Обработать ошибки
	if err != nil {
		if err == services.ErrInvalidCredentials {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// 5. Вернуть JSON ответ с токеном (200 OK)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// GetProfile возвращает профиль текущего пользователя (опционально)
// Этот метод не используется в эталонной реализации
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// TODO: Опционально - реализовать получение профиля
	// Этот эндпоинт не обязателен для базовой реализации
	// http.Error(w, "Not implemented", http.StatusNotImplemented)

	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userRepo.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to fetch user profile", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Подготовка ответа - исключаем пароль и другие чувствительные данные
	response := &model.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	// Отправляем ответ в формате JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// writeError отправляет JSON ответ с ошибкой
func writeError(w http.ResponseWriter, message string, statusCode int) {
	// TODO: Реализовать отправку ошибки в формате JSON
	// Создать структуру ErrorResponse и отправить как JSON
	//http.Error(w, message, statusCode)

	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(statusCode)
	//json.NewEncoder(w).Encode(ErrorResponse{Message: message})
	response := map[string]any{
		"error":   message,
		"status":  statusCode,
		"message": message,
	}
	WriteJSON(w, response, statusCode)
}

func WriteJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Log error but can't really recover at this point
		//logger := loggingutil.Get(context.Background())
		//logger.Error("Error encoding JSON response", "error", err)
	}
}

// getUserIDFromContext извлекает ID пользователя из контекста
func getUserIDFromContext(ctx context.Context) (int, bool) {
	// TODO: Извлечь userID из контекста
	// Ключ устанавливается в auth middleware
	// return 0, false
	const userIDKey = "userID"

	// Получаем значение из контекста
	val := ctx.Value(userIDKey)
	if val == nil {
		return 0, false
	}

	// Приводим к типу int, проверяем
	userID, ok := val.(int)
	return userID, ok
}
