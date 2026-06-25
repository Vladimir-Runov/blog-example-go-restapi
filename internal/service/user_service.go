package service

import (
	"blog-api/internal/model"
	"blog-api/internal/repository"
	"blog-api/pkg/auth"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/mail"
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

func (s *UserService) Register(ctx context.Context, req *model.UserCreateRequest) (*model.TokenResponse, error) {
	// TODO: Реализовать регистрацию пользователя
	// Шаги:
	// 1. Валидация входных данных (username >= 3 символов, email валидный, пароль >= 6 символов)
	// 2. Проверить уникальность email через репозиторий
	// 3. Проверить уникальность username через репозиторий
	// 4. Захешировать пароль используя пакет auth
	// 5. Создать модель пользователя с хешированным паролем
	// 6. Сохранить пользователя через репозиторий
	// 7. Сгенерировать JWT токен для нового пользователя
	// 8. Вернуть TokenResponse с токеном и данными пользователя
	// return nil, fmt.Errorf("not implemented")
	// 1. Валидация входных данных
	if len(req.Username) < 3 {
		return nil, fmt.Errorf("имя пользователя должно быть не менее 3 символов")
	}
	if !IsValidEmail(req.Email) {
		return nil, fmt.Errorf("некорректный формат email")
	}
	if len(req.Password) < 6 {
		return nil, fmt.Errorf("пароль должен быть не менее 6 символов")
	}

	// 2. Проверка уникальности email
	existingUser, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("ошибка при проверке email: %w", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("email уже зарегистрирован")
	}

	// 3. Проверка уникальности username
	existingUser, err = s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("ошибка при проверке имени пользователя: %w", err)
	}
	if existingUser != nil {
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
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("ошибка при сохранении пользователя: %w", err)
	}
	// 7. Генерация JWT токена
	token, err := auth.GenerateJWT(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, fmt.Errorf("ошибка при генерации токена: %w", err)
	}

	// 8. Возврат результата
	return &model.TokenResponse{
		Token: token,
		User:  &model.UserResponse{ID: user.ID, Username: user.Username, Email: user.Email},
	}, nil
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
	if !IsValidEmail(req.Email) {
		// Не раскрывать точную причину, возвращаем одинаковую ошибку
		return nil, fmt.Errorf("неверные учетные данные")
	}
	if len(req.Password) < 6 {
		return nil, fmt.Errorf("неверные учетные данные")
	}

	// 2. Найти пользователя по email
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		// В случае ошибки базы или если пользователь не найден, скрываем детали
		return nil, fmt.Errorf("неверные учетные данные")
	}
	if user == nil {
		// Пользователь не найден
		return nil, fmt.Errorf("неверные учетные данные")
	}

	// 3. Проверить пароль
	if err := auth.CheckPassword(req.Password, user.Password); err != nil {
		// Ошибка проверки пароля — скрываем детали
		return nil, fmt.Errorf("неверные учетные данные")
	}

	// 4. Генерация JWT при успешной аутентификации
	token, err := auth.GenerateJWT(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании токена") // Можно скрыть, если считается, что это внутреняя ошибка
	}

	// 5. Возвращаем TokenResponse
	return &model.TokenResponse{
		Token: token,
		User: &model.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	}, nil
}

func (s *UserService) GetByID(ctx context.Context, id int) (*model.User, error) {
	// TODO: Получить пользователя по ID через репозиторий

	return nil, fmt.Errorf("not implemented")
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	// TODO: Получить пользователя по email через репозиторий
	//return nil, fmt.Errorf("not implemented")
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			// Пользователь не найден
			return nil, nil
		}
		// Другие ошибки базы данных
		return nil, fmt.Errorf("ошибка при получении пользователя по email: %w", err)
	}
	return user, nil
}

// validateUserCreateRequest проверяет корректность данных для регистрации
func validateUserCreateRequest(req *model.UserCreateRequest) error {
	// TODO: Реализовать проверку всех полей

	return nil
}

// validateUserLoginRequest проверяет корректность данных для входа
func validateUserLoginRequest(req *model.UserLoginRequest) error {
	// TODO: Реализовать проверку полей
	// return nil
	// Проверка наличия и длины username
	if utf8.RuneCountInString(req.Username) < 3 {
		return fmt.Errorf("имя пользователя должно быть не менее 3 символов")
	}

	// Проверка валидности email
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return fmt.Errorf("некорректный email: %w", err)
	}

	// Проверка длины пароля
	if utf8.RuneCountInString(req.Password) < 6 {
		return fmt.Errorf("пароль должен быть не менее 6 символов")
	}

	return nil
}
