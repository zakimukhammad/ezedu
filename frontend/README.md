# EzEdu Frontend Application

The web frontend for **EzEdu**, built using **Astro 5** with **Preact islands** for lightweight interactive client-side components.

All user interface copy, navigation, buttons, and content are written in **Bahasa Indonesia**.

---

## 📁 Directory Structure

```
frontend/
├── public/              # Static assets (images, audio, icons)
├── src/
│   ├── components/      # Preact interactive components (LoginForm, Progress, etc.)
│   ├── layouts/         # Astro layout wrappers (BaseLayout, AuthLayout)
│   ├── lib/             # API client, utilities, and helper functions
│   ├── pages/           # Astro page routes (SSG / SSR pages)
│   └── styles/          # Global CSS design system and tokens
├── astro.config.mjs     # Astro framework configuration
├── package.json         # Node.js dependencies & scripts
└── tsconfig.json        # TypeScript configuration
```

---

## 🚀 Getting Started

### Prerequisites
- **Node.js**: `20+`
- **npm**: `10+`

### Running Locally

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Start the development server (runs on `http://localhost:4321` with proxy to backend on `:8080`):
   ```bash
   npm run dev
   ```

> **Note on Windows PowerShell execution policy**:  
> If `npm` commands fail with execution policy errors, run `npm.cmd run dev` or run `Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser` in PowerShell.

---

## 📜 Available Scripts

| Script | Command | Purpose |
| :--- | :--- | :--- |
| **Dev Server** | `npm run dev` | Starts local dev server at `localhost:4321` |
| **Build** | `npm run build` | Compiles static production output to `./dist/` |
| **Preview** | `npm run preview` | Previews the build output locally |
| **Astro CLI** | `npx astro ...` | Runs Astro CLI commands |
