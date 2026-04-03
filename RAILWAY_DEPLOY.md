# Railway Deployment Guide

## Services Setup (deploy in order)

### 1. Database (PostgreSQL)
- Add **Database** from Railway catalog
- Copy the `DATABASE_URL` connection string

### 2. Backend Service
- **Source**: Connect this repository
- **Root Directory**: `/`
- **Dockerfile**: `Dockerfile.backend`
- **Environment Variables**:
  ```
  DB_DRIVER=postgres
  DATABASE_URL=<from step 1>
  APP_NAME=NEXUS
  APP_PORT=3000
  APP_SECRET=<generate-a-random-string>
  JWT_SECRET=<generate-a-random-string>
  JWT_EXPIRE_HOURS=72
  ```
- Railway will expose this on `https://<backend>.railway.app`

### 3. Frontend Service
- Add new service → **Docker**
- **Dockerfile**: `Dockerfile.frontend`
- **Build Args**: `VITE_API_URL=https://<backend>.railway.app` (use your actual backend URL)
- Will be served on port 80

## One-time Setup
1. First register an account on the deployed backend
2. The first user becomes **admin** automatically
3. Configure Wings nodes via the admin dashboard

## Local Development
```bash
# Backend (SQLite, local)
cd nexus && CGO_ENABLED=1 go build -tags sqlite -o nexus-local ./backend/main.go && ./nexus-local

# Frontend
cd nexus/frontend && npm install && npm run dev
```