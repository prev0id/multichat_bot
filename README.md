# MultiChat Bot

MultiChatBot - это серверный бот, который пересылает сообщения между чатами YouTube и Twitch, обеспечивая синхронное общение между пользователями обеих платформ.
MultiChatBot подключается к вашим аккаунтам YouTube и Twitch, позволяя участникам обоих чатов видеть и взаимодействовать с сообщениями друг друга в режиме реального времени.


Использованные технологии: Go, HTMX, Tailwind CSS, SQLite3, Docker.

## Запуск

### Локально

Необходим go версии не ниже 1.22.0, утилита make и заполненный ./configs/local.json конфиг.

```shell
make run
```

Перейдите http://localhost:7000 (или другой порт, который указали в конфиг файле)

### В docker контейнере

Потребуется ./configs/prod.json конфиг.

```shell
docker-compose up -d
```

Перейдите https://localhost/

## Использованные зависимости

### Go пакеты

Авторизация:

+ github.com/dghubble/gologin/v2
+ github.com/dghubble/sessions
+ golang.org/x/oauth2

Роутер:

+ github.com/go-chi/chi/v5

Работа с базой данный:

+ modernc.org/sqlite
+ github.com/doug-martin/goqu/v9

YouTube/Google:

+ google.golang.org/api

Twitch:

+ github.com/gempir/go-twitch-irc/v4

Hot reload:

+ github.com/cosmtrek/air

Linter:

+ github.com/golangci/golangci-lint

### Frontend

+ [npm](https://www.npmjs.com/)
+ [HTMX](https://htmx.org/)
+ [Tailwind CSS](https://tailwindcss.com/)

### SVGs

+ https://www.svgrepo.com/svg/511185/user-02
+ https://www.svgrepo.com/svg/510970/external-link
+ https://www.svgrepo.com/svg/511122/settings
+ https://www.svgrepo.com/svg/510988/folder-code
+ https://www.svgrepo.com/svg/488713/twitch
+ https://www.svgrepo.com/svg/488595/google
+ https://www.svgrepo.com/svg/463649/settings
