# Project X: ДОВАЙБКОДИЛИСЬ

## Ветвление

Каждая новая фича или багфикс разрабатывается в отдельной ветке от `master`.

Именование веток:

```
dev-<task-id>
```

Где `<task-id>` — номер задачи из Github Boards.

Пример:
```
dev-68022509
```

### Процесс

1. Подвинуть в Github Boards задачу "В работе". Запомнить номер.
2. Создать ветку от `master`:
   ```bash
   git checkout master
   git pull
   git checkout -b dev-<номер>
   ```
3. Внести изменения, коммитить с описанием.
4. Отправить ветку и создать Merge Request в `master`.
5. После ревью и CI — мерджнуть и удалить ветку.

## Структура

```text
.
├── backend/
│   ├── cmd/app/          # точка входа
│   └── internal/         # внутренняя логика приложения
└── frontend/
    ├── public/           # статические файлы
    └── src/
        ├── app/          # Next.js App Router
        ├── entities/     # бизнес-сущности
        ├── features/     # пользовательские сценарии
        ├── widgets/      # составные UI-блоки
        └── shared/       # общие API, config, lib и UI
```

## Требования

- Go 1.26+
- Node.js 20.9+
- npm 10+

## Доска задач

[Backend](https://github.com/users/kaaarrrrrriii-ui/projects/5/views/1)

[Frontend](https://github.com/users/kaaarrrrrriii-ui/projects/4)


## Запуск frontend и backend в двух отдельных терминалах.

### Backend

```
cd backend
go mod download
go run ./cmd/app
```
доступно на ``` http://localhost:8080 ```

### Frontend

```
cd frontend
npm install
npm run dev
```
доступно на ``` http://localhost:3000 ```

Настройки портов и адрес API находятся в .env

```
PORT=8080
NEXT_PUBLIC_API_URL=http://localhost:8080```
