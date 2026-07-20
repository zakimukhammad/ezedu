# EzEdu Backend Service

The backend server for EzEdu, built with **Go** and the **Chi router**, utilizing **SQLite** for lightweight, self-contained data storage.

---

## 📁 Directory Structure

```
backend/
├── cmd/
│   └── server/          # Entry point (main package)
├── internal/
│   ├── auth/            # Authentication logic, password hashing, session management
│   ├── handler/         # HTTP handlers & REST API routes
│   ├── model/           # Data models & domain types
│   └── store/           # SQLite database persistence layer & migrations
└── data/                # SQLite database storage (gitignored)
```

---

## 🚀 Getting Started

### Prerequisites
- **Go**: `1.22+`

### Installation & Run

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```

2. Download dependencies:
   ```bash
   go mod download
   ```

3. Start the development server (runs on `http://localhost:8080`):
   ```bash
   go run ./cmd/server
   ```

### Building for Production

```bash
go build -o server ./cmd/server
```

On Windows:
```cmd
go build -o server.exe ./cmd/server
```

---

## ⚙️ Environment Variables

Optionally set env variables or use default fallbacks:
- `PORT`: Server port (default: `8080`)
- `DB_PATH`: SQLite database file path (default: `./data/ezedu.db`)
- `SESSION_SECRET`: Secret key for session security
