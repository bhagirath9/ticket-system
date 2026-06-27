# Ticket Management System Backend

A production-ready, clean, and simple REST API built using **Golang**, the **Gin Web Framework**, **GORM (ORM)**, and **SQLite** database. 

This is a backend-only application featuring complete user registration, JWT-based authentication, and a robust status workflow for ticket management with strict authorization rules.

---

## Folder Structure

```text
TicketSystem/
├── cmd/
│   └── main.go                  # Application entry point
├── config/
│   └── config.go                # Configuration loader using godotenv
├── database/
│   └── database.go              # Database GORM SQLite connection & migrations
├── models/
│   ├── user.go                  # User GORM model, requests, and response DTOs
│   └── ticket.go                # Ticket GORM model, requests, and status definitions
├── repository/
│   ├── user_repository.go       # Data access layer for users
│   └── ticket_repository.go     # Data access layer for tickets
├── services/
│   ├── auth_service.go          # Auth business logic (registration, validation, tokens)
│   └── ticket_service.go        # Ticket business logic (transitions, owner restrictions)
├── controllers/
│   ├── auth_controller.go       # Route handlers for authentication
│   └── ticket_controller.go     # Route handlers for ticket management
├── middleware/
│   └── jwt_middleware.go        # Authentication filter verifying Bearer tokens
├── routes/
│   └── routes.go                # Path definitions and route registration
├── utils/
│   └── jwt.go                   # JWT creation and parsing helpers
├── .env.example                 # Environment variables template
├── .env                         # Local environment configuration
├── Dockerfile                   # Docker build specification
└── README.md                    # Project documentation
```

---

## Environment Variables

The application reads variables from a `.env` file in the root folder.

| Variable Name  | Description | Default Value |
| -------------  | ----------- | ------------- |
| `PORT`         | Port on which the HTTP server runs | `8080` |
| `JWT_SECRET`   | Secret key used to sign authentication tokens | `your_secret_key` |
| `DATABASE_URL` | File path for SQLite database file | `tickets.db` |

---

## Local Setup

### Prerequisites

- **Go 1.21+** installed locally (if running bare-metal)
- **Docker** installed locally (if containerized)

### Steps to Run Locally (without Docker)

1. Clone or copy the project files to your workspace directory.
2. Initialize environment configurations:
   ```bash
   cp .env.example .env
   ```
3. Run the Go tidy command to resolve dependencies:
   ```bash
   go mod tidy
   ```
4. Start the application:
   ```bash
   go run cmd/main.go
   ```
   The backend will start listening on [http://localhost:8080](http://localhost:8080).

---

## Docker Commands

This project is fully dockerized with a multi-stage `Dockerfile` to keep production images tiny and optimized.

### 1. Build Docker Image
```bash
docker build -t ticket-system .
```

### 2. Run Docker Container
```bash
docker run -p 8080:8080 ticket-system
```
Once run, the application is live at [http://localhost:8080](http://localhost:8080).

---

## API Endpoints List

All request and response bodies use JSON formatting. Protected endpoints require:
`Authorization: Bearer <token>` in headers.

### Public Endpoints

#### 1. Health Check
* **Method**: `GET`
* **URL**: `/health`
* **Response**: `200 OK`
  ```json
  {
    "status": "ok"
  }
  ```

#### 2. User Registration
* **Method**: `POST`
* **URL**: `/auth/register`
* **Body**:
  ```json
  {
    "name": "Bhagirath",
    "email": "abc@gmail.com",
    "password": "mysecretpassword123"
  }
  ```
* **Response**: `201 Created`
  ```json
  {
    "message": "User Registered Successfully"
  }
  ```
* **Error**: `409 Conflict` (If email already registered), `400 Bad Request` (Missing fields or invalid password)

#### 3. User Login
* **Method**: `POST`
* **URL**: `/auth/login`
* **Body**:
  ```json
  {
    "email": "abc@gmail.com",
    "password": "mysecretpassword123"
  }
  ```
* **Response**: `200 OK`
  ```json
  {
    "token": "JWT_TOKEN_HERE"
  }
  ```
* **Error**: `401 Unauthorized` (Invalid credentials)

---

### Protected Endpoints

#### 4. Create Ticket
* **Method**: `POST`
* **URL**: `/tickets`
* **Headers**: `Authorization: Bearer <token>`
* **Body**:
  ```json
  {
    "title": "Payment Issue",
    "description": "Unable to pay fees"
  }
  ```
* **Response**: `201 Created`
  ```json
  {
    "id": 1,
    "title": "Payment Issue",
    "description": "Unable to pay fees",
    "status": "open",
    "user_id": 1,
    "created_at": "2026-06-27T11:00:00Z",
    "updated_at": "2026-06-27T11:00:00Z"
  }
  ```

#### 5. List Own Tickets
* **Method**: `GET`
* **URL**: `/tickets`
* **Headers**: `Authorization: Bearer <token>`
* **Response**: `200 OK`
  ```json
  [
    {
      "id": 1,
      "title": "Payment Issue",
      "description": "Unable to pay fees",
      "status": "open",
      "user_id": 1,
      "created_at": "2026-06-27T11:00:00Z",
      "updated_at": "2026-06-27T11:00:00Z"
    }
  ]
  ```

#### 6. Get Own Ticket by ID
* **Method**: `GET`
* **URL**: `/tickets/{id}`
* **Headers**: `Authorization: Bearer <token>`
* **Response**: `200 OK` (if owned by user), `403 Forbidden` (if ticket belongs to someone else), `404 Not Found` (if ticket doesn't exist)
  ```json
  {
    "id": 1,
    "title": "Payment Issue",
    "description": "Unable to pay fees",
    "status": "open",
    "user_id": 1,
    "created_at": "2026-06-27T11:00:00Z",
    "updated_at": "2026-06-27T11:00:00Z"
  }
  ```

#### 7. Update Own Ticket Status
* **Method**: `PATCH`
* **URL**: `/tickets/{id}/status`
* **Headers**: `Authorization: Bearer <token>`
* **Body**:
  ```json
  {
    "status": "in_progress"
  }
  ```
* **Response**: `200 OK` with updated Ticket JSON.
* **Status Rules**:
  - Valid statuses are `open`, `in_progress`, and `closed`.
  - Linear flow: `open` ➔ `in_progress` ➔ `closed`.
  - Reverting status backwards (e.g. `in_progress` to `open`) returns `400 Bad Request`.
  - Once status reaches `closed`, the ticket cannot be modified. Trying to change the status of a closed ticket returns `400 Bad Request`:
    ```json
    {
      "message": "Closed ticket cannot be reopened"
    }
    ```

---

## API Testing Flow

Here is a quick curl sequence to test all requirements sequentially.

### 1. Register User A
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"User A","email":"user_a@gmail.com","password":"password123"}'
```

### 2. Login User A (Retrieve Token A)
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user_a@gmail.com","password":"password123"}'
```

### 3. Create Ticket under User A
```bash
curl -X POST http://localhost:8080/tickets \
  -H "Authorization: Bearer <TOKEN_A>" \
  -H "Content-Type: application/json" \
  -d '{"title":"Ticket 1","description":"This is ticket 1 description"}'
```
Assume this returns ticket with ID `1`.

### 4. Register User B
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"User B","email":"user_b@gmail.com","password":"password123"}'
```

### 5. Login User B (Retrieve Token B)
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user_b@gmail.com","password":"password123"}'
```

### 6. User B Tries to Read User A's Ticket (ID 1)
```bash
curl -X GET http://localhost:8080/tickets/1 \
  -H "Authorization: Bearer <TOKEN_B>"
```
*Expected Response*: `403 Forbidden`

### 7. User A Updates Status to in_progress
```bash
curl -X PATCH http://localhost:8080/tickets/1/status \
  -H "Authorization: Bearer <TOKEN_A>" \
  -H "Content-Type: application/json" \
  -d '{"status":"in_progress"}'
```
*Expected Response*: `200 OK` (Status updated to `in_progress`)

### 8. User A Updates Status to closed
```bash
curl -X PATCH http://localhost:8080/tickets/1/status \
  -H "Authorization: Bearer <TOKEN_A>" \
  -H "Content-Type: application/json" \
  -d '{"status":"closed"}'
```
*Expected Response*: `200 OK` (Status updated to `closed`)

### 9. User A Tries to Reopen Ticket
```bash
curl -X PATCH http://localhost:8080/tickets/1/status \
  -H "Authorization: Bearer <TOKEN_A>" \
  -H "Content-Type: application/json" \
  -d '{"status":"in_progress"}'
```
*Expected Response*: `400 Bad Request` with:
```json
{
  "message": "Closed ticket cannot be reopened"
}
```

---

## Postman Collection Setup

1. Open Postman, click on **Import** and create a new request collection.
2. Define collection-wide variables:
   - `baseUrl` = `http://localhost:8080` (or your deployed URL)
   - `jwtToken` = (Leave blank; copy and paste from the login response)
3. For protected endpoints, select **Authorization** -> **Bearer Token** type and set the value to `{{jwtToken}}`.

---

## Deployment Instructions

You can deploy this Dockerized service for free on platforms like **Koyeb**, **Render**, or **Railway**:

### Deploying on Koyeb (Recommended)
1. Link your GitHub repository to your Koyeb Account.
2. Select **Go** buildpack or choose **Docker** as the deployment type (Docker is recommended since a Dockerfile is provided).
3. Set the Environment Variables:
   - `PORT=8080`
   - `JWT_SECRET=production_random_strong_key_xyz`
   - `DATABASE_URL=/secrets/tickets.db` (or just leave blank to default to local persistent disk `/app/tickets.db`)
4. Expose port `8080` to public HTTP traffic.
5. Deploy! Koyeb builds the container automatically and provisions a public URL (e.g. `https://your-service-name.koyeb.app`).
6. Test your live health endpoint: `https://your-service-name.koyeb.app/health`.

### Deploying on Render
1. Create a new **Web Service** on Render and connect your repository.
2. Choose **Docker** as the runtime (Render will automatically detect the `Dockerfile`).
3. Under Environment variables, add:
   - `JWT_SECRET`
   - `PORT` (Render defaults to routing to the exposed Docker port)
4. Select the Free tier plan and click **Create Web Service**.
5. Once building completes, your service is publicly online!

---

## Deployed URL
- **Production URL**: *(To be populated upon candidate deployment)*
- **Health Check Endpoint**: `/health`
