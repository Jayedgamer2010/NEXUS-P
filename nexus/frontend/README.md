# NEXUS Frontend - Phase 2: Admin Panel

Complete admin dashboard built with React, TypeScript, Vite, and xterm.js. Features real-time server console with WebSocket proxy and live resource monitoring.

## Tech Stack

- **React 18** - UI library
- **TypeScript** - Type safety
- **Vite** - Build tool and dev server
- **React Router v6** - Routing
- **Axios** - HTTP client
- **Zustand** - State management
- **xterm.js** - Terminal emulator for console
- **Tailwind not used** - All custom CSS per design spec

## Design Spec

Dark theme with:
- Background: #0a0a0f
- Card surfaces: #1a1a2e
- Accent: Purple gradient #7c3aed → #3b82f6
- Font: Inter (UI), JetBrains Mono (console)
- Sharp borders, 8px max radius

## Project Structure

```
frontend/src/
├── api/
│   ├── client.ts          - Axios instance with JWT interceptor
│   ├── auth.ts            - Auth endpoints
│   ├── admin.ts           - Admin CRUD endpoints
│   └── wings.ts           - Server resources polling
├── store/
│   └── authStore.ts       - Zustand auth state with localStorage
├── hooks/
│   ├── useConsole.ts      - WebSocket console connection
│   └── useServerStats.ts  - Poll server resources every 3s
├── components/
│   ├── layout/
│   │   ├── Sidebar.tsx    - Navigation sidebar
│   │   ├── Header.tsx     - Top header bar
│   │   └── Layout.tsx     - Main layout wrapper
│   ├── ui/
│   │   ├── StatusBadge.tsx - Status indicator pill
│   │   ├── StatCard.tsx    - Live stats card
│   │   ├── DataTable.tsx   - Reusable table with pagination
│   │   └── Modal.tsx       - Reusable modal
│   └── console/
│       ├── TerminalDisplay.tsx - xterm.js terminal
│       └── ConsoleInput.tsx    - Command input bar
├── pages/
│   ├── auth/
│   │   ├── Login.tsx
│   │   └── Register.tsx
│   ├── admin/
│   │   ├── Dashboard.tsx
│   │   ├── Servers.tsx
│   │   ├── ServerDetail.tsx  - Live console + stats
│   │   ├── Nodes.tsx
│   │   ├── Users.tsx
│   │   └── Eggs.tsx
│   └── errors/
│       ├── 404.tsx
│       └── 403.tsx
├── router/
│   └── index.tsx          - React Router v6 protected routes
├── types/
│   └── index.ts           - TypeScript interfaces
├── App.tsx
└── main.tsx
```

## Setup

### Prerequisites

- Node.js 18+
- Backend running on http://localhost:3000

### Installation

```bash
cd frontend
npm install
```

### Configuration

Create `.env` file:

```env
VITE_API_URL=http://localhost:3000
```

### Development

```bash
npm run dev
```

Open http://localhost:5173

The dev server proxies `/api` and `/ws` to the backend.

### Build

```bash
npm run build
```

Output in `dist/` folder.

## Features Implemented

### Authentication
- JWT-based auth with localStorage persistence
- Protected routes with admin role guard
- Login and registration pages
- Automatic 401 redirect to login

### Admin Dashboard
- 4 stat cards: Total Servers, Nodes, Users, Running Servers
- Recent servers table with status badges
- Real-time data fetching

### Server Management
- Full server list with DataTable
- Create/Edit/Delete server modals
- Power actions (Start/Stop/Restart/Kill)
- Inline status indicators

### Server Detail (Most Important)
- Server header with name, UUID, status badge
- 4 power action buttons
- 3 live stat cards (CPU %, RAM, Disk) - updates every 3s
- Real-time console with xterm.js
- WebSocket console input for sending commands
- Auto-reconnect on disconnect

### Node, User, Egg Management
- Complete CRUD for all admin resources
- Modals for create/edit operations
- DataTables with pagination
- Confirmation for deletions

### Errors
- Custom 403 Forbidden page
- Custom 404 Not Found page

## API Integration

All API calls use the centralized `api/client.ts` with:
- Automatic JWT header attachment
- 401 interception → localStorage clear → redirect
- Standardized response handling

WebSocket console uses raw WebSocket API with JWT in query string.

## Performance Optimizations

- React.memo on components (can be added where needed)
- useServerStats: clearInterval on unmount
- useConsole: clean up WebSocket + timeouts on unmount
- Lazy loading ready (pages can be wrapped in React.lazy)
- xterm: fitAddon resize optimized

## TypeScript Types

Full interface coverage matching Go backend models:
- User, Node, Server, Egg, Allocation, Ticket, CoinTransaction
- ApiResponse, PaginatedResponse
- PowerAction enum

## Known Limitations (Phase 1)

- Server create modal: allocation_id hardcoded to 1 (need allocation picker)
- Dashboard node/user counts are placeholders (need separate count endpoints)
- Server stats show 0 until Wings connection established
- No confirmation dialogs on some destructive actions

## Next Phase (Phase 3: Client Dashboard)

- User-facing server list
- Client server console (same component, different auth level)
- Coin management UI
- Ticket system
