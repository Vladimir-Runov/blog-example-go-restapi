package service

import (
	"blog-example-go-restapi/internal/model"
	"blog-example-go-restapi/internal/repository"
	"blog-example-go-restapi/pkg/auth"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

type UserService struct {
	userRepo   repository.UserRepository
	jwtManager *auth.JWTManager
}

func NewUserService(userRepo repository.UserRepository, jwtManager *auth.JWTManager) *UserService {
	return &UserService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// ValidateEmail проверяет формат email (базовая проверка)
// Базовая проверка, возвращает ошибку если Email не соответствует требованиям
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}

	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(emailRegex, email)
	if err != nil {
		return fmt.Errorf("error checking email format: %v", err)
	}
	if !matched {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

func (s *UserService) Register(ctx context.Context, req *model.UserCreateRequest) (*model.TokenResponse, error) {
	// 1. Валидация входных данных
	if len(req.Username) < 3 {
		return nil, fmt.Errorf("имя пользователя должно быть не менее 3 символов")
	}
	if err := ValidateEmail(req.Email); err != nil {
		return nil, fmt.Errorf("неверный email: %w", err)
	}
	if len(req.Password) < 6 {
		return nil, fmt.Errorf("пароль должен быть не менее 6 символов")
	}

	// 2. Проверка уникальности email
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("ошибка при проверке email: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("email уже занят")
	}

	// 3. Проверка уникальности username
	exist, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("ошибка при проверке имени пользователя: %w", err)
	}
	if exist {
		return nil, fmt.Errorf("имя пользователя уже занято")
	}

	// 4. Хеширование пароля
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("ошибка при хешировании пароля: %w", err)
	}

	// 5. Создание модели пользователя
	user := &model.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 6. Сохранение пользователя
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("ошибка при сохранении пользователя: %w", err)
	}

	// 7. Генерация JWT токена
	//func NewJWTManager(secretKey string, ttlHours int) *JWTManager {
	//func (m *JWTManager) GenerateToken(userID int, email string, username string) (string, time.Time, error) {
	token, _, err := auth.NewJWTManager("", 24).GenerateToken(user.ID, user.Email, user.Username)
	if err != nil {
		return nil, fmt.Errorf("ошибка при генерации токена: %w", err)
	}

	// 8. Возврат результата
	// UserResponse - структура для ответа с данными пользователя (без пароля)
	//type UserResponse struct {
	//	ID        int       `json:"id"`
	//	Username  string    `json:"username"`
	//	Email     string    `json:"email"`
	//	CreatedAt time.Time `json:"created_at"`
	//}

	return &model.TokenResponse{
		Token: token,
		User:  model.UserResponse{ID: user.ID, Username: user.Username, Email: user.Email}}, nil
}

func (s *UserService) Login(ctx context.Context, req *model.UserLoginRequest) (*model.TokenResponse, error) {
	// TODO: Реализовать вход пользователя
	// Шаги:
	// 1. Валидация входных данных
	// 2. Найти пользователя по email через репозиторий
	// 3. Проверить пароль используя функцию из пакета auth
	// 4. Сгенерировать JWT токен при успешной аутентификации
	// 5. Вернуть TokenResponse
	// ВАЖНО: При ошибке не раскрывать, что именно неправильно (email или пароль)
	//return nil, fmt.Errorf("not implemented")
	// 1. Валидация входных данных
	if err := ValidateEmail(req.Email); err != nil {
		return nil, fmt.Errorf("неверный email: %w", err)
	}
	if len(req.Password) < 6 {
		return nil, fmt.Errorf("неверные учетные данные")
	}

	// 2. Найти пользователя по email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		// В случае ошибки базы или если пользователь не найден, скрываем детали
		return nil, fmt.Errorf("неверные учетные данные (2)")
	}

	if user == nil {
		// Пользователь не найден
		return nil, fmt.Errorf("неверные учетные данные (1)")
	}

	// 3. Проверить пароль
	//func CheckPassword(password, hash string) bool {
	if ok := auth.CheckPassword(req.Password, user.Password); !ok {
		// Ошибка проверки пароля — скрываем детали
		return nil, fmt.Errorf("неверные учетные данные")
	}

	// 4. Генерация JWT при успешной аутентификации
	token, _, err := auth.NewJWTManager("", 24).GenerateToken(user.ID, user.Email, user.Username)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании токена") // Можно скрыть, если считается, что это внутреняя ошибка
	}

	// 5. Возвращаем TokenResponse
	return &model.TokenResponse{
		Token: token,
		User:  model.UserResponse{ID: user.ID, Username: user.Username, Email: user.Email}}, nil
}

func (s *UserService) GetByID(ctx context.Context, id int) (*model.User, error) {
	// TODO: Получить пользователя по ID через репозиторий
	user, err := s.GetByID(ctx, id) // Получаем пользователя через репозиторий
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении пользователя: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("пользователь с ID %d не найден", id)
	}

	return user, nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := s.GetByEmail(ctx, email) // Получаем пользователя через репозиторий
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении пользователя по email: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("пользователь с email %s не найден", email)
	}

	return user, nil
}

// validateUserCreateRequest проверяет корректность данных для регистрации
func validateUserCreateRequest(req *model.UserCreateRequest) error {
	// TODO: Реализовать проверку всех полей
	if req == nil {
		return errors.New("запрос не может быть nil")
	}

	// Проверка имени
	if strings.TrimSpace(req.Username) == "" {
		return errors.New("имя не может быть пустым")
	}

	// Проверка электронной почты
	if strings.TrimSpace(req.Email) == "" {
		return errors.New("электронная почта не может быть пустой")
	}

	if err := ValidateEmail(req.Email); err != nil {
		return errors.New("недопустимый формат электронной почты")
	}

	//	_, err := mail.ParseAddress(req.Email)
	//	if err != nil {
	//		return errors.New("недопустимый формат электронной почты")
	//	}

	// Проверка пароля
	if len(req.Password) < 6 {
		return errors.New("пароль должен содержать не менее 6 символов")
	}

	// Можно добавить дополнительные проверки, например, на уникальность имени пользователя или электронной почты

	return nil
}

// validateUserLoginRequest проверяет корректность данных для входа
func validateUserLoginRequest(req *model.UserLoginRequest) error {

	// Проверка валидности email
	//if _, err := mail.ParseAddress(req.Email); err != nil {
	if err := ValidateEmail(req.Email); err != nil {
		return fmt.Errorf("некорректный email: %w", err)
	}

	// Проверка длины пароля
	if utf8.RuneCountInString(req.Password) < 6 {
		return fmt.Errorf("пароль должен быть не менее 6 символов")
	}

	return nil
}
