package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

// Claims представляет данные, хранимые в JWT токене
type Claims struct {
	UserID               int    `json:"user_id"`
	Email                string `json:"email"`
	Username             string `json:"username"`
	jwt.RegisteredClaims        // TODO: Добавить стандартные JWT claims
	// Подсказка: используйте jwt.RegisteredClaims или jwt.StandardClaims
}

// JWTManager управляет созданием и валидацией JWT токенов
type JWTManager struct {
	secretKey []byte
	ttl       time.Duration
}

// NewJWTManager создает новый экземпляр JWT менеджера
func NewJWTManager(secretKey string, ttlHours int) *JWTManager {
	// TODO: Инициализировать JWTManager
	// - Преобразовать secretKey в []byte
	// - Преобразовать ttlHours в time.Duration
	return &JWTManager{
		secretKey: []byte(secretKey),
		ttl:       time.Duration(ttlHours) * time.Hour,
	}
}

// GenerateToken создает новый JWT токен для пользователя
func (m *JWTManager) GenerateToken(userID int, email, username string) (string, time.Time, error) {
	// TODO: Реализовать генерацию JWT токена
	// Шаги:
	// 1. Создать Claims с данными пользователя
	// 2. Установить время истечения токена (текущее время + ttl)
	// 3. Создать токен используя алгоритм подписи (например, HS256)
	// 4. Подписать токен секретным ключом
	// 5. Вернуть подписанную строку токена и время истечения
	//
	// Подсказка: используйте библиотеку github.com/golang-jwt/jwt/v5
	//return "", time.Time{}, errors.New("not implemented")

	expirationTime := time.Now().Add(m.ttl) // Время истечения токена
	claims := Claims{
		UserID:   userID,
		Email:    email,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime), // Установить время истечения
			IssuedAt:  jwt.NewNumericDate(time.Now()),     // Установить время выпуска
		},
	}

	// Шаг 3: Создать токен используя алгоритм подписи (например, HS256)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Шаг 4: Подписать токен секретным ключом
	signedToken, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", time.Time{}, err // Вернуть ошибку в случае сбоя
	}

	// Шаг 5: Вернуть подписанную строку токена и время истечения
	return signedToken, expirationTime, nil
}

// ValidateToken проверяет и парсит JWT токен
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	// TODO: Реализовать валидацию и парсинг JWT токена
	// Шаги:
	// 1. Распарсить токен с проверкой подписи
	// 2. Извлечь claims из токена
	// 3. Проверить время истечения токена
	// 4. Вернуть claims если токен валидный
	//
	// Обработать ошибки:
	// - Невалидная подпись -> ErrInvalidToken
	// - Истекший токен -> ErrExpiredToken
	// - Другие ошибки -> ErrInvalidToken

	//return nil, errors.New("not implemented")
	// Шаг 1: Распарсить токен с проверкой подписи
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверка метода подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return m.secretKey, nil // Возвращаем секретный ключ для проверки подписи
	})

	if err != nil {
		// Обработка ошибок парсинга токена
		if err == jwt.ErrSignatureInvalid {
			return nil, ErrInvalidToken // Невалидная подпись
		}
		return nil, ErrInvalidToken // Другие ошибки
	}

	// Шаг 2: Извлечь claims из токена
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken // Если токен не валиден или не удалось извлечь claims
	}

	// Шаг 3: Проверить время истечения токена
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, ErrExpiredToken // Истекший токен
	}

	// Шаг 4: Вернуть claims если токен валидный
	return claims, nil
}

// RefreshToken обновляет существующий токен (опциональное задание)
func (m *JWTManager) RefreshToken(tokenString string) (string, time.Time, error) {
	// TODO: Реализовать обновление токена (продвинутое задание)
	// Шаги:
	// 1. Валидировать существующий токен
	// 2. Извлечь данные пользователя из старого токена
	// 3. Сгенерировать новый токен с теми же данными
	// 4. Вернуть новый токен

	//return "", time.Time{}, errors.New("not implemented")
	// Шаг 1: Валидировать существующий токен
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return "", time.Time{}, err // Вернуть ошибку, если токен не валиден
	}

	// Шаг 2: Извлечь данные пользователя из старого токена
	userID := claims.UserID // Предполагается, что в claims есть поле UserID
	// Здесь можно добавить дополнительные данные, если они есть в claims

	// Шаг 3: Сгенерировать новый токен с теми же данными
	newClaims := &Claims{
		UserID: userID,
		// Здесь можно добавить дополнительные поля для claims
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.tokenLifetime)), // Установка нового времени истечения
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	tokenString, err = newToken.SignedString(m.secretKey)
	if err != nil {
		return "", time.Time{}, err // Вернуть ошибку, если не удалось подписать новый токен
	}

	// Шаг 4: Вернуть новый токен и время истечения
	return tokenString, newClaims.ExpiresAt.Time, nil

}

// GetUserIDFromToken быстро извлекает ID пользователя из токена без полной валидации
func (m *JWTManager) GetUserIDFromToken(tokenString string) (int, error) {
	// TODO: Извлечь UserID из токена (опциональное задание)
	// Может быть полезно для быстрой проверки

	//return 0, errors.New("not implemented")
	// Шаг 1: Разобрать токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверка метода подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secretKey, nil
	})

	// Шаг 2: Проверить наличие ошибок при разборе токена
	if err != nil {
		return 0, err // Вернуть ошибку, если токен не может быть разобран
	}

	// Шаг 3: Извлечь UserID из токена
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["user_id"].(float64); ok { // Предполагается, что UserID хранится как float64
			return int(userID), nil // Вернуть UserID как int
		}
	}

	return 0, errors.New("user ID not found in token")
}
