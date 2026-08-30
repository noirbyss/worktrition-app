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
- nutrition-service gRPC: `localhost:50052`
- workout-service gRPC: `localhost:50054`
- ai-service gRPC: `localhost:50056`
- gamification-service gRPC: `localhost:50055`

AI через gateway:

- `POST /ai/generations` c JSON `{ "plan_type": "all" | "workout" | "nutrition" }`
- `GET /ai/generations/{generation_id}`

Gamification через gateway:

- `GET /gamification/character`
- `POST /gamification/rewards/workout` c JSON `{ "is_strength": true | false }`
- `POST /gamification/rewards/meal`
- `POST /gamification/rewards/water`

Для остановки:

```powershell
docker compose down
```

Если нужно удалить и volume базы данных тоже:

```powershell
docker compose down -v
```
