# Fly Crypto API Documentation

## Base URL
`http://localhost:8080`

## Authentication
The API uses Bearer token authentication. After login and verification, include the token in all protected endpoints:
```
Authorization: Bearer <your_token>
```

## Endpoints

### Health Check
#### GET /ping
Check if the API is running.

**Response**
```json
{
    "message": "pong"
}
```

### Authentication

#### POST /login
Initiate login process by requesting a verification code.

**Request Body**
```json
{
    "email": "user@example.com"
}
```

**Response**
```json
{
    "message": "Verification code sent"
}
```

#### POST /verify
Verify email with the received code to get an authentication token.

**Request Body**
```json
{
    "email": "user@example.com",
    "code": "123456"
}
```

**Response**
```json
{
    "token": "your_auth_token"
}
```

### Notifications (Protected Endpoints)

#### POST /notifications
Create a new price notification.

**Headers**
```
Authorization: Bearer <your_token>
```

**Request Body**
```json
{
    "coin_symbol": "BTC",
    "target_price": 50000.00
}
```

**Notes**
- `coin_symbol`: Supported values are "BTC", "ETH", "SOL"
- `target_price`: Target price for notification

**Response**
```json
{
    "message": "Notification created successfully",
    "data": {
        "id": 1,
        "user_id": 123,
        "coin_symbol": "BTC",
        "target_price": 50000.00,
        "is_above": true
    }
}
```

#### GET /notifications
Get all notifications for the authenticated user.

**Headers**
```
Authorization: Bearer <your_token>
```

**Response**
```json
{
    "notifications": [
        {
            "id": 1,
            "user_id": 123,
            "coin_symbol": "BTC",
            "target_price": 50000.00,
            "is_above": true
        }
    ]
}
```

#### DELETE /notifications/:id
Delete a specific notification.

**Headers**
```
Authorization: Bearer <your_token>
```

**Parameters**
- `id`: Notification ID to delete

**Response**
```json
{
    "message": "Notification deleted successfully"
}
```

## Error Responses
All endpoints may return error responses in the following format:

```json
{
    "error": "Error message description"
}
```

Common HTTP status codes:
- 200: Success
- 400: Bad Request
- 401: Unauthorized
- 500: Internal Server Error