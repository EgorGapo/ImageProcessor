# Image Processing API

REST API сервис для обработки изображений с асинхронной обработкой через RabbitMQ, аутентификацией и кешированием в Redis.

##  Требования

- Go 1.24+
- PostgreSQL 12+
- Redis 6+
- RabbitMQ 3.8+
- Docker & Docker Compose

## API Endpoints

### Аутентификация

#### POST `/auth/register` - Регистрация
Создает новую учетную запись пользователя. Не требует аутентификации.

**Запрос:**
```json
{
  "username": "user@example.com",
  "password": "securepassword123"
}
```

**Ответы:**
- `200 OK`: `{"value": "Registration successful"}`
- `400 Bad Request`: Ошибка валидации данных
- `500 Internal Server Error`: Ошибка на сервере

---

#### POST `/auth/login` - Вход в систему
Аутентифицирует пользователя и возвращает токен сессии. Не требует аутентификации.

**Запрос:**
```json
{
  "username": "user@example.com",
  "password": "securepassword123"
}
```

**Ответ (200 OK):**
```json
{
  "value": "log in was succesful",
  "token": "uuid-session-token-here"
}
```

**Ошибки:**
- `400 Bad Request`: Некорректные учетные данные или ошибка парсинга
- `500 Internal Server Error`: Ошибка сервера

**Middleware:** Нет

---

###  Управление Задачами (требуют аутентификации)

#### POST `/task` - Создать задачу обработки
Создает новую задачу обработки изображения с указанными фильтрами. **Требует аутентификации**.

**Middleware:** `AuthMiddleware` - Проверяет наличие валидного токена в заголовке `Authorization: Bearer {token}`

**Запрос:**
```json
{
  "image": "data:image/jpeg;base64,/9j/4AAQSkZJRg...",
  "filter": {
    "name": "blur",
    "parameters": {
      "radius": 10
    }
  }
}
```

**Поддерживаемые фильтры:**
- `blur` - Размытие изображения
  - Параметры: `{"radius": число}`
- `grayscale` - Преобразование в оттенки серого
  - Параметры: `{}`
- `sepia` - Эффект сепия
  - Параметры: `{"intensity": 0-100}`
- `sharpen` - Повышение резкости
  - Параметры: `{"amount": число}`

**Ответ (200 OK):**
```json
{
  "value": "task-uuid-1234-5678"
}
```

**Ошибки:**
- `400 Bad Request`: Ошибка парсинга, некорректные фильтры
- `401 Unauthorized`: Отсутствует или невалидный токен
- `500 Internal Server Error`: Ошибка при создании задачи в БД или отправке в очередь

**Процесс:**
1. Проверяется токен аутентификации
2. Парсятся параметры фильтра из тела запроса
3. Создается новая задача в PostgreSQL
4. Задача отправляется в RabbitMQ очередь для асинхронной обработки
5. Возвращается ID задачи

---

#### GET `/task/status/{taskID}` - Получить статус задачи
Получает текущий статус задачи обработки. **Требует аутентификации**.

**Middleware:** `AuthMiddleware` - Проверяет наличие валидного токена

**Параметры пути:**
- `taskID` (string, обязательный) - Идентификатор задачи

**Заголовки:**
- `Authorization: Bearer {token}` (обязательный)

**Ответ (200 OK):**
```json
{
  "value": "processing"
}
```

**Возможные статусы:**
- `pending` - Задача в очереди, ожидает обработки
- `processing` - Задача в процессе обработки
- `completed` - Обработка завершена успешно
- `failed` - Ошибка при обработке

**Ошибки:**
- `400 Bad Request`: Некорректный ID задачи, ошибка парсинга
- `401 Unauthorized`: Отсутствует или невалидный токен
- `500 Internal Server Error`: Ошибка при получении статуса из БД

---

#### GET `/task/result/{taskID}` - Получить результат обработки
Получает обработанное изображение в формате PNG. **Требует аутентификации**.

**Middleware:** `AuthMiddleware` - Проверяет наличие валидного токена

**Параметры пути:**
- `taskID` (string, обязательный) - Идентификатор задачи

**Заголовки:**
- `Authorization: Bearer {token}` (обязательный)

**Ответ (200 OK):**
```
Бинарные данные изображения PNG
Content-Type: image/png
```

**Ошибки:**
- `202 Accepted`: Результат еще не готов, задача все еще обрабатывается
  - Сообщение: `"No result yet, please check the status"`
- `400 Bad Request`: Некорректный ID задачи
- `401 Unauthorized`: Отсутствует или невалидный токен
- `500 Internal Server Error`: Ошибка кодирования изображения или получения результата

**Примечание:** Результат доступен только после того, как статус задачи изменится на `completed`.

---

###  Внутренние операции

#### POST `/commit` - Завершить обработку задачи (Внутреннее)
Сохраняет результаты обработки изображения. Используется потребителем RabbitMQ для сохранения результатов.

**Запрос:**
```json
{
  "id": "task-uuid",
  "image_base": "data:image/jpeg;base64,...",
  "filter_name": "blur",
  "filter_parameters": "{\"radius\": 10}",
  "status": "completed",
  "result": "data:image/png;base64,iVBORw0KGgoAAAANSUhE..."
}
```

**Ответ (200 OK):**
```json
{
  "value": "task was done, you can check"
}
```

**Ошибки:**
- `400 Bad Request`: Ошибка парсинга тела запроса
- `500 Internal Server Error`: Ошибка при сохранении задачи в БД

**Процесс:**
1. Парсятся данные задачи из тела запроса
2. Данные сохраняются в PostgreSQL
3. Результат становится доступным через `/task/result/{taskID}`

---

## Полный поток обработки

```
1. Пользователь регистрируется (POST /auth/register)
   ↓
2. Пользователь входит (POST /auth/login) → получает токен
   ↓
3. Пользователь создает задачу (POST /task) с фильтром
   ↓
4. API сохраняет задачу в PostgreSQL
   ↓
5. API отправляет задачу в RabbitMQ очередь
   ↓
6. Потребитель (ImageProcessor) получает задачу из очереди
   ↓
7. Потребитель применяет фильтр к изображению
   ↓
8. Потребитель вызывает POST /commit с результатом
   ↓
9. API сохраняет результат в PostgreSQL
   ↓
10. Пользователь проверяет статус (GET /task/status/{taskID})
    ↓
11. После завершения получает результат (GET /task/result/{taskID})
```

---

##  База данных

```sql
-- Пользователи
CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    login VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL
);

-- Задачи обработки
CREATE TABLE tasks (
    id VARCHAR(255) PRIMARY KEY,
    image_base TEXT NOT NULL,
    filter_name VARCHAR(255),
    filter_parameters TEXT,
    status VARCHAR(50),
    result TEXT
);

-- Сессии
CREATE TABLE sessions (
    user_id VARCHAR(255) NOT NULL,
    session_id VARCHAR(255) PRIMARY KEY,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

##  Docker сервисы

- **PostgreSQL**: порт 5432 (для хранения пользователей, задач, сессий)
- **Redis**: порт 6379 (для кеширования сессий)
- **RabbitMQ**: порт 5672 (очередь задач), Admin UI: http://localhost:15672

##  Тестирование

```bash
go test ./...
go test -cover ./...
```
