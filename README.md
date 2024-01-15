# MultiChat Bot

## Progress

структура репозитория проекта основана
на [Standard Go Project Layout](https://github.com/golang-standards/project-layout)

### Backend

Go 1.21.3.

Дизайн системы:
![system design](./docs/assets/architecture.png)

#### AppManager

#### Twitch

1) Manager. Предоставляет интерфейс для работы с Twitch, как API, так и IRC.
2) Message Service. Предоставляет интерфейс для IRC. Нужен для работы с rate limit'ами.
3) API Manager. Предоставляет интерфейс для работы с API.
4) IRC Client. Websocket клиент принимающий и отправляющий сообщения.
5) Processor. Обрабатывает различные типы сообщения приходящих из IRC Client, нужен для удобной единой точки вызова
   функционала AppManager или twitch.Manager.

#### DataBase

#### Youtube

TODO

### Frontend

+ singe page website
+ [light theme](https://www.realtimecolors.com/dashboard?colors=1c0e03-ffffff-1361a4-d7bff8-197bd2&fonts=Ubuntu-Ubuntu)
+ [dark theme](https://www.realtimecolors.com/dashboard?colors=fceee3-000000-5ba8ec-1f0740-2d90e6&fonts=Ubuntu-Ubuntu)
+ [font](https://fonts.google.com/specimen/Ubuntu)
+ [svg collection](https://www.svgrepo.com/collection/coolicons-line-oval-icons/1)
+ [svg with google/twitch](https://www.svgrepo.com/collection/phosphor-bold-icons/)
+ [logo](https://www.svgrepo.com/svg/324471/robot-artificial-intelligence-android)

### Twitch

+ зарегистрировал бота, разбираюсь с api
+ есть лимит на частоту отправления сообщений ботом

### Youtube

+ регистрировал когда-то давно dev аккаунт в гугловой апишке, надо будет вспоминать как там все работает

### VKPlay

+ похоже публичной апи пока нет. Есть альфа тест, на который можно
  подать [заявление](https://vk.com/wall-212496568_43917)

## Extra

+ в силу жестких ограничений на использовании апи одним аккаунтом бота на твиче, сделать поддержку пула из ботов которые
  будут выступать в качестве единой системы

## TODO

[ ] Парсер itc сообщений твича
[ ] Клиент твича с поддержкой отправки сообщений
[ ] Автополучение токенов, на данный момент они захардкожены в конфиге

[ ] клиент YT

extra планы
[ ] Доп возможности по конфигурированию для клиента твича
[ ] Unit тесты возможно отрефатич большое количество кода)
[ ] базовые e2e тесты для клиентов