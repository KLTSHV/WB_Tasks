# Calendar HTTP Server

HTTP-сервер для работы с календарём событий. Поддерживает CRUD для событий и выборки событий на день/неделю/месяц. Хранение данных в памяти, без внешней бд.

## Возможности

- `POST /create_event` — создать событие
- `POST /update_event` — обновить событие
- `POST /delete_event` — удалить событие
- `GET /events_for_day` — события на день
- `GET /events_for_week` — события на неделю (7 дней с указанной даты)
- `GET /events_for_month` — события на месяц (календарный месяц указанной даты)
- Middleware логирования: метод, URL, время выполнения

Формат:
- POST: `application/x-www-form-urlencoded`
- GET: параметры через query string

Ответы:
- Успех: `{"result":"..."}`
- Ошибка: `{"error":"..."}`

Коды:
- `200` — успех
- `400` — ошибка ввода (например, неверная дата)
- `503` — бизнес-ошибка (например, удалить несуществующее)
- `500` — прочие ошибки

## Требования

- Go 1.20+

## Запуск

### Через флаг
```bash
go run ./cmd -port=8080

Через переменную окружения

PORT=8080 go run ./cmd

По умолчанию используется порт 8080.

Примеры запросов (curl)

Создать событие

curl -i -X POST "http://localhost:8080/create_event" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "user_id=1&date=2026-01-17&event=Call%20girlfriend"

Обновить событие

curl -i -X POST "http://localhost:8080/update_event" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "id=1&user_id=1&date=2026-01-18&event=Buy%20fowers"

Удалить событие

curl -i -X POST "http://localhost:8080/delete_event" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "id=1&user_id=1"

Получить события на день

curl -i "http://localhost:8080/events_for_day?user_id=1&date=2026-01-18"

Получить события на неделю

curl -i "http://localhost:8080/events_for_week?user_id=1&date=2026-01-17"

Получить события на месяц

curl -i "http://localhost:8080/events_for_month?user_id=1&date=2026-01-17"

Тесты

Запуск всех тестов:

go test ./...

Проверка на гонки данных:

go test ./... -race
