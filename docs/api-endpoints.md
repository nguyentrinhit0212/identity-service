# Identity Service API Documentation

This document outlines all available API endpoints in the Identity Service.

## Authentication

All authenticated endpoints require a Bearer token in the Authorization header:
```
Authorization: Bearer <access_token>
```

### Error Responses
All endpoints may return these error responses:
```typescript
{
  success: false,
  error: {
    message: string;
    details?: string;
  }
}
```

Common error codes:
- `unauthorized`: Missing or invalid authentication
- `forbidden`: Insufficient permissions
- `invalid_request`: Invalid request parameters
- `not_found`: Resource not found
- `internal_error`: Server error

## Authentication Endpoints

### OAuth Authentication

#### GET /api/auth/:provider/login
Initiates OAuth login flow for specified provider. Currently supported providers: `google`

Query Parameters:
```typescript
{
  redirectUrl: string; // URL to redirect after OAuth completion
}
```

Example Request:
```
GET /api/auth/google/login?redirectUrl=http://localhost:3000/oauth/callback
```

Example Response:
```json
{
  "url": "https://accounts.google.com/o/oauth2/v2/auth?client_id=...&redirect_uri=..."
}
```

#### GET /api/auth/:provider/callback
OAuth callback endpoint that handles the provider's response.

Query Parameters:
```typescript
{
  code: string;    // OAuth authorization code
  state: string;   // State parameter for security
}
```

Example Response (redirects to frontend with tokens):
```
302 Redirect to: http://localhost:3000?access_token=...&refresh_token=...
```

### Traditional Authentication

#### POST /api/auth/login
User login with email/password credentials.

Request:
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

Success Response (200 OK):
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "userId": "123e4567-e89b-12d3-a456-426614174000",
  "tenantId": "123e4567-e89b-12d3-a456-426614174001",
  "accessToken": "eyJhbGciOiJSUzI1NiIs...",
  "refreshToken": "eyJhbGciOiJSUzI1NiIs...",
  "expiresAt": "2023-01-01T00:00:00Z",
  "createdAt": "2023-01-01T00:00:00Z",
  "lastUsedAt": "2023-01-01T00:00:00Z",
  "ipAddress": "192.168.1.1",
  "userAgent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"
}
```

Error Response (401 Unauthorized):
```json
{
  "error": "Invalid credentials",
  "code": "unauthorized"
}
```

#### POST /api/auth/logout
Logout current session. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Success Response (200 OK):
```json
{
  "message": "Logged out successfully"
}
```

#### POST /api/auth/refresh
Refresh access token using refresh token.

Headers:
```
Refresh-Token: <refresh_token>
```

Success Response (200 OK):
```json
{
  "accessToken": "eyJhbGciOiJSUzI1NiIs...",
  "refreshToken": "eyJhbGciOiJSUzI1NiIs...",
  "expiresIn": 900,
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "user@example.com",
    "name": "John Doe",
    "status": "active",
    "role": "user",
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-01-01T00:00:00Z"
  }
}
```

Error Response (401 Unauthorized):
```json
{
  "error": "Invalid refresh token",
  "code": "unauthorized"
}
```

#### GET /api/auth/session
Get current session info. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Success Response (200 OK):
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "userId": "123e4567-e89b-12d3-a456-426614174000",
  "tenantId": "123e4567-e89b-12d3-a456-426614174001",
  "accessToken": "eyJhbGciOiJSUzI1NiIs...",
  "refreshToken": "eyJhbGciOiJSUzI1NiIs...",
  "expiresAt": "2023-01-01T00:00:00Z",
  "createdAt": "2023-01-01T00:00:00Z",
  "lastUsedAt": "2023-01-01T00:00:00Z",
  "ipAddress": "192.168.1.1",
  "userAgent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"
}
```

### Password Management

#### POST /api/auth/forgot-password
Request password reset email.

Request:
```json
{
  "email": "user@example.com"
}
```

Success Response (200 OK):
```json
{
  "message": "Password reset email sent"
}
```

Note: Always returns success even if email doesn't exist (security best practice).

#### POST /api/auth/reset-password
Reset password using reset token.

Request:
```json
{
  "token": "eyJhbGciOiJSUzI1NiIs...",
  "newPassword": "newSecurePassword123"
}
```

Success Response (200 OK):
```json
{
  "message": "Password reset successfully"
}
```

Error Response (400 Bad Request):
```json
{
  "error": "Invalid or expired token",
  "code": "invalid_token"
}
```

#### POST /api/auth/verify-email
Verify email address using verification token.

Request:
```json
{
  "token": "eyJhbGciOiJSUzI1NiIs..."
}
```

Success Response (200 OK):
```json
{
  "message": "Email verified successfully"
}
```

### Protected Session Management

#### GET /api/auth/sessions
List all active sessions. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Success Response (200 OK):
```json
{
  "success": true,
  "data": {
    "sessions": [
      {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "userId": "123e4567-e89b-12d3-a456-426614174000",
        "tenantId": "123e4567-e89b-12d3-a456-426614174001",
        "accessToken": "eyJhbGciOiJSUzI1NiIs...",
        "refreshToken": "eyJhbGciOiJSUzI1NiIs...",
        "expiresAt": "2023-01-01T00:00:00Z",
        "createdAt": "2023-01-01T00:00:00Z",
        "lastUsedAt": "2023-01-01T00:00:00Z",
        "ipAddress": "192.168.1.1",
        "userAgent": "Mozilla/5.0..."
      }
    ]
  }
}
```

#### DELETE /api/auth/sessions/:id
Revoke specific session. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Path Parameters:
- `id`: Session ID to revoke

Success Response (200 OK):
```json
{
  "message": "Session revoked successfully"
}
```

#### DELETE /api/auth/sessions
Revoke all sessions except current. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Success Response (200 OK):
```json
{
  "message": "All sessions revoked successfully"
}
```

### Security Settings

#### GET /api/auth/security
Get security settings. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Success Response (200 OK):
```json
{
  "success": true,
  "data": {
    "mfaEnabled": true,
    "ipWhitelistEnabled": false,
    "passwordExpiryDays": 90,
    "sessionTimeoutMinutes": 15,
    "lastLogin": "2023-01-01T00:00:00Z"
  }
}
```

#### PUT /api/auth/security
Update security settings. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Request:
```json
{
  "mfaEnabled": true,
  "ipWhitelistEnabled": true,
  "passwordExpiryDays": 60,
  "sessionTimeoutMinutes": 30
}
```

Success Response (200 OK):
```json
{
  "message": "Security settings updated successfully"
}
```

### Multi-Factor Authentication (MFA)

#### POST /api/auth/mfa/enable
Enable MFA. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Success Response (200 OK):
```json
{
  "success": true,
  "data": {
    "mfaSecret": "ABCDEFGHIJKLMNOP",
    "qrCode": "data:image/png;base64,..."
  }
}
```

#### POST /api/auth/mfa/disable
Disable MFA. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Request:
```json
{
  "password": "currentPassword123"
}
```

Success Response (200 OK):
```json
{
  "message": "MFA disabled successfully"
}
```

#### POST /api/auth/mfa/verify
Verify MFA token. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Request:
```json
{
  "token": "123456"
}
```

Success Response (200 OK):
```json
{
  "message": "MFA token verified successfully"
}
```

Error Response (401 Unauthorized):
```json
{
  "error": "Invalid MFA token",
  "code": "invalid_mfa_token"
}
```

## Security Endpoints

### IP Whitelist Management

#### GET /api/security/whitelist
List all whitelisted IPs. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Success Response (200 OK):
```json
{
  "ips": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "ip": "192.168.1.1",
      "description": "Office IP",
      "createdAt": "2023-01-01T00:00:00Z"
    },
    {
      "id": "123e4567-e89b-12d3-a456-426614174001",
      "ip": "192.168.1.2",
      "description": "Home IP",
      "createdAt": "2023-01-01T00:00:00Z"
    }
  ]
}
```

#### POST /api/security/whitelist
Add a new IP to the whitelist. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Request:
```json
{
  "ip": "192.168.1.3",
  "description": "New IP"
}
```

Success Response (201 Created):
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174002",
  "ip": "192.168.1.3",
  "description": "New IP",
  "createdAt": "2023-01-01T00:00:00Z"
}
```

#### DELETE /api/security/whitelist/:id
Remove an IP from the whitelist. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Path Parameters:
- `id`: ID of the IP to remove

Success Response (200 OK):
```json
{
  "message": "IP removed successfully"
}
```

### API Keys Management

#### GET /api/security/api-keys
List all API keys. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Success Response (200 OK):
```json
{
  "keys": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "name": "Key 1",
      "prefix": "abc123",
      "expiresAt": "2023-01-01T00:00:00Z",
      "createdAt": "2023-01-01T00:00:00Z",
      "lastUsedAt": "2023-01-01T00:00:00Z"
    },
    {
      "id": "123e4567-e89b-12d3-a456-426614174001",
      "name": "Key 2",
      "prefix": "def456",
      "expiresAt": "2023-01-01T00:00:00Z",
      "createdAt": "2023-01-01T00:00:00Z",
      "lastUsedAt": "2023-01-01T00:00:00Z"
    }
  ]
}
```

#### POST /api/security/api-keys
Create a new API key. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Request:
```json
{
  "name": "New API Key",
  "expiresAt": "2023-12-31T23:59:59Z"
}
```

Success Response (201 Created):
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174002",
  "name": "New API Key",
  "prefix": "ghi789",
  "key": "ghi789...",
  "expiresAt": "2023-12-31T23:59:59Z",
  "createdAt": "2023-01-01T00:00:00Z"
}
```

#### GET /api/security/api-keys/:id
Get details of a specific API key. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Path Parameters:
- `id`: ID of the API key

Success Response (200 OK):
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "Key 1",
  "prefix": "abc123",
  "expiresAt": "2023-01-01T00:00:00Z",
  "createdAt": "2023-01-01T00:00:00Z",
  "lastUsedAt": "2023-01-01T00:00:00Z"
}
```

#### PUT /api/security/api-keys/:id
Update an existing API key. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Path Parameters:
- `id`: ID of the API key

Request:
```json
{
  "name": "Updated API Key",
  "expiresAt": "2024-12-31T23:59:59Z"
}
```

Success Response (200 OK):
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "Updated API Key",
  "prefix": "abc123",
  "expiresAt": "2024-12-31T23:59:59Z",
  "createdAt": "2023-01-01T00:00:00Z",
  "lastUsedAt": "2023-01-01T00:00:00Z"
}
```

#### DELETE /api/security/api-keys/:id
Delete an API key. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Path Parameters:
- `id`: ID of the API key

Success Response (200 OK):
```json
{
  "message": "API key deleted successfully"
}
```

### Audit Logs

#### GET /api/security/audit-logs
List all audit logs. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Success Response (200 OK):
```json
{
  "logs": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "action": "login",
      "userId": "123e4567-e89b-12d3-a456-426614174000",
      "userEmail": "user@example.com",
      "ipAddress": "192.168.1.1",
      "userAgent": "Mozilla/5.0",
      "metadata": {},
      "createdAt": "2023-01-01T00:00:00Z"
    },
    {
      "id": "123e4567-e89b-12d3-a456-426614174001",
      "action": "logout",
      "userId": "123e4567-e89b-12d3-a456-426614174000",
      "userEmail": "user@example.com",
      "ipAddress": "192.168.1.2",
      "userAgent": "Mozilla/5.0",
      "metadata": {},
      "createdAt": "2023-01-01T00:00:00Z"
    }
  ]
}
```

#### GET /api/security/audit-logs/:id
Get details of a specific audit log entry. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Path Parameters:
- `id`: ID of the audit log entry

Success Response (200 OK):
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "action": "login",
  "userId": "123e4567-e89b-12d3-a456-426614174000",
  "userEmail": "user@example.com",
  "ipAddress": "192.168.1.1",
  "userAgent": "Mozilla/5.0",
  "metadata": {},
  "createdAt": "2023-01-01T00:00:00Z"
}
```

### Security Policies

#### GET /api/security/policies
List all security policies. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Success Response (200 OK):
```json
{
  "policies": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "name": "Password Policy",
      "description": "Minimum 8 characters, at least one uppercase letter",
      "createdAt": "2023-01-01T00:00:00Z"
    },
    {
      "id": "123e4567-e89b-12d3-a456-426614174001",
      "name": "Session Timeout",
      "description": "Sessions expire after 15 minutes of inactivity",
      "createdAt": "2023-01-01T00:00:00Z"
    }
  ]
}
```

#### PUT /api/security/policies
Update security policies. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Request:
```json
{
  "policies": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "name": "Password Policy",
      "description": "Minimum 10 characters, at least one uppercase letter, one number",
    }
  ]
}
```

Success Response (200 OK):
```json
{
  "message": "Security policies updated successfully"
}
```

#### POST /api/security/policies/test
Test a security policy. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Request:
```json
{
  "type": "password",
  "value": "TestPassword123!"
}
```

Success Response (200 OK):
```json
{
  "result": {
    "valid": true,
    "errors": []
  }
}
```

## Tenant Endpoints

### Tenant Management

#### GET /api/tenants
List all tenants with pagination and filters. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Query Parameters:
- `page`: Page number for pagination (default: 1)
- `limit`: Number of tenants per page (default: 10)
- `search`: Search term for filtering tenants by name or slug
- `filter`: Additional filters as key-value pairs

Success Response (200 OK):
```json
{
  "tenants": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "name": "Tenant Name",
      "slug": "tenant-name",
      "type": "team",
      "createdAt": "2023-01-01T00:00:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "limit": 10
}
```

#### POST /api/tenants
Create a new tenant. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Request:
```json
{
  "name": "New Tenant",
  "type": "team"
}
```

Success Response (201 Created):
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "New Tenant",
  "slug": "new-tenant",
  "type": "team",
  "createdAt": "2023-01-01T00:00:00Z"
}
```

#### GET /api/tenants/:id
Get details of a specific tenant. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Path Parameters:
- `id`: ID of the tenant

Success Response (200 OK):
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "Tenant Name",
  "slug": "tenant-name",
  "type": "team",
  "createdAt": "2023-01-01T00:00:00Z"
}
```

#### PUT /api/tenants/:id
Update a tenant. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Path Parameters:
- `id`: ID of the tenant

Request:
```json
{
  "name": "Updated Tenant Name"
}
```

Success Response (200 OK):
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "Updated Tenant Name",
  "slug": "updated-tenant-name",
  "type": "team",
  "createdAt": "2023-01-01T00:00:00Z"
}
```

#### DELETE /api/tenants/:id
Delete a tenant. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Path Parameters:
- `id`: ID of the tenant

Success Response (200 OK):
```json
{
  "message": "Tenant deleted successfully"
}
```

### User Endpoints

#### User Management

##### GET /api/users
List all users with pagination and filters. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Query Parameters:
- `page`: Page number for pagination (default: 1)
- `limit`: Number of users per page (default: 10)
- `search`: Search term for filtering users by name or email
- `filter`: Additional filters as key-value pairs

Success Response (200 OK):
```json
{
  "users": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "name": "John Doe",
      "email": "john.doe@example.com",
      "status": "active",
      "createdAt": "2023-01-01T00:00:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "limit": 10
}
```

##### POST /api/users
Create a new user. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Request:
```json
{
  "name": "Jane Doe",
  "email": "jane.doe@example.com",
  "password": "securepassword123"
}
```

Success Response (201 Created):
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174001",
  "name": "Jane Doe",
  "email": "jane.doe@example.com",
  "status": "active",
  "createdAt": "2023-01-01T00:00:00Z"
}
```

##### GET /api/users/:id
Get details of a specific user. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Path Parameters:
- `id`: ID of the user

Success Response (200 OK):
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "John Doe",
  "email": "john.doe@example.com",
  "status": "active",
  "createdAt": "2023-01-01T00:00:00Z"
}
```

##### PUT /api/users/:id
Update a user. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Path Parameters:
- `id`: ID of the user

Request:
```json
{
  "name": "Johnathan Doe"
}
```

Success Response (200 OK):
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "Johnathan Doe",
  "email": "john.doe@example.com",
  "status": "active",
  "createdAt": "2023-01-01T00:00:00Z"
}
```

##### DELETE /api/users/:id
Delete a user. Requires authentication.

Headers:
```
Authorization: Bearer <access_token>
```

Path Parameters:
- `id`: ID of the user

Success Response (200 OK):
```json
{
  "message": "User deleted successfully"
}
``` 