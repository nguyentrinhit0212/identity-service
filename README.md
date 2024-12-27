# Identity Service

A secure and scalable identity management service built in Go, featuring OAuth2 integration, multi-tenant support, and advanced security features.

## Features

- **Authentication**
  - JWT-based authentication with RSA key rotation
  - OAuth2 support (Google, with extensible provider system)
  - Session management with refresh tokens
  - Multi-factor authentication (MFA/2FA)

- **User Management**
  - User profile management
  - Password management with secure hashing
  - Email verification
  - Password reset functionality

- **Multi-tenancy**
  - Tenant-based access control
  - Tenant switching capability
  - Per-tenant user settings

- **Security**
  - RSA key rotation for JWT signing
  - Secure session management
  - Rate limiting
  - MFA/2FA support
  - Security settings management

## Architecture

### Core Components

- **Services**
  - `AuthService`: Handles authentication and token management
  - `UserService`: Manages user operations
  - `TenantService`: Handles multi-tenant functionality
  - `SecurityService`: Manages security settings and MFA

- **Repositories**
  - Session management
  - User data storage
  - Tenant information
  - Security settings

- **OAuth Providers**
  - Extensible OAuth provider system
  - Built-in Google OAuth support
  - Easy integration for additional providers

### Security Features

- JWT token management with RSA key pairs
- Automatic key rotation (default: 24 hours)
- Secure session handling
- Password hashing using bcrypt

## Setup

1. Clone the repository
2. Configure environment variables
3. Set up the database
4. Run the service

## Configuration

Key configuration parameters:

```go
const (
    defaultTokenDuration     = 24 * time.Hour
    defaultKeySize          = 2048
    defaultKeyRotationPeriod = 24 * time.Hour
)
```

## API Endpoints

### Authentication
- `POST /auth/login`: User login
- `POST /auth/logout`: User logout
- `POST /auth/refresh`: Refresh access token
- `GET /auth/session`: Get current session

### OAuth
- `GET /auth/oauth/{provider}`: Initiate OAuth flow
- `GET /auth/oauth/{provider}/callback`: OAuth callback

### User Management
- `GET /users/profile`: Get user profile
- `PUT /users/profile`: Update user profile
- `PUT /users/password`: Update password

### Security
- `GET /security/settings`: Get security settings
- `PUT /security/settings`: Update security settings
- `POST /security/mfa/enable`: Enable MFA
- `POST /security/mfa/disable`: Disable MFA
- `POST /security/mfa/verify`: Verify MFA token

## Development

### Key Components

1. **JWT Key Manager**
   - Handles RSA key pair generation and rotation
   - Manages signing and verification of tokens
   - Maintains current and previous keys for smooth rotation

2. **OAuth System**
   - Extensible provider interface
   - Built-in Google implementation
   - Easy to add new providers

3. **Session Management**
   - Secure session tracking
   - Automatic session cleanup
   - Multiple device support

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

[Add your license here]
```
