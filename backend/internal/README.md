# EzEdu Backend Internal Packages

Internal Go application packages implementing core application domain logic and APIs.

---

## 📂 Packages

- **`auth/`**: Authentication services, password hashing (bcrypt), JWT/session cookie handling, and middleware.
- **`handler/`**: HTTP handlers and routing for REST API endpoints (Auth, Categories, Lessons, Progress, etc.).
- **`model/`**: Go data structs and JSON response contracts.
- **`store/`**: SQLite database store, queries, and schema migration logic.
