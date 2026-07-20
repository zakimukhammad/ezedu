# EzEdu Deployment & Reverse Proxy Configuration

Deployment configurations for running **EzEdu** behind **Caddy 2** reverse proxy in both local development and production environments.

---

## 📁 Files

- **`Caddyfile`**: Production Caddy server configuration with TLS/HTTPS, static asset caching, API proxying to `:8080`, rate limiting, and security headers.
- **`Caddyfile.dev`**: Development configuration for local testing with automatic proxying to local services.

---

## 🚀 Usage

### Development Mode

Run Caddy locally with dev config:
```bash
caddy run --config deploy/Caddyfile.dev
```

### Production Deployment

1. Build static frontend assets and Go binary:
   ```bash
   cd frontend && npm run build
   cd ../backend && go build -o server ./cmd/server
   ```

2. Copy compiled static files (`frontend/dist/`) to `/opt/ezedu/public` (or your chosen web root).

3. Set your target domain:
   ```bash
   export EZEDU_DOMAIN=ezedu.yourdomain.com
   ```

4. Run Caddy with production configuration:
   ```bash
   caddy run --config deploy/Caddyfile
   ```
