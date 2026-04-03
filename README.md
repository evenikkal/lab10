# Лабораторная работа №10 — FastAPI (Python) vs Gin (Go)

**Студент:** Никишина Евгения Александровна  
**Группа:** 221131  
**Тема:** Веб-разработка: FastAPI (Python) vs Gin (Go) + gRPC + JWT

---

## Структура репозитория

```
.
├── task_m1_gin_api/          # М1: REST API на Gin (книги)
│   ├── cmd/server/           #     Точка входа
│   └── internal/             #     handler, model, repository
├── task_m3_validation/       # М3: Валидация входных данных в Go
├── task_m5_json_exchange/    # М5: Передача сложных JSON-структур
│   ├── go_service/           #     Go-сервис (Orders API)
│   └── python_client/        #     Python-клиент
├── task_v1_grpc/             # В1: gRPC-сервер (Go) + клиент (Python)
│   ├── proto/                #     Исходный .proto файл
│   ├── go_server/            #     Go gRPC-сервер + сгенерированные pb-файлы
│   └── python_client/        #     Python gRPC-клиент + pb2-файлы
├── task_v3_jwt/              # В3: JWT-аутентификация в Go + верификация из Python
│   ├── go_service/           #     Go-сервис с /login и /protected
│   └── python_client/        #     Python-клиент с PyJWT верификацией
├── PROMPT_LOG.md             # Лог промптов
└── README.md
```

---

## Требования

| Инструмент | Версия |
|------------|--------|
| Go         | ≥ 1.21 |
| Python     | ≥ 3.10 |
| protoc     | любая (только для регенерации proto) |

---

## М1 — REST API на Gin

Управление коллекцией книг с in-memory хранилищем и потокобезопасным доступом (`sync.RWMutex`).

**Эндпоинты:** `GET /health`, `GET /books`, `GET /books/:id`, `POST /books`

```bash
cd task_m1_gin_api
go mod tidy
go test ./... -v
go run cmd/server/main.go    # сервер на :8080
```

Примеры запросов:

```bash
curl http://localhost:8080/health
curl http://localhost:8080/books
curl http://localhost:8080/books/1
curl -X POST http://localhost:8080/books \
     -H 'Content-Type: application/json' \
     -d '{"title":"1984","author":"George Orwell","year":1949}'
```

---

## М3 — Валидация входных данных

Два эндпоинта с binding-валидацией через теги Gin.

- `POST /users` — name (required, min=2, max=50), email (required, валидный), age (required, 18–100)
- `POST /products` — title (required, min=3), price (required, >0), quantity (required, ≥0), category (oneof: electronics, food, clothing, other)

```bash
cd task_m3_validation
go mod tidy
go test ./... -v
go run main.go               # сервер на :8081
```

Примеры запросов:

```bash
# Валидный пользователь
curl -X POST http://localhost:8081/users \
     -H 'Content-Type: application/json' \
     -d '{"name":"Alice","email":"alice@example.com","age":25}'

# Невалидный (age < 18)
curl -X POST http://localhost:8081/users \
     -H 'Content-Type: application/json' \
     -d '{"name":"Alice","email":"alice@example.com","age":15}'
```

---

## М5 — Передача сложных JSON-структур

Go-сервис принимает `Order` с вложенными `OrderItem[]` и `Address`, автоматически вычисляет сумму.

**Эндпоинты:** `POST /orders`, `GET /orders/:id`, `GET /orders`

```bash
# Go-сервис
cd task_m5_json_exchange/go_service
go mod tidy
go test ./... -v
go run main.go               # сервер на :8082

# Python-клиент (в другом терминале)
cd task_m5_json_exchange/python_client
pip install -r requirements.txt
python client.py             # нужен запущенный Go-сервис
python -m pytest test_client.py -v   # тесты без сервера
```

---

## В1 — gRPC-сервер (Go) + клиент (Python)

Go-сервер реализует Greeter-сервис с методом `SayHello`. Python-клиент вызывает его через gRPC.

```bash
# Go-сервер
cd task_v1_grpc/go_server
go mod tidy
go test ./... -v             # тесты через bufconn (без реального порта)
go run main.go               # gRPC сервер на :50051

# Python-клиент (в другом терминале)
cd task_v1_grpc/python_client
pip install -r requirements.txt
python client.py             # нужен запущенный Go-сервер
python -m pytest test_client.py -v   # тесты без сервера
```

### Регенерация proto-файлов (опционально)

```bash
# Go
cd task_v1_grpc/go_server
protoc --proto_path=../proto \
       --go_out=pb --go_opt=paths=source_relative \
       --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
       ../proto/hello.proto

# Python
cd task_v1_grpc/python_client
python -m grpc_tools.protoc \
       --proto_path=../proto \
       --python_out=. --grpc_python_out=. \
       ../proto/hello.proto
```

---

## В3 — JWT-аутентификация

Go-сервис выдаёт HS256 JWT на `/login`. Защищённые эндпоинты `/protected` и `/profile` требуют Bearer-токен. Python-клиент верифицирует токен независимо через PyJWT с общим секретом.

Тестовые пользователи: `alice / password123`, `bob / qwerty`.

```bash
# Go-сервис
cd task_v3_jwt/go_service
go mod tidy
go test ./... -v
go run main.go               # сервер на :8083

# Python-клиент (в другом терминале)
cd task_v3_jwt/python_client
pip install -r requirements.txt
python client.py             # нужен запущенный Go-сервис
python -m pytest test_client.py -v   # тесты без сервера
```

Примеры запросов:

```bash
# Получить токен
curl -X POST http://localhost:8083/login \
     -H 'Content-Type: application/json' \
     -d '{"username":"alice","password":"password123"}'

# Вызвать защищённый эндпоинт
TOKEN="<token из ответа>"
curl http://localhost:8083/protected -H "Authorization: Bearer $TOKEN"
curl http://localhost:8083/profile   -H "Authorization: Bearer $TOKEN"
```

---

## Запуск всех тестов

```bash
# Go
for dir in task_m1_gin_api task_m3_validation \
           task_m5_json_exchange/go_service \
           task_v1_grpc/go_server task_v3_jwt/go_service; do
    echo "=== $dir ==="
    (cd $dir && go mod tidy && go test ./... -v)
done

# Python
for dir in task_m5_json_exchange/python_client \
           task_v1_grpc/python_client \
           task_v3_jwt/python_client; do
    echo "=== $dir ==="
    (cd $dir && python -m pytest test_client.py -v)
done
```

---