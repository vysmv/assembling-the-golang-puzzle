# API Development

## Введение

(Слайд №1\introduction)

```
Привет всем заглянувшим!  
Ролик будет без голоса — только заметки и действия.
```

(Слайд №2\introduction)

```text
В этом видео мы создам REST API, позволяющее работать с сущностью «видео».
И в итоге получим следующий набор маршрутов:
┌────────┬─────────────────────────────────┬──────────────────────────────────────────────────────┐
│ Method │ URL Pattern                     │ Action                                               │
├────────┼─────────────────────────────────┼──────────────────────────────────────────────────────┤
│ GET    │ /v1/healthcheck                 │ Показать информацию о состоянии и версии приложения  │
│ GET    │ /v1/movies                      │ Показать информацию обо всех фильмах                 │
│ POST   │ /v1/movies                      │ Создать новый фильм                                  │
│ GET    │ /v1/movies/:id                  │ Показать информацию о конкретном фильме              │
│ PATCH  │ /v1/movies/:id                  │ Обновить информацию о конкретном фильме              │
│ DELETE │ /v1/movies/:id                  │ Удалить конкретный фильм                             │
│ POST   │ /v1/users                       │ Зарегистрировать нового пользователя                 │
│ PUT    │ /v1/users/activated             │ Активировать конкретного пользователя                │
│ PUT    │ /v1/users/password              │ Обновить пароль конкретного пользователя             │
│ POST   │ /v1/tokens/authentication       │ Создать новый токен аутентификации                   │
│ POST   │ /v1/tokens/password-reset       │ Создать новый токен для сброса пароля                │
│ GET    │ /debug/vars                     │ Показать метрики приложения                          │
└────────┴─────────────────────────────────┴──────────────────────────────────────────────────────┘
```

(Слайд №3\introduction)

```text
Будут использоваться:
- Go версии 1.26.3
- Редактор VS Code
- Инструмент нагрузочного тестирования hey
- Система контроля версий Git
- Postman/curl
```

## Начало работы

### Слайды

(Слайд №1\getting-started)

```text
Первым делом:
- создадим директорию проекта
- инициализируем модуль
- создадим базовую структуру проекта
- выполним тестовый запуск
```

(Слайд №2\getting-started)

```text
Краткое пояснение структуры проекта:

bin/         — готовые бинарные файлы приложения
cmd/api/     — код API: запуск сервера, HTTP, аутентификация
internal/    — внутренние переиспользуемые пакеты (БД, email, validation)
migrations/  — SQL-миграции базы данных
remote/      — конфиги и скрипты production-сервера
go.mod       — зависимости и модули проекта
Makefile     — автоматизация задач (build, audit, migrations)
```

### Действия
1. 
```bash
mkdir -p $HOME/greenlight && code ./greenlight
go mod init github.com:vysmv/greenlight
mkdir -p bin cmd/api internal migrations remote
touch Makefile
touch cmd/api/main.go
```

2. 
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello world!")
}
```
3. 
```bash
go run ./cmd/api
```

### Слайды

(Слайд №3\getting-started)

Следующим шагом напишем код, необходимый для подготовки и запуска сервера.

Пакет main будет собирать зависимости в структуру application и запускать сервер.

Также создадим первый эндпоинт:

/v1/healthcheck

Этот ресурс будет предоставлять базовую информацию об API:
- status
- environment
- version

(Слайд №4\getting-started)

План:
1. Создадим в файле cmd/api/main.go:
   - version
   - config
   - application
2. Используем пакет flag для возможности изменения значений env и port в конфигурации приложения.
3. Создадим собственный логгер.
4. Установим зависимости в application.
5. Создадим свой ServeMux и зарегистрируем первый маршрут.
6. Создадим healthcheckHandler.
7. Запустим сервер.
8. Выполним тестовый запуск.

Поехали!

### Действия

cmd/api/main.go
```go
package main

import (
    "flag"
    "fmt"
    "log/slog"
    "net/http"
    "os"
    "time"
)

const version = "1.0.0"

type config struct {
    port int
    env  string
}

type application struct {
    config config
    logger *slog.Logger
}

func main() {
    var cfg config

    flag.IntVar(&cfg.port, "port", 4000, "API server port")
    flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
    flag.Parse()

    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

    app := &application{
        config: cfg,
        logger: logger,
    }

    mux := http.NewServeMux()
    mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)

    srv := &http.Server{
        Addr:         fmt.Sprintf(":%d", cfg.port),
        Handler:      mux,
        IdleTimeout:  time.Minute,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
    }

    logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)
    
    err := srv.ListenAndServe()
    logger.Error(err.Error())
    os.Exit(1)
}
```

```bash
touch cmd/api/healthcheck.go
```

cmd/api/healthcheck.go
```go
package main

import (
    "fmt"
    "net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "status: available")
    fmt.Fprintf(w, "environment: %s\n", app.config.env)
    fmt.Fprintf(w, "version: %s\n", version)
}
```

