# Music Library

Проект предназначен для управления музыкальной библиотекой, включающей добавление, удаление, обновление и получение данных о песнях.

## Структура проекта

Проект организован с четким разделением ответственности, что упрощает поддержку и расширение функциональности.

### Описание папок и файлов
- **cmd/**:
  - **main.go**: Основной файл приложения.
- **config/**: Настройки конфигурации проекта.
- **docs/**: Документация API.
- **internal/http_server/handlers/**: Обработчики HTTP-запросов.
  - **add_song/**: Обработчик для добавления песни.
  - **delete_song/**: Обработчик для удаления песни.
  - **get_all_data/**: Обработчик для получения всех данных.
  - **get_song/**: Обработчик для получения конкретной песни.
  - **update_song/**: Обработчик для обновления песни.
- **lib/**: Библиотеки и утилиты.
  - **logger/**: Утилиты для логирования.
  - **response/**: Утилиты для формирования ответов.
  - **utils/**: Общие утилиты.
- **mocks/**: Мок-реализация ответа от внешнего API для тестирования.
- **models/**: Модели данных.
- **storage/pg/**: Реализация хранения данных в PostgreSQL.
- **migrations/**: Файлы миграций базы данных.
- **.env**: Файл переменных окружения.
- **.env.example**: Пример файла переменных окружения.
- **.gitignore**: Файл исключений для git.
- **go.mod**: Файл с описанием зависимостей модуля.
- **go.sum**: Файл с контрольными суммами зависимостей.
- **README.md**: Документация и описание проекта.

### Инструкции по установке

1. Склонируйте репозиторий:

   git clone https://github.com/VldslvKtv/music-library.git

2. Перейдите в директорию проекта:

    cd music-library

3. Создайте и настройте файл .env (пример в .env.example).

4. Установите зависимости:

    go mod download

5. Запустите проект:
    
    go run cmd/main.go



