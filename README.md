#  ✨Koda Shortlink
> This is an application to shorten/shorten the URL you entered with a system that redirects to the destination, and here is the dashboard page to see your data for 7 days and you can do link management.


🚀 Features
- 🔐 Authentication (Login & Register)
- 🔗 Shortlink Management
- 📊 Analytics & Tracking
- 🔄 Token-based Access
- 🛡 Security & Access Control
- ⚡ Performance Optimization

## 🛠️ Tech Stack
![Go](https://img.shields.io/badge/-Go-00ADD8?logo=go&logoColor=white&style=for-the-badge)
![Gin](https://img.shields.io/badge/-Gin-00ADD8?logo=go&logoColor=white&style=for-the-badge)
![PostgreSQL](https://img.shields.io/badge/-PostgreSQL-4169E1?logo=postgresql&logoColor=white&style=for-the-badge)
![Swagger](https://img.shields.io/badge/Swagger-UI-85EA2D?logo=swagger&logoColor=black&style=for-the-badge)
![JWT](https://img.shields.io/badge/JWT-000000?logo=jsonwebtokens&logoColor=white&style=for-the-badge)
![Argon2](https://img.shields.io/badge/Argon2-0A7E8C?style=for-the-badge)
![Go Migrate](https://img.shields.io/badge/Go%20Migrate-01B3E3?logo=go&logoColor=white&style=for-the-badge)


##  🔐 .env Configuration
```
# Database
DBUSER=youruser
DBPASS=yourpass
DBHOST=localhost
DBPORT=5432
DBNAME=tickitz

# JWT hash
JWT_SECRET=your_jwt_secret

# Redish
REDISUSER=<redis_user>
REDISPASS=<redis_pass>
REDISPORT=6379
REDISHOST=<redis_host>

# Vercel
DATABASE_URL=<your_url_database>
REDIS_URL=<your_redis_url>

# CORS
CORS_ORIGIN1=http://localhost:5173
CORS_ORIGIN2=<url_frontend_2>

# SHORT BASE URL
BASE_URL=http://localhost:8011
```

## 📦 How to Install & Run Project
### 1. First, clone this repository: 
```
https://github.com/federus1105/koda-b4--final-backend.git
```
### 2. Install Dependencies
```go
go mod tidy
```
### 3. Setup your environment
### 4. Do the Database Migration
### 5. Run Server/Project
```go
go run .\cmd\main.go 
```
### 6. Init Swagger
```go
swag init -g ./cmd/main.go
```
### Open Swagger Documentation in Browser
#### ⚠️ Make sure the server is running
```http://localhost:8011/swagger/index.html```


<br>


## 🗃️ How to run Database Migrations
### ⚠️ Attention: This only applies to PostgreSQL, because enums can only be used in PostgreSQL.
### 1. Install Go migrate
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest;
```
### 2. Create database
```bash
CREATE DATABASE <database_name>;
```
### 3. Migrations Up
```bash
migrate -path ./db/migrations -database "postgres://user:password@localhost:5432/database?sslmode=disable" up
```
### 4. Migrations Down
```bash
migrate -path ./db/migrations -database "postgres://user:password@localhost:5432/database?sslmode=disable" down
```

## 📌 Koda Shortlink - API Endpoints
`http://localhost:8080/api/v1`

| Endpoint                  | Method | Auth             | Keterangan                                  |
|----------------------------|--------|-----------------|--------------------------------------------|
| `/auth/register`           | POST   | No              | Register user                               |
| `/auth/login`              | POST   | No              | Login, get it `access_token` & `refresh_token` |
| `/auth/logout`             | POST   | Bearer Token    | Logout user                                 |
| `/auth/refresh`            | POST   | No              | Refresh `access_token`                      |
| `/dashboard/stats`         | GET    | Bearer Token    | get statistic dashboard                |
| `/profile`                 | GET    | Bearer Token    | get data profile user                  |
| `/links`                   | POST   | Optional        | create shortlink                              |
| `/links`                   | GET    | Bearer Token    | List all shortlink user                   |
| `/links/:shortcode`        | GET    | Bearer Token    | Detail shortlink                            |
| `/links/:shortcode`        | DELETE | Bearer Token    | delete shortlink                             |
| `/:shortcode`              | GET    | No              | Redirect ke URL original                        |

---

## 👨‍💻 Made with by
📫 [federusrudi@gmail.com](mailto:federusrudi@gmail.com)  
💼 [LinkedIn](https://www.linkedin.com/in/federus-rudi/)  

## 📜 License
Released under the **MIT License**.  
You’re free to use, modify, and distribute this project — just don’t forget to give a little credit

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
