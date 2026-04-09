# MythicalDash Deployment Guide for Pterodactyle Ports

## Prerequisites
- Ubuntu 20.04/22.04/24.04 server (or compatible Debian) running in your VPS.
- Pterodactyl panel already listening on **ports 80 and 443**.
- At least 2 GB RAM, 2 CPU cores, 20 GB storage.
- Root or sudo access.

---

## Step 1: System Update & Install Dependencies
```bash
sudo apt update && sudo apt upgrade -y
sudo apt install -y git curl wget nodejs npm unzip nginx mariadb-server redis-server php8.2-fpm
```

---

## Step 2: Clone the Repository
```bash
cd /var/www
sudo git clone https://github.com/your-username/dash.jsnexusp.online.git mythdash
cd mythdash
```
> Replace the URL with your fork's HTTPS URL.

---

## Step 3: Install the Correct Node.js Version (using nvm)
```bash
# Install nvm (Node Version Manager)
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash
export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \ . "$NVM_DIR/nvm.sh"
# Install Node 20 (LTS) and set as default
nvm install 20
nvm use 20
nvm alias default 20
```

---

## Step 4: MySQL / MariaDB Database Setup
```bash
sudo mysql -u root -p <<EOF
CREATE DATABASE mythdash CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'mythdash'@'localhost' IDENTIFIED BY 'StrongPassword123!';
GRANT ALL PRIVILEGES ON mythdash.* TO 'mythdash'@'localhost';
FLUSH PRIVILEGES;
EOF
```
> **IMPORTANT**: Replace `StrongPassword123!` with a strong, unique password.

---

## Step 5: Environment Configuration
```bash
cd /var/www/mythdash/backend
cp .env.example .env
nano .env
```
Add/update the following variables in `.env`:
```env
# Application
APP_URL=https://dash.yourdomain.com   # Or your server IP
APP_PORT=3000                         # Internal port MythicalDash will listen on

# Database (MySQL/MariaDB)
DB_TYPE=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USER=mythdash
DB_PASSWORD=StrongPassword123!   # Must match the password you set above
DB_NAME=mythdash

# Pterodactyl API (required for integration)
PTERODACTYL_BASE_URL=https://your-pterodactyl-domain.com
PTERODACTYL_API_KEY=YOUR_PTERODACTYL_API_KEY

# Security keys (generate strong random strings)
JWT_SECRET=$(openssl rand -hex 32)
ENCRYPTION_KEY=$(openssl rand -hex 32)
```
Save and exit.

---

## Step 6: Install Backend Dependencies
```bash
npm install
```
If you see any Supabase packages listed, they have already been removed in a previous step. The project now uses `mysql2`/`knex` (or the official ORM used by the reference implementation).

---

## Step 7: Run Database Migrations
The fork provides a migration script. Run it:
```bash
# From the backend directory
node database/migrate.js   # or, if the project uses Laravel artisan:
php artisan migrate --force
```
Check the console for any errors. All tables should now exist in the `mythdash` database.

---

## Step 8: Install & Build Front‑end
```bash
cd /var/www/mythdash/frontend
npm install
npm run build   # Produces a production build in ./dist (or ./build)
```
If the project uses Vite/React/Vue, the output folder may be `dist`. Adjust the path in the Nginx config accordingly.

---

## Step 9: Process Management with PM2
```bash
# Install PM2 globally
sudo npm install -g pm2

# Start the backend server (default entry is server.js – adjust if different)
cd /var/www/mythdash/backend
pm2 start server.js --name mythdash

# Save the process list so it restarts on boot
pm2 save

# Generate startup script for your init system (systemd)
pm pm2 startup systemd -u $USER --hp $HOME
```
Verify the app is running:
```bash
pm2 status
pm2 logs mythdash   # watch for any startup errors
```

---

## Step 10: Nginx Reverse Proxy (Ports 3000 & 3001)
Because ports **80** and **443** are occupied by Pterodactyl, we expose MythicalDash on **port 3001** (HTTPS can be on 8443 if desired).

```bash
sudo nano /etc/nginx/sites-available/mythdash
```
Paste the following configuration (replace `dash.yourdomain.com` with your actual domain or use the server IP):
```nginx
server {
    listen 3001;
    server_name dash.yourdomain.com;

    location / {
        proxy_pass http://127.0.0.1:3000;   # Backend internal port
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}
```
Enable the site and reload Nginx:
```bash
sudo ln -s /etc/nginx/sites-available/mythdash /etc/nginx/sites-enabled/
sudo nginx -t && sudo systemctl reload nginx
```
Now the dashboard is reachable at `http://dash.yourdomain.com:3001`.

---

## Step 11: (Optional) SSL on a Non‑Standard Port
Since 443 is taken, you can serve HTTPS on **8443**:
```bash
sudo nano /etc/nginx/sites-available/mythdash-ssl
```
Add:
```nginx
server {
    listen 8443 ssl http2;
    server_name dash.yourdomain.com;

    ssl_certificate /etc/letsencrypt/live/dash.yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/dash.yourdomain.com/privkey.pem;

    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```
Enable and reload:
```bash
sudo ln -s /etc/nginx/sites-available/mythdash-ssl /etc/nginx/sites-enabled/
sudo nginx -t && sudo systemctl reload nginx
```
Access via `https://dash.yourdomain.com:8443`.

---

## Step 12: Firewall Rules
Open the new ports while keeping the existing ones for Pterodactyl:
```bash
sudo ufw allow 3000/tcp   # Backend internal port (optional, only for testing)
sudo ufw allow 3001/tcp   # HTTP reverse‑proxy port
sudo ufw allow 8443/tcp   # HTTPS reverse‑proxy port (if you set up SSL)
sudo ufw reload
```

---

## Step 13: Verify the Installation
```bash
# Verify backend is listening
sudo netstat -tlnp | grep 3000

# Verify Nginx proxy works
curl -I http://your-server-ip:3001
```
You should receive a `200 OK` response and the JSON API payload from MythicalDash.

---

## Step 14: Maintenance Commands
```bash
# Restart the app
pm2 restart mythdash

# View logs
pm2 logs mythdash

# Stop the app
pm2 stop mythdash

# Remove the app from PM2
pm2 delete mythdash
```

---

## Troubleshooting Tips
- **Database connection errors** – double‑check `.env` DB credentials and that the MySQL service is running (`systemctl status mariadb`).
- **Port conflicts** – run `sudo lsof -i -P -n | grep LISTEN` to see what services occupy each port.
- **Nginx 502 Bad Gateway** – ensure PM2 reports the process as `online` and that `server.js` is listening on port 3000.
- **Pterodactyl integration failures** – confirm the API key and base URL are correct and that the panel allows API access.

---

## File Structure Overview
```
/var/www/mythdash/
├─ backend/          # PHP/Node backend
│   ├─ server.js
│   ├─ .env          # Database & API credentials
│   └─ database/...
├─ frontend/         # Vue/React source
│   ├─ package.json
│   └─ dist/          # Production build output
└─ README.md
```

---

## Support & Contributions
- Open an issue on the GitHub repo: https://github.com/your-username/dash.jsnexusp.online/issues
- Refer to the official MythicalDash documentation: https://docs.mythicalsystems.gq

---

*End of deployment guide.*