# API Development

## Введение

(Слайд №1\introduction)

```
Привет всем заглянувшим!  
Ролик будет без голоса — только заметки и действия.
```

(Слайд №2\introduction)

В этом видео мы создам REST API, позволяющее работать с сущностью «видео».
И в итоге получим следующий набор маршрутов:

| Method | URL Pattern                   | Action                                              |
|--------|-------------------------------|-----------------------------------------------------|
| GET    | /v1/healthcheck               | Показать информацию о состоянии и версии приложения |
| GET    | /v1/movies                    | Показать информацию обо всех фильмах                |
| POST   | /v1/movies                    | Создать новый фильм                                 |
| GET    | /v1/movies/:id                | Показать информацию о конкретном фильме             |
| PATCH  | /v1/movies/:id                | Обновить информацию о конкретном фильме             |
| DELETE | /v1/movies/:id                | Удалить конкретный фильм                            |
| POST   | /v1/users                     | Зарегистрировать нового пользователя                |
| PUT    | /v1/users/activated           | Активировать конкретного пользователя               |
| PUT    | /v1/users/password            | Обновить пароль конкретного пользователя            |
| POST   | /v1/tokens/authentication     | Создать новый токен аутентификации                  |
| POST   | /v1/tokens/password-reset     | Создать новый токен для сброса пароля               |
| GET    | /debug/vars                   | Показать метрики приложения                         |


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

```bash
go run ./cmd/api
curl -i localhost:4000/v1/healthcheck
```

### Слайды

(Слайд №5/getting-started)

Мы начали с использования маршрутизатора ServeMux из пакета net/http.
Но в дальнейшем будем использовать пакет httprouter.

Причины:

- По умолчанию http.ServeMux возвращает ответы с ошибкой 404 Not Found
  с заголовком:
  Content-Type: text/plain; charset=utf-8

  Это поведение можно переопределить через middleware
  или собственный обработчик ошибок, чтобы возвращать JSON-ответы:
  Content-Type: application/json

- http.ServeMux не предоставляет встроенной обработки
  405 Method Not Allowed.

- http.ServeMux не поддерживает автоматическую обработку
  OPTIONS-запросов и формирование заголовка Allow.


(Слайд №6/getting-started)

План:
1. Заменить маршрутизатор.
2. Создать два новых эндпоинта.

| Method | URL Pattern      | Handler            | Action                         |
|--------|------------------|--------------------|----------------------------------|
| GET    | /v1/healthcheck  | healthcheckHandler | Показать информацию о приложении |
| POST   | /v1/movies       | createMovieHandler | Создать новый фильм              |
| GET    | /v1/movies/:id   | showMovieHandler   | Показать конкретный фильм        |

3. Создать дополнительный метод для чтения параметра id из URL-адреса.
4. Выполнить тестовые запросы.

### Действия

touch cmd/api/routes.go
```go
package main

import (
    "net/http"

    "github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
    router := httprouter.New()

    router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
    router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
    router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)

    return router
}
```

DELETE in cmd/api/main.go
```go
mux := http.NewServeMux()
mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)
```

UPDATE in cmd/api/main.go
```go
srv := &http.Server{
    Addr:         fmt.Sprintf(":%d", cfg.port),
    Handler:      app.routes(), // <-
    IdleTimeout:  time.Minute,
    ReadTimeout:  5 * time.Second,
    WriteTimeout: 10 * time.Second,
    ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
}
```

```bash
 touch cmd/api/movies.go
 ```

 cmd/api/movies.go
 ```go
 package main

import (
    "fmt"
    "net/http"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "create a new movie")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
    id, err := app.readIDParam(r)
    if err != nil {
        http.NotFound(w, r)
        return
    }

    fmt.Fprintf(w, "show the details of movie %d\n", id)
}
```

```bash
touch cmd/api/helpers.go
```

cmd/api/helpers.go
```go
package main

import (
    "errors"
    "net/http"
    "strconv"

    "github.com/julienschmidt/httprouter"
)
 
func (app *application) readIDParam(r *http.Request) (int, error) {
    params := httprouter.ParamsFromContext(r.Context())

    id, err := strconv.Atoi(params.ByName("id"))
    if err != nil || id < 1 {
        return 0, errors.New("invalid id parameter")
    }

    return id, nil
}
```

REBUILD
```bash
go run ./cmd/api
```

TEST
```bash
curl localhost:4000/v1/healthcheck
curl -X POST localhost:4000/v1/movies
curl localhost:4000/v1/movies/123
curl -i -X ​​POST localhost:4000/v1/healthcheck
curl -i -X ​​OPTIONS localhost:4000/v1/healthcheck
curl -i localhost:4000/v1/movies/abc
```

## Отправка JSON-ответов

### Слайды

**(Слайд №1 sending-json-responses)**  
Следующим шагом поработаем с ответами нашего API.  
Добавим поддержку ответов в формате JSON.

**(Слайд №2 sending-json-responses)**  
JSON — это самостоятельный формат обмена данными.  
У него есть собственные типы:

- object
- array
- string
- number
- true
- false
- null

```json
{
  "name": "Alice",
  "age": 25,
  "admin": true,
  "tags": ["go", "linux"],
  "profile": {
    "country": "UA"
  },
  "deleted": null
}
```

**(Слайд №3 sending-json-responses)**  
В Go есть пакет `encoding/json`.

По сути, это двусторонняя система преобразования::

- Go → JSON
- JSON → Go

**(Слайд №4 sending-json-responses)**  
Правила при направлении Go → JSON:

```go
jsonData, err := json.Marshal(user)
```

- `struct` → JSON object
- `string` → string
- `[]byte` → JSON-строка, закодированная в Base64
- `int`, `float64` → number
- `bool` → true/false
- `slice` → array
- `map` → JSON object
- `pointer` → pointer → если nil, то null, иначе сериализуется значение, 
на которое указывает pointer по общим правилам
- `chan`, `func`, `complex*` → не поддерживаются
- `time.Time` → JSON-строка в формате RFC3339

---

**(Слайд №5 sending-json-responses)**  
Правила при направлении JSON → Go:

```go
json.Unmarshal(jsonStr, &target)
```

Программист должен предоставить подходящую структуру под полученный JSON.  
После этого будут применяться правила, обратные вышеописанным:

- JSON object → `struct` или `map`
- JSON array → `slice` или `array`
- JSON string → `string`
- JSON boolean → `bool`
- JSON number → числовой Go-тип
- JSON null → nil-able Go-тип

**(Слайд №6 sending-json-responses)**  
Когда Go кодирует определённый тип в JSON, он проверяет, есть ли у этого типа метод `MarshalJSON()`.

Если метод существует, Go вызывает его, чтобы определить способ кодирования.

То есть, если нам нужен собственный способ кодирования, мы можем:

- создать свой тип
- реализовать для него интерфейс:

```go
type Marshaler interface {
    MarshalJSON() ([]byte, error)
}
```

Если тип удовлетворяет интерфейсу, Go вызовет метод `MarshalJSON()` и использует возвращённый `[]byte` как готовое JSON-значение.

Но в этом видео я этого делать не буду.

---

**(Слайд №7 sending-json-responses)**  
План:

1. Создать вспомогательный метод writeJSON в cmd/api/helpers.go.

   Он будет:
   - создавать JSON,
   - устанавливать нужные заголовки,
   - записывать JSON в тело ответа.

2. Создать структуру Movie в internal/data, которая будет преобразовываться в JSON.

Добавим теги структуры, чтобы:
- изменить имена полей в JSON-объекте;
- скрыть некоторые поля.

3. Внести изменения в обработчики, чтобы они работали с новым методом writeJSON.

4. Добавить систему обработки ошибок, чтобы все ошибки возвращались в JSON-формате.

### Действия

cmd/api/helpers.go
ADD  
"encoding/json" // New import
```go
type envelope map[string]any

func (app *application) writeJSON(
    w http.ResponseWriter, 
    status int, 
    data envelope, 
    headers http.Header,
    ) error {
    js, err := json.Marshal(data)
    if err != nil {
        return err
    }

    js = append(js, '\n')

    for key, values := range headers {
        for _, value := range values {
            w.Header().Add(key, value)
        }
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    w.Write(js)

    return nil
}
```

```bash
mkdir internal/data && touch internal/data/movies.go
```

internal/data/movies.go
```go
package data

import (
    "time"
)

type Movie struct {
    ID        int       `json:"id"`
    CreatedAt time.Time `json:"-"` 
    Title     string    `json:"title"`
    Year      int       `json:"year,omitzero"`    
    Runtime   int       `json:"runtime,omitzero"` 
    Genres    []string  `json:"genres,omitzero"`  
    Version   int       `json:"version"`
}
```

### Слайды

**(Слайд №8 sending-json-responses)**  
С помощью тегов вида `json:"id"` мы определяем, как имена полей будут записаны в JSON.

Атрибут `omitzero` позволяет игнорировать поле, если оно содержит значение по умолчанию (нулевое значение).

Атрибут `"-"` безусловно исключает поле из JSON. 

### Действия

cmd/api/healthcheck.go

DELETE 
- "fmt"
```go
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
    envelope := envelope{
        "status": "available",
        "system_info": map[string]string{
            "environment": app.config.env,
            "version":     version,
        },
    }

    err := app.writeJSON(w, http.StatusOK, envelope, nil)
    if err != nil {
        app.serverErrorResponse(w, r, err)
    }
}
```

File: cmd/api/movies.go

ADD
- "time" // New import
- "github.com/vysmv/greenlight/internal/data" // New import
```go
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
    id, err := app.readIDParam(r)
    if err != nil {
        app.notFoundResponse(w, r)
        return
    }

    movie := data.Movie{
        ID:        id,
        CreatedAt: time.Now(),
        Title:     "Casablanca",
        Runtime:   102,
        Genres:    []string{"drama", "romance", "war"},
        Version:   1,
    }

    err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
    if err != nil {
        app.serverErrorResponse(w, r, err)
    }
}
```

```bash
touch cmd/api/errors.go
```
File: cmd/api/errors.go
```go
package main

import (
    "fmt"
    "net/http"
)

func (app *application) logError(r *http.Request, err error) {
    var (
        method = r.Method
        uri    = r.URL.RequestURI()
    )

    app.logger.Error(err.Error(), "method", method, "uri", uri)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
    envelope := envelope{"error": message}

    err := app.writeJSON(w, status, envelope, nil)
    if err != nil {
        app.logError(r, err)
        w.WriteHeader(500)
    }
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
    app.logError(r, err)

    message := "the server encountered a problem and could not process your request"
    app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
    message := "the requested resource could not be found"
    app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
    message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
    app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}
```

File: cmd/api/routes.go
```go
func (app *application) routes() http.Handler {
    router := httprouter.New()

    //ADD
    router.NotFound = http.HandlerFunc(app.notFoundResponse)

    // ADD
    router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

    router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
    router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
    router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)

    // ADD
    return app.recoverPanic(router)
}
```

```bash
touch cmd/api/middleware.go
```

File: cmd/api/middleware.go
```go
package main

import (
    "fmt"
    "net/http"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            pv := recover()
            if pv != nil {
                w.Header().Set("Connection", "close")
                app.serverErrorResponse(w, r, fmt.Errorf("%v", pv))
            }
        }()

        next.ServeHTTP(w, r)
    })
}
```