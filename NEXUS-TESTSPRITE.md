# NEXUS Panel - Test Specification

## Application Overview
NEXUS Panel is a Pterodactyl Panel admin dashboard built with React + Vite.

## Test Scenarios

### Auth - Login (/login)
- Elements: username/email input, password input, login button, link to register
- Test: login button clickable (not disabled)
- Test: form validation (empty fields show error)
- Test: links to register page

### Auth - Register (/register)
- Elements: username input, email input, password input, confirm password, register button, link to login
- Test: register button is disabled when form invalid
- Test: register button becomes enabled when all fields valid
- Test: passwords must match validation

### Admin Dashboard (/admin/dashboard)
- Elements: 4 stat cards (Total Servers, Running Server, Active Nodes, Active Users), recent activity table
- Test: stat cards are visible
- Test: recent activity table loads

### Admin Servers (/admin/servers)
- Elements: data table, Create Server button
- Test: table loads with server list
- Test: Create Server button opens modal
- Test: modal has required fields (name, node, egg, memory, disk, cpu)
- Test: status badges show correct colors (green=running, red=stopped, yellow=installing)

### Admin Server Detail (/admin/servers/:id)
- Elements: server name, status badge, power buttons (Start/Stop/Restart/Kill), 3 stat cards, console terminal, command input
- Test: power buttons are clickable
- Test: console terminal is visible
- Test: stat cards show CPU/RAM/Disk

### Admin Nodes (/admin/nodes)
- Elements: data table, Create Node button
- Test: table loads
- Test: Create Node button opens modal

### Admin Users (/admin/users)
- Elements: data table, Create User button
- Test: table loads
- Test: Create User button opens modal

### Admin Eggs (/admin/eggs)
- Elements: data table, Create Egg button
- Test: table loads

## UI Components to Verify
- Sidebar: visible on all admin pages, correct active state highlighting
- Header: visible on all admin pages
- Status badges: correct color coding
- Modals: open and close correctly, ESC key closes modal
- Data tables: pagination controls visible if more than 10 items
- Forms: all required fields marked, submit disabled until valid

## Design Verification
- Background color: #0a0a0f (dark)
- Primary accent: #7c3aed (purple)
- Card surfaces: #1a1a2e
- Font: Inter
- No elements overflowing viewport on mobile (375px width)
- No elements overflowing viewport on desktop (1440px width)

## API Endpoints to Test
- POST /api/auth/register → expect 200 with token
- POST /api/auth/login → expect 200 with token
- GET /api/auth/me → expect 200 with user object (requires Bearer token)
- GET /api/admin/servers → expect 200 (requires admin Bearer token)
- GET /api/admin/nodes → expect 200 (requires admin Bearer token)
- GET /api/admin/users → expect 200 (requires admin Bearer token)

## Expected User Flow
1. Visit /login
2. Register new account at /register
3. Login with credentials
4. Land on /admin/dashboard
5. Create a node via /admin/nodes
6. Create an egg via /admin/eggs
7. Create a server via /admin/servers
8. View server detail and test power controls

## Known Limitations for Testing
- WebSocket console requires active Wings daemon connection
- Live stats require Wings daemon running on the node
- First registered user automatically gets admin role
