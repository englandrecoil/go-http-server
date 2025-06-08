# go-http-server
[![en](https://img.shields.io/badge/lang-en-green?style=flat)](https://github.com/englandrecoil/go-http-server/blob/main/README.md)

**go-http-server** — это минималистичный HTTP-бэкенд на языке Go, имитирующий простейшую социальную сеть. Проект задуман как песочница для изучения веб-разработки на Go: обработки HTTP-запросов, построения REST API и работы с JSON — исключительно на стандартной библиотеке.

## ✨ Возможности
- Регистрация, аутентификация и авторизация пользователей  
- Создание, удаление и просмотр постов
- Вебхук стороннего сервиса оплаты
- JWT, refresh-токены  

## 🚀 Необходимые компоненты
Для запуска программы нужно установить **две** вещи:
1. Go версии 1.23 или выше ([официальная инструкция](https://go.dev))
2. PostgreSQL версии 16 или выше ([инструкция по установке здесь](https://www.postgresql.org))

Выполните команды `go version` и `postgres --version` в терминале, чтобы убедиться, что Go и PostgreSQL установлены корректно.

## 💽 Установка
Чтобы начать использовать:
1. Клонируйте репозиторий:
```
git clone https://github.com/englandrecoil/go-http-server
```
2. Перейдите в созданную папку:
```
cd go-http-server
```
3. Создайте .env файл в корне папки go-http-server и не забудьте подставить свои значения. С его структурой и примером можно ознакомиться ниже:
```
# Ссылка на БД
DB_URL="postgres://nikitateresenko:@localhost:5432/chirpy?sslmode=disable"
# Оставьте это значение по умолчанию
PLATFROM="dev"
# "Секрет", который используется при валидации и создании JWT и refresh-токенов
SECRET="NvUIUpWg8dUDClS/h3n2hMd1zJwkVRBqTz57rgJsbWlF5X7ulCnL4CEliEjCMV+4
RbjRvFODK8ZEB/cotg8/AA=="
# Ключ стороннего сервиса оплаты
POLKA_KEY="f271c81ff7084ee5b99a5091b42d486e"
```
4. Соберите и запустите сервер командой:
```
go build && ./go-http-server
```
После успешного запуска сервер будет доступен снаружи как localhost:8080

## 📌 API эндпоинты:
| Метод | Эндпоинт | Описание |
| --- | --- | --- |
| GET | /api/healthz | Просмотр состояния сервера |
| POST | /api/polka/webhooks | Обновление статуса пользователя при покупке подписки |
| GET | /admin/metrics | Получение метрик |
| POST | /admin/reset | Удаление всех пользователей из БД |
| POST | /api/users | Создание пользователя |
| PUT | /api/users | Обновление пользовательских данных |
| POST | /api/login | Аутентификация и авторизация пользователя |
| POST | /api/refresh | Создание refresh-токена |
| POST | /api/revoke | Отмена refresh-токена |
| POST | /api/chirps | Создание поста |
| GET | /api/chirps | Получение всех постов |
| GET | /api/chirps/{chirpID} | Получение конкретного поста по ID |
| DELETE | /api/chirps/{chirpID} | Удаление поста по ID |
