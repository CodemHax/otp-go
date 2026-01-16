# OTP-Go

A simple Go backend for sending and verifying OTP codes via email.

## Features
- Send OTP codes to email addresses using SMTP
- Verify OTP codes securely
- Store OTPs in Redis with configurable connection
- Environment-based configuration for SMTP and Redis

## Setup

1. Copy `.env.example` to `.env` and fill in your SMTP and Redis credentials:
   ```
   cp .env.example .env
   ```
   Example `.env` variables:
   ```
   SMTP_HOST=mail.smtp2go.com
   SMTP_PORT=2525
   SMTP_USER=your_smtp_username
   SMTP_PASS=your_smtp_password
   SMTP_FROM=your_email@example.com
   ENC_KEY=your_32_byte_encryption_key
   REDIS_HOST=127.0.0.1
   REDIS_PORT=6379
   REDIS_PASSWORD=
   REDIS_DB=0
   ```
2. Install Go dependencies:
   ```
   go mod tidy
   ```
3. Run the server:
   ```
   go build -o main.exe main.go
   ./main.exe
   ```

## Endpoints

- `POST /otp` — Send OTP to email
  - JSON body: `{ "email": "user@example.com" }`
  - Response: `{ "message": "OTP sent successfully" }`
- `POST /verify` — Verify OTP
  - JSON body: `{ "email": "user@example.com", "code": "123456" }`
  - Response: `{ "message": "otp verified successfully" }` or error
- `GET /health` — Check Redis connection
  - Response: `{ "redis": "connected" }`

## Redis Configuration

Redis connection is configured via environment variables. You can use a local or remote Redis instance. Make sure the credentials in `.env` match your Redis setup.

## Python Verification Example

Install requests:
```
pip install requests
```

Use `verify_otp.py`:
```python
from verify_otp import verify_otp
verify_otp("user@example.com", "123456")
```

## Security Notes
- OTPs are hashed before storage and can be encrypted if you set `ENC_KEY`.
- Always use strong SMTP and Redis passwords in production.

## License
MIT
