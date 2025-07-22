# 🚀 GoFiber Starterkit with MySQL

A ready-to-use starterkit for building REST APIs using [GoFiber](https://gofiber.io) and MySQL.  
It provides a clean architecture structure, validation, logging, and database connection setup to help you focus on your business logic.

---

## 📁 Project Structure

```
├── cmd
│   └── main.go                          # Application entry point
├── config                               # App configurations
│   ├── app.go
│   ├── environment.go
│   ├── fiber.go
│   ├── logrus.go
│   ├── mysql.go
│   └── validator.go
├── Dockerfile                           # Docker setup
├── go.mod
├── go.sum
├── internal                             # Clean architecture layout
│   ├── domain
│   │   └── example
│   │       ├── dto
│   │       ├── entity
│   │       ├── handler
│   │       │   └── example_handler.go
│   │       ├── repository
│   │       │   └── resource
│   │       └── service
│   └── infrastructure
│       └── rest
│           ├── middleware
│           └── routes
│               └── route.go
├── pkg
│   └── response_status_exception.go     # Utility helpers
└── test                                 # Unit tests
```

---

## ✅ Features

- ⚡ [Fiber](https://gofiber.io) framework (fast & lightweight)
- 🐬 ORM [Gorm](https://gorm.io) and MySQL with connection pooling
- 🔍 Input validation with [validator](https://pkg.go.dev/github.com/go-playground/validator/v10)
- 📦 Modular clean architecture
- 🐳 Docker ready
- 📄 Auto JSON response [formatter](https://pkg.go.dev/github.com/goccy/go-json)
- 🧪 Ready for unit testing

---

## ⚙️ Installation & Setup

### 1. Clone the Repository

```bash
git clone https://github.com/nurmanhadi/gofiber-starter-kit-mysql.git
cd gofiber-starter-kit-mysql
```

---

### 2. Create `.env` File

Create a `.env` file in the root directory:

```dotenv
APP_NAME="Your App Name"

DB_MYSQL_URL="username:password@tcp(localhost:3306)/yourdb?charset=utf8mb4&parseTime=True&loc=Local"
DB_POOL_MAX_IDLE_CONNS=5
DB_POOL_MAX_OPEN_CONNS=15
DB_POOL_MAX_IDLE_TIME=10    # in minutes
DB_POOL_MAX_LIFETIME=30     # in minutes
```

---

### 3. Remove Git History & Reinitialize

```bash
rm -rf .git
git init
```

---

### 4. Rename Module in `go.mod`

Edit the `go.mod` file:

**Before:**
```go
module gofiber-starterkit-mysql
```

**After:**
```go
module your-module
```

Then run:

```bash
go mod tidy
```

finnaly you can chnage `import path`

---

## 🐳 Run with Docker

```bash
docker build -t gofiber-app .
docker run -p 3000:3000 --env-file .env gofiber-app
```

---

## 📮 Example Endpoint

`GET /`

---

## 🤝 Contributing

Pull requests are welcome!  
Feel free to fork and customize this starterkit for your own projects.

---

## 📝 License

MIT License. Use, modify, and distribute as needed.

---

## 👤 Author

Muhammad Nurman Hadi  
GitHub: [@nurmanhadi](https://github.com/nurmanhadi)