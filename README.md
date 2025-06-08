# go-http-server
[![en](https://img.shields.io/badge/lang-en-green?style=flat)](https://github.com/englandrecoil/go-http-server/blob/main/README.md)
[![ru](https://img.shields.io/badge/lang-ru-blue?style=flat)](https://github.com/englandrecoil/go-http-server/blob/main/README.ru.md)

**go-http-server** is a minimalist HTTP backend written in Go, simulating a basic social network. The project is intended as a sandbox for learning web development in Go — handling HTTP requests, building REST APIs, and working with JSON — using only the standard library.

## ✨ Features
- User registration, authentication, and authorization  
- Creating, deleting, and viewing posts  
- Webhook for a third-party payment service  
- JWT and refresh tokens  

## 🚀 Requirements
To run the program, you need to have **two** things installed:
1. Go version 1.23 or later ([official guide](https://go.dev))
2. PostgreSQL version 16 or later ([installation guide here](https://www.postgresql.org))

Run `go version` and `postgres --version` in your terminal to ensure Go and PostgreSQL are correctly installed.

## 💽 Installation
To get started:
1. Clone the repository:
```
git clone https://github.com/englandrecoil/go-http-server
```
2. Navigate into the created folder:
```
cd go-http-server
```
3. Create a `.env` file in the root of the `go-http-server` folder and make sure to insert your values. Here's the structure and an example:
```
# Database URL
DB_URL="postgres://nikitateresenko:@localhost:5432/chirpy?sslmode=disable"
# Leave this as default
PLATFROM="dev"
# "Secret" used to create and validate JWT and refresh tokens
SECRET="NvUIUpWg8dUDClS/h3n2hMd1zJwkVRBqTz57rgJsbWlF5X7ulCnL4CEliEjCMV+4
RbjRvFODK8ZEB/cotg8/AA=="
# Key from external payment service
POLKA_KEY="f271c81ff7084ee5b99a5091b42d486e"
```
4. Build and run the server with:
```
go build && ./go-http-server
```
After a successful launch, the server will be available at `localhost:8080`.

## 📌 API Endpoints:
| Method | Endpoint | Description |
| --- | --- | --- |
| GET | /api/healthz | Check server health |
| POST | /api/polka/webhooks | Update user status after subscription purchase |
| GET | /admin/metrics | Get server metrics |
| POST | /admin/reset | Delete all users from the database |
| POST | /api/users | Create a new user |
| PUT | /api/users | Update user data |
| POST | /api/login | Authenticate and authorize a user |
| POST | /api/refresh | Generate a refresh token |
| POST | /api/revoke | Revoke a refresh token |
| POST | /api/chirps | Create a new post |
| GET | /api/chirps | Retrieve all posts |
| GET | /api/chirps/{chirpID} | Get a specific post by ID |
| DELETE | /api/chirps/{chirpID} | Delete a post by ID |
