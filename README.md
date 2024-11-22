Folder structure

```
identity-service/
├── cmd/
│   └── main.go                # Entry point chính
├── config/
│   └── config.go              # Cấu hình ứng dụng (DB, JWT, v.v.)
├── db/
│   ├── migrations/            # Thư mục chứa các file migration SQL
│   │   ├── 001_create_users.up.sql
│   │   ├── 001_create_users.down.sql
│   │   ├── 002_create_user_credentials.up.sql
│   │   └── 002_create_user_credentials.down.sql
│   └── db.go                  # Kết nối cơ sở dữ liệu
├── internal/
│   ├── auth/
│   │   ├── jwt.go             # Logic JWT (generate, validate, revoke)
│   │   ├── login.go           # Logic đăng nhập
│   │   ├── oauth.go           # Xử lý OAuth2
│   │   └── ip_validation.go   # Xác thực IP
│   ├── models/
│   │   ├── user.go            # Model cho bảng `users`
│   │   ├── user_credentials.go# Model cho bảng `user_credentials`
│   │   ├── oauth_provider.go  # Model cho bảng `oauth_providers`
│   │   ├── jwt_token.go       # Model cho bảng `jwt_tokens`
│   │   ├── ip_whitelist.go    # Model cho bảng `ip_whitelist`
│   │   └── failed_login.go    # Model cho bảng `failed_logins`
│   ├── repositories/
│   │   ├── user_repository.go # CRUD liên quan đến users
│   │   └── token_repository.go# CRUD liên quan đến tokens
│   ├── services/
│   │   ├── auth_service.go    # Service xử lý logic authentication
│   │   ├── token_service.go   # Service xử lý JWT và token management
│   │   └── user_service.go    # Service xử lý logic user
│   └── utils/
│       └── hash.go            # Hàm hash mật khẩu/token
├── Makefile                   # Quản lý build/migration tiện lợi
├── go.mod                     # Module Go
└── go.sum                     # Module dependency
```
