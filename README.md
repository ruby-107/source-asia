Rate Limited API Service (Golang)

Overview

This project is a simple API service that:

Accepts user requests
Applies rate limiting (5 requests per user per minute)
Uses a queue + worker for async processing
Stores data in PostgreSQL
Handles concurrent requests safely

How to Run
1. Clone the repo
   git clone https://github.com/ruby-107/source-asia.git
   cd source-asia

2. Create .env file
 PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=your_db


3. Create table
   CREATE TABLE production.users (
   id SERIAL PRIMARY KEY,
   user_id INT,
   payload VARCHAR(100),
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );

4. Run project
   go mod tidy
   go run cmd/main.go

APIs
🔹 POST /request

Request

{
"user_id": 1,
"payload": "test"
}

Success Response

{
"success": true,
"data": {
"user_id": 1,
"payload": "test",
"status": "accepted"
}
}

Rate Limit Error

{
"success": false,
"error": "rate limit exceeded"
}

GET /stats

Response

[
{
"user_id": 1,
"request_count": 5
}
]

Design (Simple)
Rate Limiter → map + mutex (thread-safe)
Queue → buffered channel (handles burst traffic)
Worker → processes requests in background
Retry Logic → retries DB insert 3 times if failed

Flow
Client → API → Rate Limit → Queue → Worker → DB


Improvements
Use Redis for distributed rate limiting
Add multiple workers (worker pool)
Add request tracking


# source-asia
