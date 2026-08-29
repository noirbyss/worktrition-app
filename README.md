# Worktrition

Worktrition - приложение, которое составит тебе персональный план тренировок и питания.

## Запуск через Docker

Требуется установленный Docker с Docker Compose.

Из корня проекта выполните:

```powershell
docker compose up --build
```

После запуска будут доступны:

- фронтенд: `http://localhost:5173`
- API gateway: `http://localhost:8080`
- user-service gRPC: `localhost:50051`

Для остановки:

```powershell
docker compose down
```

Если нужно удалить и volume базы данных тоже:

```powershell
docker compose down -v
```