# Blog API - Шаблон проектной работы

## Описание проекта (blog-api)

Вам необходимо реализовать REST API для блога с функциональностью:

* Аутентификация пользователей (JWT)
* CRUD операции для постов
* Комментарии к постам
* Авторизация (только автор может редактировать/удалять свои посты и комментарии)

## Структура проекта

```
blog-api/
├── cmd/api/              # Точка входа приложения
│   └── main.go
├── internal/             # Внутренние пакеты приложения
│   ├── model/           # Модели данных
│   ├── handler/         # HTTP хендлеры
│   ├── service/         # Бизнес-логика
│   ├── repository/      # Работа с БД
│   └── middleware/      # HTTP middleware
├── pkg/                 # Переиспользуемые пакеты
│   ├── auth/           # JWT и пароли
│   └── database/       # Подключение к БД
├── migrations/         # SQL миграции
├── docker-compose.yml  # PostgreSQL и Adminer
├── .env.example        # Пример конфигурации
├── go.mod
└── README.md
```

## Начало работы

### 1\. Подготовка окружения

```bash
# Клонировать шаблон
cp -r template my-blog-api
cd my-blog-api

# Установить зависимости
go mod download
go mod tidy

# Создать файл конфигурации
copy .env.example .env

# Запустить PostgreSQL
docker-compose up -d

# Проверить что БД работает (опционально)
docker-compose logs postgres
```

### 2\. Порядок реализации

#### Этап 1: Базовая инфраструктура

1. **pkg/database/postgres.go**

   * Реализовать подключение к БД
   * Реализовать функцию миграций
2. **pkg/auth/password.go**

   * Реализовать хеширование паролей (bcrypt)
   * Реализовать проверку пароля
3. **pkg/auth/jwt.go**

   * Реализовать генерацию JWT токенов
   * Реализовать валидацию токенов

#### Этап 2: Репозитории

1. **internal/repository/user\_repo.go**

   * Завершить реализацию всех методов
   * SQL запросы уже подготовлены
2. **internal/repository/post\_repo.go**

   * Завершить реализацию CRUD операций
   * Добавить методы пагинации
3. **internal/repository/comment\_repo.go**

   * Завершить реализацию работы с комментариями

#### Этап 3: Бизнес-логика

1. **internal/service/user\_service.go**

   * Регистрация с валидацией
   * Вход с проверкой пароля
   * Генерация JWT токена
2. **internal/service/post\_service.go**

   * CRUD операции с проверкой прав
   * Пагинация и фильтрация
3. **internal/service/comment\_service.go**

   * Создание комментариев
   * Проверка прав на редактирование

#### Этап 4: HTTP слой

1. **internal/handler/auth\_handler.go**

   * Обработка регистрации и входа
   * Возврат JWT токена
2. **internal/handler/post\_handler.go**

   * REST эндпоинты для постов
   * Обработка ошибок
3. **internal/handler/comment\_handler.go**

   * Эндпоинты для комментариев

#### Этап 5: Middleware

1. **internal/middleware/auth.go**

   * JWT проверка
   * Добавление user\_id в контекст
2. **internal/middleware/logging.go**

   * Логирование запросов
   * Recovery от паник
   * CORS заголовки

#### Этап 6: Сборка приложения

1. **cmd/api/main.go**

   * Инициализация всех компонентов
   * Настройка маршрутов
   * Запуск сервера

## API Эндпоинты

### Публичные (без аутентификации)

* `POST /api/register` - регистрация пользователя
* `POST /api/login` - вход пользователя
* `GET /api/posts` - список постов
* `GET /api/posts/{id}` - получить пост
* `GET /api/posts/{id}/comments` - комментарии к посту

### Защищенные (требуют JWT токен)

* `POST /api/posts` - создать пост
* `PUT /api/posts/{id}` - обновить пост (только автор)
* `DELETE /api/posts/{id}` - удалить пост (только автор)
* `POST /api/posts/{id}/comments` - создать комментарий к посту

## Требования к реализации

### Обязательные требования

* ✅ Все основные эндпоинты работают
* ✅ JWT аутентификация реализована
* ✅ Проверка прав доступа работает
* ✅ Валидация входных данных
* ✅ Обработка ошибок
* ✅ Пагинация для списков

### Дополнительные требования (для высокой оценки)

* 📊 Кеширование часто запрашиваемых данных
* 🔍 Поиск и фильтрация постов
* 📝 Подробное логирование
* ⚡ Оптимизированные SQL запросы
* 🧪 Юнит-тесты для критической логики
* 📚 API документация (Swagger/OpenAPI)

## Полезные команды

```bash
# Запуск приложения
go run cmd/api/main.go

# Запуск с hot-reload (установить air)
air

# Тестирование
go test ./...

# Проверка на ошибки
go vet ./...
golangci-lint run

# Форматирование кода
go fmt ./...

# База данных
docker-compose up -d    # Запустить
docker-compose down      # Остановить
docker-compose logs -f   # Логи

# Примеры запросов
# Регистрация
curl -X POST http://localhost:8080/api/register \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
  -H "Content-Type: application/json" \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'

# Вход
curl -X POST http://localhost:8080/api/login \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
  -H "Content-Type: application/json" \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
  -d '{"email":"test@example.com","password":"password123"}'

# Создание поста (с токеном)
curl -X POST http://localhost:8080/api/posts \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
  -H "Authorization: Bearer YOUR\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\_TOKEN" \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
  -H "Content-Type: application/json" \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
  -d '{"title":"My Post","content":"Post content"}'
```

## Где искать подсказки

1. **TODO комментарии** - в каждом файле есть указания что нужно реализовать
2. **Интерфейсы** - в `repository/interfaces.go` описаны все методы
3. **Модели** - в `model/` определены все структуры данных
4. **SQL запросы** - базовые запросы уже есть в репозиториях
5. **Примеры из solution** - можете подсмотреть в готовое решение при затруднениях

## Частые ошибки

1. **Не забудьте обработку ошибок** - всегда проверяйте err != nil
2. **SQL injection** - используйте placeholder'ы ($1, $2) в SQL запросах
3. **Контекст** - передавайте context во все методы работы с БД
4. **Закрытие ресурсов** - используйте defer для rows.Close()
5. **Права доступа** - проверяйте что пользователь может редактировать только свои данные

## Критерии оценки

### Минимум для зачета (60%)

* Работают эндпоинты регистрации и входа
* Можно создавать и получать посты
* JWT токены генерируются и проверяются

### Хорошо (80%)

* Все CRUD операции работают
* Реализована проверка прав доступа
* Корректная обработка ошибок
* Пагинация работает

### Отлично (100%)

* Код хорошо структурирован
* Добавлены дополнительные функции
* Есть тесты
* Оптимизированы запросы к БД
* Документирован API

## Полезные ссылки

* [Go database/sql tutorial](http://go-database-sql.org/)
* [JWT in Go](https://github.com/golang-jwt/jwt)
* [Chi router](https://github.com/go-chi/chi)
* [bcrypt in Go](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
* [PostgreSQL documentation](https://www.postgresql.org/docs/)

## Вопросы?

Если возникли вопросы:

1. Проверьте TODO комментарии в коде
2. Посмотрите примеры в готовом решении
3. Обратитесь к преподавателю

Удачи в реализации! 🚀



docker-compose up -d
time="2026-06-30T23:34:24+03:00" level=warning msg="C:\\\\\\\\\\\\\\\\Users\\\\\\\\\\\\\\\\Admin\\\\\\\\\\\\\\\\Documents\\\\\\\\\\\\\\\\go\\\\\\\\\\\\\\\\git\\\\\\\_netology\\\\\\\\\\\\\\\\blog-example-go-restapi\\\\\\\\\\\\\\\\docker-compose.yml: the attribute `version` is obsolete, it will be ignored, please remove it to avoid potential confusion"
\\\\\\\[+] up 2/2
✔ Container blog\\\\\\\_postgres Running                                                                                  0.0s
✔ Container blog\\\\\\\_adminer  Running                                                                                  0.0s



docker-compose ps
time="2026-06-30T23:35:01+03:00" level=warning msg="C:\\\\\\\\\\\\\\\\Users\\\\\\\\\\\\\\\\Admin\\\\\\\\\\\\\\\\Documents\\\\\\\\\\\\\\\\go\\\\\\\\\\\\\\\\git\\\\\\\_netology\\\\\\\\\\\\\\\\blog-example-go-restapi\\\\\\\\\\\\\\\\docker-compose.yml: the attribute `version` is obsolete, it will be ignored, please remove it to avoid potential confusion"
NAME            IMAGE                COMMAND                  SERVICE    CREATED         STATUS                   PORTS
blog\\\\\\\_adminer    adminer:latest       "entrypoint.sh docke…"   adminer    9 minutes ago   Up 9 minutes             0.0.0.0:8081->8080/tcp, \\\\\\\[::]:8081->8080/tcp
blog\\\\\\\_postgres   postgres:15-alpine   "docker-entrypoint.s…"   postgres   9 minutes ago   Up 9 minutes (healthy)   0.0.0.0:5432->5432/tcp, \\\\\\\[::]:5432->5432/tcp





curl -X POST http://localhost:8080/api/register -H "Content-Type: application/json" -d "{\\"username\\": \\"tester\_001\\", \\"password\\": \\"tester\_pass123\\", \\"email\\": \\"tester\_001@example.com\\"}"

{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJlbWFpbCI6InRlc3Rlcl8wMDFAZXhhbXBsZS5jb20iLCJ1c2VybmFtZSI6InRlc3Rlcl8wMDEiLCJleHAiOjE3ODI5NDE5OTAsImlhdCI6MTc4Mjg1NTU5MH0.0aAzqcbeorTSBsFW41dqTdHWVys3XPVdnbx3rz\_DAig"}





curl -X POST http://localhost:8080/api/login -H "Content-Type: application/json" -d "{\\"email\\": \\"tester\_001@example.com\\", \\"password\\": \\"tester\_pass123\\"}"

{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJlbWFpbCI6InRlc3Rlcl8wMDFAZXhhbXBsZS5jb20iLCJ1c2VybmFtZSI6InRlc3Rlcl8wMDEiLCJleHAiOjE3ODI5NDM0OTYsImlhdCI6MTc4Mjg1NzA5Nn0.jdB3xKSRTAT8dy3MMrb3Zs8-2v42u78-KO37mLk-2k8"}







C:\\Users\\Admin\\Documents\\go\\git\_netology\\blog-example-go-restapi>docker exec -it blog\_postgres bash

98c4dfdef89a:/#  psql -U bloguser -d blogdb\_test

psql: error: connection to server on socket "/var/run/postgresql/.s.PGSQL.5432" failed: FATAL:  database "blogdb\_test" does not exist

98c4dfdef89a:/# CREATE DATABASE blogdb\_test;

bash: CREATE: command not found

98c4dfdef89a:/# psql -U bloguser -d postgres

psql (15.18)

Type "help" for help.



postgres=# CREATE DATABASE blogdb\_test;

CREATE DATABASE

postgres=#    \\l

&#x20;                                                List of databases

&#x20;   Name     |  Owner   | Encoding |  Collate   |   Ctype    | ICU Locale | Locale Provider |   Access privileges

\-------------+----------+----------+------------+------------+------------+-----------------+-----------------------

&#x20;blogdb      | bloguser | UTF8     | en\_US.utf8 | en\_US.utf8 |            | libc            |

&#x20;blogdb\_test | bloguser | UTF8     | en\_US.utf8 | en\_US.utf8 |            | libc            |





go.exe test -test.fullpath=true -timeout 30s -run ^TestAuthHandler\_Register$ blog-example-go-restapi/internal/handler

ok  	blog-example-go-restapi/internal/handler	1.173s



