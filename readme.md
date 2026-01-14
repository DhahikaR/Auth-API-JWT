# 🚀 Auth API — Fiber + JWT (Golang)

RESTful Authentication API menggunakan **Golang, Fiber, JWT, dan GORM**, dengan dukungan **PostgreSQL** sebagai database utama dan **SQLite (in-memory)** untuk unit testing, role-based authorization, Swagger documentation, serta unit testing dengan test coverage 81%.

---

## ✨ Fitur Utama

### 🔐 Authentication

- Register user
- Login (JWT generation)
- Verifikasi token via middleware
- Claim & expiry validation
- Struktur response seragam (WebResponse)

### 👤 User Management

- /users/me → lihat & update profile sendiri
- Admin: CRUD seluruh user
- Update user
- Delete user
- Find user by ID / email

### 🛡 Middleware

- JWT Middleware → verifikasi token
- Admin Middleware → batasi akses admin saja
- Ownership Guard → user hanya bisa akses datanya sendiri

### 📘 Swagger Documentation

Swagger otomatis dengan anotasi Go:

http://127.0.0.1:3000/docs/

(Hasil swag init berada di folder docs/)

---

## 📂 Struktur Folder

```
auth-api-jwt/
│
├── config/  # DB setup & config
│   ├── database.go
|
├── controller/ # Fiber controllers (request handlers) + Swagger annotations
│   ├── auth_controller.go
│   ├── user_controller.go
│   ├── auth_controller_docs.go
│   └── user_controller_docs.go
│
├── service/ # Business logic layer
│   ├── auth_service.go
│   ├── auth_service_impl.go
│   ├── user_service.go
│   └── user_service_impl.go
│
├── repository/ # GORM data access layer
│   ├── auth_repository.go
│   ├── auth_repository_impl.go
│   ├── user_repository.go
│   └── user_repository_impl.go
│
├── models/
│   ├── domain/ # Database models
│   └── web/ # Request & response DTOs
│
├── middleware/
│   ├── jwt_middleware.go
│   └── admin_middleware.go
│
├── helper/
│   ├── response.go
│   ├── error.go
│   ├── tx.go
│   ├── validation.go
│   ├── model.go
│   └── validation.go
│
├── utils/
│   ├── hash.go
│   └── jwt.go
│
├── routes/
│   ├── auth_routes.go
│   └── user_routes.go
│
├── docs/ # Generated Swagger files
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
│
├── test/ # Unit tests
│   ├── auth_service_test.go
│   ├── user_service_test.go
│   ├── jwt_middleware_test.go
│   └── ...
│
├── main.go
└── go.mod
```

---

## ⚙️ Instalasi & Setup

1️⃣ Clone Repo

```bash
git clone https://github.com/DhahikaR/Auth-API-JWT.git
cd Auth-API-JWT
```

2️⃣ Install Dependencies

```bash
go mod tidy
```

3️⃣ Setup Environment

Buat file .env:

```bash
JWT_SECRET=your_secret_key

#Super Admin (opsional)
ADMIN_EMAIL=admin@example.com
ADMIN_PASSWORD=12345678
ADMIN_NAME=Super Admin
```

4️⃣ Jalankan server

```bash
go run main.go
```

Server berjalan di:

http://127.0.0.1:3000

---

## 🔥 Swagger Documentation

http://127.0.0.1:3000/docs/

Regenerasi dokumentasi:

```bash
swag init --parseDependency --parseInternal
```

Komentar anotasi berada pada:

```bash
controller/auth_controller_docs.go
controller/user_controller_docs.go
swagger_root.go
```

---

## 🔐 Authentication Flow

- Register
  POST /auth/register

- Login
  POST /auth/login

Response berisi JWT:

- Authorization: Bearer <token>
- Protected routes

Semua endpoint /users membutuhkan token valid.

---

## 👨‍💼 Penjelasan Mekanisme Super Admin

Di main.go terdapat mekanisme opsional untuk auto-seeding Super Admin:

```bash
 func seedSuperAdmin(db *gorm.DB) {
 	email := os.Getenv("ADMIN_EMAIL")
 	password := os.Getenv("ADMIN_PASSWORD")
 	fullName := os.Getenv("ADMIN_NAME")

 	if email == "" || password == "" {
 		log.Println("Super admin environment variables not set. Skipping seeder...")
 		return
 	}

 	var count int64
 	db.Model(&domain.User{}).Where("role = ?", "admin").Count(&count)

 	if count == 0 {
 		hashed, _ := utils.HashPassword(password)

 		db.Create(&domain.User{
 			Email:        email,
 			PasswordHash: hashed,
 			FullName:     fullName,
 			Role:         "admin",
 			IsVerified:   true,
 		})

 		log.Println("Super admin created:", email)
 	} else {
 		log.Println("Super admin already exists. Skipping seeder...")
 	}
 }
```

---

## 🎯 Mengapa Super Admin Dibuat?

- Hanya admin yang boleh membuat user lain (via endpoint /users)
- Tanpa admin awal, API tidak bisa digunakan untuk CRUD user
- Seeder memastikan terdapat minimal 1 admin permanen

🔧 Cara mengaktifkan seeder:

Di main.go:

```bash
seedSuperAdmin(db)
```

---

## 🧑‍⚖️ Role Akses

- user hanya bisa akses /users/me
- admin boleh CRUD semua user

---

## 📌 API Endpoints

### 🔐 Auth

Method Endpoint Deskripsi

- POST /auth/register Register user
- POST /auth/login Login & JWT

### 👤 User

Method Endpoint Role Deskripsi

- GET /users/me user/admin lihat profil sendiri
- PUT /users/me user/admin update profil sendiri
- GET /users/:id admin/user\* user hanya bisa miliknya sendiri
- POST /users admin create user
- PUT /users/:id admin update user
- DELETE /users/:id admin delete user

---

## 🧪 Testing

Jalankan seluruh test:

```bash
go test ./... -v -coverpkg=./...
```

### SQLite digunakan sebagai in-memory database untuk speed & isolation.

Testing mencakup:

- Controller
- Service
- Middleware
- Repository
- Helper
- Exception handling

---

## 🛡 Keamanan

- JWT HS256
- Token expiry
- Password hashing (bcrypt via utils.HashPassword)
- Validasi input struct
- Role-based authorization
- Error standardization (WebResponse)

---

## 🧑‍💻 Author

**Dhahika Rahmadani**  
Backend Developer • Go Enthusiast  
📧 [dhahikardani@gmail.com](mailto:dhahikardani@gmail.com)
