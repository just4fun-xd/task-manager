# 📝 Task Manager API

<p align="center">
  <img src="https://img.shields.io/github/actions/workflow/status/just4fun-xd/task-manager/deploy.yml?branch=main&style=for-the-badge&logo=github-actions&label=Deploy%20Status" alt="Deploy Status">
  <img src="https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version">
  <img src="https://img.shields.io/badge/Server%20Status-Online-success?style=for-the-badge&logo=raspberry-pi&logoColor=white" alt="Server Status">
</p>

<p align="center">
  <b>Сервис управления задачами, развернутый на Raspberry Pi 🍓</b>
  <br />
  <a href="http://5.35.45.233:8000/tasks"><strong>Попробовать Live Demo »</strong></a>
</p>

---

## 🛠 Стек технологий

<p align="center">
  <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white" />
  <img src="https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white" />
  <img src="https://img.shields.io/badge/GitHub%20Actions-2088FF?style=for-the-badge&logo=github-actions&logoColor=white" />
  <img src="https://img.shields.io/badge/Raspberry%20Pi-C51A4A?style=for-the-badge&logo=raspberry-pi&logoColor=white" />
</p>

## ✨ Особенности проекта
- **Слоистая архитектура (Clean Architecture):** Четкое разделение на API, Service и Repository уровни.
- **Бизнес-валидация:** Строгий контроль состояний задач (нельзя перепрыгнуть из `New` сразу в `Done`).
- **Группировка:** Поддержка создания групп и привязки задач к ним.
- **CI/CD Pipeline:** Автоматическая сборка под **ARMv7** и деплой на Raspberry Pi через SSH.
- **Миграции:** Автоматическое управление схемой базы данных при старте.

## 📁 Структура проекта
```text
.
├── cmd/task-manager/      # Точка входа в приложение
├── internal/
│   ├── api/               # HTTP хэндлеры и роутинг (chi)
│   ├── config/            # Загрузка конфигурации из .env
│   └── task/              # Бизнес-логика, модели и репозитории
├── migrations/            # SQL файлы миграций
├── docker-compose.prod.yaml # Конфигурация для Raspberry Pi
└── dockerfile             # Сборка Go-приложения
```

---

## 📡 API Endpoints

### 📋 Задачи (`/tasks`)

| Метод  | Эндпоинт       | Описание                                           |
| :----- | :------------- | :------------------------------------------------- |
| **GET**| `/tasks`       | Получить все задачи. Фильтр: `?group_id=1`         |
| **POST**| `/tasks`      | Создать новую задачу                               |
| **GET**| `/tasks/{id}`  | Детальная информация о задаче                      |
| **PUT**| `/tasks/{id}`  | Обновить задачу (с валидацией переходов статуса)   |
| **DELETE**| `/tasks/{id}`| Удалить задачу                                     |

### 📁 Группы (`/groups`)

| Метод  | Эндпоинт    | Описание                                  |
| :----- | :---------- | :---------------------------------------- |
| **GET**| `/groups`   | Список всех групп                         |
| **POST**| `/groups`  | Создать новую группу (категорию)          |

## 📄 Примеры запросов (JSON)

### Создание задачи (`POST /tasks`)

    ```json
    {
    "name": "Купить малинку",
    "description": "Нужна еще одна Raspberry Pi для кластера",
    "group_id": 1
    }
    ```

### Обновление задачи (PUT /tasks/{id})

Доступные статусы: new, in_progress, done.

    ```json
    {
    "name": "Купить малинку",
    "description": "Заказал на маркетплейсе",
    "status": "in_progress",
    "group_id": 1
    }
    ```

### Создание группы (POST /groups)

    ```json
    {
    "name": "Покупки"
    }
    ```

## ⚙️ Установка и запуск

### 💻 Локальная разработка

1. **Склонируйте репозиторий:**
   ```bash
   git clone [https://github.com/just4fun-xd/task-manager.git](https://github.com/just4fun-xd/task-manager.git)
   cd task-manager
   ```

2. **Настройте переменные окружения:**
    ```bash
    cp .env.example .env
    ```
    Отредактируйте .env, указав свои данные для БД.

3. **Запустите контейнеры:**
    ```bash
    docker-compose up -d
    ```

## 🧪 Тестирование

    Для запуска модульных тестов (сервисный слой и логика статусов):
    ```bash
    go test ./internal/task/...
    ```

## 📝 Лицензия
Проект распространяется под лицензией MIT.

