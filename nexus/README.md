# NEXUS - Unified Game Server Management Panel

Phase 2: Admin Panel + Real-time Frontend (Complete)

## Overview

NEXUS is a lightweight game server management panel designed for 512MB RAM VPS deployments. It provides centralized control over multiple Pterodactyl Wings nodes with user management, server provisioning, and console access.

**Tech Stack:**
- **Backend:** Go + Fiber + GORM (Port 3000)
- **Database:** SQLite (default), MySQL optional
- **Auth:** JWT tokens
- **Frontend:** React + TypeScript + Vite (Port 5173)
- **Console:** xterm.js with WebSocket real-time
- **State:** Zustand

## Project Structure

```
nexus/
├── backend/                  # Go backend (Port 3000)
│   ├── main.go
│   ├── config/config.go
│   ├── database/database.go
│   ├── models/
│   │   ├── user.go
│   │   ├── node.go
│   │   ├── server.go
│   │   ├── egg.go
│   │   ├── allocation.go
│   │   ├── ticket.go
│   │   └── coin_transaction.go
│   ├── routes/routes.go
│   ├── controllers/
│   │   ├── auth_controller.go
│   │   ├── user_controller.go
│   │   ├── server_controller.go
│   │   ├── node_controller.go
│   │   └── egg_controller.go
│   ├── middleware/
│   │   ├── auth.go
│   │   └── admin.go
│   ├── wings/
│   │   ├── client.go
│   │   ├── websocket.go
│   │   └── types.go
│   └── utils/
│       ├── jwt.go
│       └── response.go
├── frontend/                 # React TypeScript frontend (Port 5173)
│   ├── package.json
│   ├── vite.config.ts
│   ├── tsconfig.json
│   ├── src/
│   │   ├── api/
│   │   ├── store/
│   │   ├── hooks/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── router/
│   │   ├── types/
│   │   ├── App.tsx
│   │   └── main.tsx
│   ├── README.md            # Detailed frontend docs
│   └── .env.example
├── .env.example
├── .gitignore
└── README.md
```

## Installation

### Prerequisites

- Go 1.21+
- Node.js 18+ (for frontend later)
- Database (SQLite default, or MySQL 8+)

### Setup Steps

1. **Clone and enter directory:**
   ```bash
   cd nexus/backend
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Configure environment:**
   ```bash
   cp .env.example .env
   # Edit .env with your settings
   ```

   Key configuration:
   - `DB_DRIVER`: `sqlite` (default) or `mysql`
   - For SQLite: `DB_PATH=./nexus.db`
   - For MySQL: set `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASS`
   - `JWT_SECRET`: Generate a strong random secret
   - `WINGS_TOKEN_ID` and `WINGS_TOKEN`: Credentials for each Wings node

4. **Run the server:**
   ```bash
   go run ./backend/main.go
   ```

   Server starts on `:3000` (or `APP_PORT` configured)

5. **Initialize database:**
   - The database auto-migrates on first run
   - First user can self-register as admin

6. **Install frontend dependencies:**
   ```bash
   cd ../frontend
   npm install
   ```

7. **Configure frontend:**
   ```bash
   cp .env.example .env
   # Ensure VITE_API_URL=http://localhost:3000
   ```

8. **Run frontend dev server:**
   ```bash
   npm run dev
   ```
   Frontend runs on http://localhost:5173

9. **Access the panel:**
   - Open http://localhost:5173
   - Register first user (self-assigns admin role)
   - Login and access admin panel

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login
- `GET  /api/auth/me` - Get current user (protected)

### Admin Routes (require admin role)
- `GET    /api/admin/users` - List all users
- `POST   /api/admin/users` - Create user
- `GET    /api/admin/users/:id` - Get user
- `PATCH  /api/admin/users/:id` - Update user
- `DELETE /api/admin/users/:id` - Delete user

- `GET    /api/admin/nodes` - List all nodes
- `POST   /api/admin/nodes` - Create node
- `GET    /api/admin/nodes/:id` - Get node
- `GET    /api/admin/nodes/:id/stats` - Get node statistics
- `PATCH  /api/admin/nodes/:id` - Update node
- `DELETE /api/admin/nodes/:id` - Delete node

- `GET    /api/admin/servers` - List all servers
- `POST   /api/admin/servers` - Create server
- `GET    /api/admin/servers/:id` - Get server
- `PATCH  /api/admin/servers/:id` - Update server
- `DELETE /api/admin/servers/:id` - Delete server
- `POST   /api/admin/servers/:id/power` - Send power action

- `GET    /api/admin/eggs` - List all eggs
- `POST   /api/admin/eggs` - Create egg
- `GET    /api/admin/eggs/:id` - Get egg
- `PATCH  /api/admin/eggs/:id` - Update egg
- `DELETE /api/admin/eggs/:id` - Delete egg

### Client Routes (authenticated)
- `GET  /api/client/servers` - List user's servers
- `GET  /api/client/servers/:uuid` - Get server details
- `GET  /api/client/servers/:uuid/resources` - Get live resource usage

### Health Check
- `GET /health` - Health status

## Data Models

### User
- `id`, `uuid`, `username`, `email`, `password` (bcrypt)
- `role`: "admin" or "client"
- `coins` (int), `root_admin` (bool)
- `created_at`, `updated_at`

### Node
- `id`, `uuid`, `name`, `fqdn`
- `scheme`: "http" or "https"
- `wings_port`: default 8080
- `memory`, `disk` (in MB)
- `memory_overalloc`, `disk_overalloc` (percentage)
- `token_id`, `token` (Wings API credentials)
- `created_at`, `updated_at`

### Server
- `id`, `uuid`, `name`
- `user_id` (FK), `node_id` (FK), `egg_id` (FK), `allocation_id` (FK)
- `memory`, `disk`, `cpu` (resource limits)
- `status`: "installing", "running", "stopped", "error"
- `suspended` (bool)
- `created_at`, `updated_at`

### Egg
- `id`, `uuid`, `name`, `description`
- `docker_image`, `startup_command`
- `created_at`, `updated_at`

### Allocation
- `id`, `node_id` (FK), `ip`, `port`
- `assigned` (bool), `server_id` (nullable FK)
- `created_at`, `updated_at`

### Ticket
- `id`, `user_id` (FK), `subject`
- `status`: "open", "closed", "pending"
- `priority`: "low", "medium", "high"
- `created_at`, `updated_at`

### CoinTransaction
- `id`, `user_id` (FK), `amount` (can be negative)
- `reason`
- `created_at`

## Wings API Client

The `wings` package provides:
- **HTTP Client:** `GetServerDetails`, `CreateServer`, `DeleteServer`, `SendPowerAction`, `GetServerResources`
- **WebSocket Proxy:** Bidirectional console connection between browser and Wings daemon

### Wings API Authentication

Uses `Authorization: Bearer {token}` where token is the node's `token` (decrypted). Base URL is constructed as `{scheme}://{fqdn}:{wings_port}`.

## Performance Optimizations (512MB VPS Target)

- **GORM:** Logger disabled in production
- **Database connection pool:** max 10 open, max 5 idle connections
- **Fiber timeouts:** Read, Write, Idle timeouts configured
- **No unnecessary dependencies** - minimal go.mod
- **Goroutines:** Per-connection WebSocket handling (cleaned up on disconnect)
- **Indexed queries:** All foreign keys have database indexes

## Development

### Adding New Routes

1. Add handler method in appropriate `controllers/*.go`
2. Register in `routes/routes.go`

### Database Migrations

Auto-migration runs on startup. For production changes:
- Add new columns with `gorm:"default:..."` or nullable
- Do NOT remove old columns immediately

### Testing API

```bash
# Register
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","email":"admin@example.com","password":"password","role":"admin"}'

# Login
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password"}'

# Use token
TOKEN="<jwt_token>"
curl -H "Authorization: Bearer $TOKEN" http://localhost:3000/api/auth/me
```

## Phase 2 Status (Admin Panel Complete)

- [x] All npm dependencies defined
- [x] API layer complete with interceptors (axios + JWT)
- [x] Zustand auth store with localStorage persistence
- [x] React Router v6 with protected admin routes
- [x] Login & Register pages via Stitch MCP design
- [x] Admin Dashboard with 4 stat cards via Stitch
- [x] Admin Servers list with DataTable via Stitch
- [x] Admin Server Detail with live console + stats via Stitch
- [x] Admin Nodes page via Stitch
- [x] Admin Users page via Stitch
- [x] Admin Eggs page via Stitch
- [x] useConsole hook with xterm.js WebSocket
- [x] useServerStats hook with 3s polling
- [x] Vite proxy configured for /api and /ws
- [x] Frontend starts with `npm run dev`
- [x] Frontend builds with `npm run build`

## Next Phases

- **Phase 3:** Client dashboard (user-facing server list + console)
- **Phase 4:** Advanced features (backups, file manager, tickets, coins)

## License

Proprietary - All Rights Reserved
