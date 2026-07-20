# EzEdu — Interactive Kids' Learning App

Platform belajar interaktif untuk anak usia 4–12 tahun.

## Tech Stack

- **Frontend**: Astro 5 + Preact (islands architecture)
- **Backend**: Go + Chi router
- **Database**: SQLite (WAL mode)
- **Reverse Proxy**: Caddy 2

## Quick Start

### Prerequisites
- Go 1.22+
- Node.js 20+

### Backend
```bash
cd backend
go mod download
go run ./cmd/server    # Starts on :8080
```

### Frontend
```bash
cd frontend
npm install
npm run dev            # Dev server on :4321 (proxies /api to :8080)
```

### Production Build
```bash
cd frontend && npm run build    # Static output to dist/
cd backend && go build -o server ./cmd/server
```

## Project Structure

```
ezedu/
├── backend/
│   ├── cmd/server/         # Entry point
│   ├── internal/
│   │   ├── auth/           # Auth service & middleware
│   │   ├── handler/        # HTTP handlers
│   │   ├── model/          # Data models
│   │   └── store/          # SQLite persistence
│   └── data/               # SQLite DB (gitignored)
├── frontend/
│   ├── src/
│   │   ├── components/     # Preact interactive components
│   │   ├── layouts/        # Astro layouts
│   │   ├── lib/            # API client, utilities
│   │   ├── pages/          # Astro pages (SSG)
│   │   └── styles/         # CSS design system
│   └── public/             # Static assets
└── deploy/
    ├── Caddyfile           # Production config
    └── Caddyfile.dev       # Development config
```

## UI Language

All user-facing text is in **Bahasa Indonesia** (Indonesian).
