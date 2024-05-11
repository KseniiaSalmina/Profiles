# Profile managment service
Сервис для хранения и управления пользовательскими профилями.

В качестве хранилища используется собственная in-memry база данных.

## Объекты

Данные о пользователях принимаются в следующем виде:

	Email    string `json:"email"`    //required
	Username string `json:"username"` //required
	Password string `json:"password"` //required
	Admin    bool   `json:"admin"`

Профили отдаются в виде:

    ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Admin    bool   `json:"admin"`

При выдаче нескольких профилей сервер отдаёт страницу вида:

    Users       []UserResponse `json:"users"`
	PageNo      int            `json:"page_number"`
	Limit       int            `json:"limit"`
	PagesAmount int            `json:"pages_amount"`

## API
Сервис работает с форматом JSON.

Доступные методы:

    GET /user - возвращает страницу пользователей любому зарегистрированному пользователю. Принимает параметры pageNo и limit, при их отстуствии проставит дефолтные значения (pageNo = 1, limit = 30)
	POST /user - создаёт нового пользователя по запросу любого пользователя с правами администратора, возвращает id (формат uuid)
	GET /user/:id - возвращает профиль конкретного пользователя, доступен для любого зарегистрированного пользователя
	PUT /user/:id - обновляет пользователя по запросу любого пользователя с правами администратора, параметр id обновить нельзя
	DELETE /user/:id - удаляет пользователя по запросу любого пользователя с правами администратора

## Переменные окружения

Сервис умеет считывать переменные из файла .env в директории исполняемого файла (в корне проекта).

В примерах указаны дефолтные значения. Если программа не сможет считать пользовательские env, то возьмет их (предназначены только для тестового запуска).

Переменные сервера:

    SERVER_LISTEN=:8088
    SERVER_READ_TIMEOUT=5s
    SERVER_WRITE_TIMEOUT=5s
    SERVER_IDLE_TIMEOUT=30s

Переменные форматтера (включают в себя данные первого пользователя-администратора):

    Salt          string `env:"FMT_SALT" envDefault:"MyUniqueSalt"`
	AdminUsername string `env:"DB_USERNAME" envDefault:"Admin"`
	AdminPassword string `env:"DB_PASS" envDefault:"qwerty"`
	AdminEmail    string `env:"DB_Email" envDefault:"qwerty@email.com"`

Переменные логгера:

    LOG_LEVEL=debug

## Makefile

Подготовлены следующие команды:

    build - загружает зависимости из go.mod и собирает бинарник
    run - выполняет команду build и запускает приложение
    test - запускает все тесты