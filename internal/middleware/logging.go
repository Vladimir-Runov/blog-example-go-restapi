package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/google/uuid" 
)

// LoggingMiddleware provides request logging, CORS, recovery and other utility middleware
type LoggingMiddleware struct {
	logger *log.Logger
}

// NewLoggingMiddleware creates a new logging middleware instance
func NewLoggingMiddleware(logger *log.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

// Logger logs all HTTP requests
func (m *LoggingMiddleware) Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Реализовать логирование запросов
		// Шаги:
		// 1. Засечь время начала запроса
		// 2. Создать wrapper для ResponseWriter чтобы захватить статус код
		// 3. Вызвать следующий handler с wrapped writer
		// 4. После выполнения залогировать: метод, путь, IP, статус, время выполнения

		// Временная реализация
		//next(w, r)
		// 1. Засечь время начала запроса
		start := time.Now()

		// 2. Создать wrapper для ResponseWriter чтобы захватить статус код
		wrappedWriter := &ResponseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

		// 3. Вызвать следующий handler с wrapped writer
		next(wrappedWriter, r)

		// 4. После выполнения залогировать: метод, путь, IP, статус, время выполнения
		duration := time.Since(start)
		fmt.Printf("Method: %s, Path: %s, IP: %s, Status: %d, Duration: %v\n",
			r.Method, r.URL.Path, r.RemoteAddr, wrappedWriter.statusCode, duration)
	}
}

// Recovery восстанавливается после паник
func (m *LoggingMiddleware) Recovery(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Реализовать восстановление после паник
		// Шаги:
		// 1. Использовать defer с recover() для перехвата паник
		// 2. При панике залогировать ошибку
		// 3. Опционально: добавить stack trace
		// 4. Вернуть клиенту 500 Internal Server Error
		// 5. Вызвать следующий handler

		defer func() {
			if err := recover(); err != nil {
				// 1. При панике залогировать ошибку
				m.logger.Printf("Recovered from panic: %v", err)

				// 2. Опционально: добавить stack trace
				// Можно использовать пакет runtime для получения stack trace
				// buf := make([]byte, 1024)
				// n := runtime.Stack(buf, true)
				// m.logger.Printf("Stack trace: %s", buf[:n])

				// 3. Вернуть клиенту 500 Internal Server Error
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		// Временная реализация Вызвать следующий handler
		next(w, r)
	}
}

// CORS добавляет CORS заголовки
func (m *LoggingMiddleware) CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Реализовать CORS заголовки
		// Шаги:
		// 1. Добавить необходимые CORS заголовки (Origin, Methods, Headers, Max-Age)
		// 2. Обработать preflight запросы (OPTIONS метод) - вернуть 204
		// 3. Для остальных методов вызвать следующий handler

		// 1. Добавить необходимые CORS заголовки
		//1. CORS Заголовки:
		//– Access-Control-Allow-Origin: Указывает, какие источники могут обращаться к ресурсу. В данном случае установлен на "*", что разрешает все источники.
		//– Access-Control-Allow-Methods: Перечисляет методы, разрешенные для использования с данным ресурсом.
		//– Access-Control-Allow-Headers: Указывает заголовки, которые могут быть использованы при выполнении запроса.
		//– Access-Control-Max-Age: Указывает время (в секундах), в течение которого результаты preflight запроса могут быть кэшированы.
		w.Header().Set("Access-Control-Allow-Origin", "*") // Разрешить все источники
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400") // Кэширование preflight запроса на 24 часа

		// 2. Обработать preflight запросы (OPTIONS метод)
		if r.Method == http.MethodOptions { отка preflight запросов:
		//   – Если метод запроса — OPTIONS, возвращаем статус 204 No Content и завершаем выполнение функции. Это позволяет клиенту знать, что pre-flight запрос был успешен.
			w.WriteHeader(http.StatusNoContent) // Возвращаем 204 No Content
			return
		}

		// // 3. Для остальных методов вызвать следующий handler /Временная реализация
		next(w, r)
	}
}

// RequestID добавляет уникальный ID к каждому запросу - middleware будет добавлять уникальный идентификатор к каждому запросу и логировать его, что может помочь в отслеживании и отладке запросов.
func (m *LoggingMiddleware) RequestID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Реализовать генерацию Request ID
		// Шаги:
		// 1. Сгенерировать уникальный ID (UUID или timestamp+random)
		// 2. Добавить ID в контекст запроса для использования в логах
		// 3. Добавить ID в заголовок ответа X-Request-ID
		// 4. Залогировать запрос с Request ID
		// 5. Вызвать следующий handler
 		
		// 1. Сгенерировать уникальный ID
        requestID := uuid.New().String() // Генерация UUID github.com/google/uuid 

        // 2. Добавить ID в контекст запроса
        ctx := context.WithValue(r.Context(), "requestID", requestID)
        r = r.WithContext(ctx)

        // 3. Добавить ID в заголовок ответа X-Request-ID , чтобы клиент мог видеть уникальный идентификатор запроса.
        w.Header().Set("X-Request-ID", requestID)

        // 4. Залогировать запрос с Request ID
        log.Printf("RequestID: %s, Method: %s, URL: %s", requestID, r.Method, r.URL)

		// 5. Вызвать следующий handler 		
		next(w, r)
	}
}

// RateLimiter ограничивает количество запросов от одного клиента
func (m *LoggingMiddleware) RateLimiter(maxRequests int, window time.Duration) 
	func(http.HandlerFunc) http.HandlerFunc {
	// TODO: Реализовать rate limiting (продвинутое задание)
	// Шаги:
	// 1. Создать хранилище для отслеживания запросов по IP адресам
	 m.requestData = make(map[string]*requestInfo)
	

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			// 2. Использовать mutex для безопасного доступа к хранилищу
            m.mu.Lock()
            defer m.mu.Unlock()

			ip := r.RemoteAddr // Получаем IP клиента
			// 3. Для каждого запроса:
			//    - Получить IP клиента
			//    - Проверить количество запросов в текущем окне времени
			//    - Если превышен лимит - вернуть 429 Too Many Requests
			//    - Иначе увеличить счетчик и пропустить запрос

            // Проверяем или создаем запись для IP
            info, exists := m.requestData[ip]
            if !exists {
                info = &requestInfo{count: 0, firstSeen: time.Now()}
                m.requestData[ip] = info
            }

            // Проверяем, не истекло ли окно времени
            if time.Since(info.firstSeen) > window {
                info.count = 0 // Сбрасываем счетчик
                info.firstSeen = time.Now() // Обновляем время первого запроса
            }

            // Проверяем количество запросов
            if info.count >= maxRequests {
                http.Error(w, "429 Too Many Requests", http.StatusTooManyRequests)
                return
            }

            // Увеличиваем счетчик и пропускаем запрос
            info.count++

			next(w, r)
		}
	}
}

// ContentTypeJSON устанавливает Content-Type: application/json для всех ответов
func (m *LoggingMiddleware) ContentTypeJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Установить Content-Type: application/json для всех ответов
		w.Header().Set("Content-Type", "application/json")

        // Вызываем следующий обработчик
        next(w, r)		
	}
}

// getClientIP извлекает IP адрес клиента
// позволяет корректно извлекать IP-адрес клиента в различных сценариях, включая использование прокси-серверов.
func getClientIP(r *http.Request) string {
	// TODO: Извлечь реальный IP адрес клиента
	// Проверить заголовки: X-Forwarded-For, X-Real-IP, затем RemoteAddr
	// Учесть что X-Forwarded-For может содержать несколько IP
	//return r.RemoteAddr
	// Проверяем заголовок X-Forwarded-For
	//– Сначала проверяется заголовок X-Forwarded-For. Если он не пустой, то предполагается, что он может содержать несколько IP-адресов, разделенных запятыми. Мы берем первый IP и убираем лишние пробелы.
	//– Затем проверяется заголовок X-Real-IP. Если он установлен, то возвращается его значение.
	//
	//	2. Использование RemoteAddr:
	//	– Если оба заголовка отсутствуют, используется RemoteAddr, который содержит IP-адрес клиента. Функция net.SplitHostPort разбивает строку на IP и порт. Если возникает ошибка при разбиении, возвращается оригинальное значение RemoteAddr.
    xForwardedFor := r.Header.Get("X-Forwarded-For")
    if xForwardedFor != "" {
        // X-Forwarded-For может содержать несколько IP, берем первый
        ips := strings.Split(xForwardedFor, ",")
        return strings.TrimSpace(ips[0])
    }

    // Проверяем заголовок X-Real-IP
    xRealIP := r.Header.Get("X-Real-IP")
    if xRealIP != "" {
        return xRealIP
    }

    // используем RemoteAddr
    ip, _, err := net.SplitHostPort(r.RemoteAddr)
    if err != nil {
        return r.RemoteAddr // Возвращаем RemoteAddr в случае ошибки
    }
    return ip	
}

// responseWriter обертка для захвата статус кода
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

// WriteHeader сохраняет статус код
func (rw *responseWriter) WriteHeader(code int) {
	if !rw.written {
		rw.statusCode = code
		rw.ResponseWriter.WriteHeader(code)
		rw.written = true
	}
}

// Write вызывает WriteHeader если еще не был вызван
func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

// newResponseWriter создает новую обертку
func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		written:        false,
	}
}
