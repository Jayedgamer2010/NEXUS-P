# NEXUS Dashboard - VPS Deployment Guide

Deploy the NEXUS Dashboard to your own VPS with Supabase as the database backend.

## Prerequisites

- A VPS with Ubuntu/Debian (tested on Ubuntu 22.04)
- Domain name pointing to your VPS (e.g., `dash.jsnexusp.online`)
- Supabase account and project
- Basic Linux command line knowledge

## Step 1: Prepare Your VPS

SSH into your VPS and run:

```bash
# Update system
apt update && apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com | sh
usermod -aG docker $USER
newgrp docker

# Install Docker Compose
apt install docker-compose-plugin -y

# Create project directory
mkdir -p /opt/nexus
```

## Step 2: Set Up Supabase

1. Go to https://supabase.com/dashboard and create a new project
2. Go to **SQL Editor** in your Supabase dashboard
3. Copy the contents of `nexus/supabase/schema.sql` and paste/run it
4. Go to **Project Settings → API** and copy:
   - **Project URL** (e.g., `https://xxxxx.supabase.co`)
   - **anon/public key**
   - **service_role key** (keep this secret — only for backend)

## Step 3: Upload Project Files

```bash
# On your VPS, clone the repo or upload files
cd /opt/nexus

# If using git:
git clone <your-repo-url> .

# Or scp from your machine:
# scp -r /path/to/NEXUS-P/* root@your-vps-ip:/opt/nexus/
```

## Step 4: Configure Environment

```bash
cd /opt/nexus/NEXUS-P/nexus
cp .env.example .env
nano .env
```

Edit the `.env` file with your values:

```env
# Supabase
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_SERVICE_KEY=your_supabase_service_role_key
SUPABASE_ANON_KEY=your_supabase_anon_key

# Pterodactyl (optional - only if you have Pterodactyl installed)
PTERODACTYL_URL=http://your-pterodactyl-url
PTERODACTYL_API_KEY=your-pterodactyl-api-key

# Redis
REDIS_PASSWORD=nexus_redis_password

# App
APP_NAME=NEXUS
```

## Step 5: Deploy Containers

```bash
cd /opt/nexus/NEXUS-P/nexus
docker compose up -d --build
```

Wait for it to finish building (takes 2-3 minutes). This will:
- Build PHP backend
- Build Vue frontend
- Start Redis container

## Step 6: Configure Nginx Reverse Proxy

### Remove any broken existing config
```bash
rm -f /etc/nginx/sites-enabled/mythicaldash
```

### Create the nginx config
```bash
nano /etc/nginx/sites-available/dash.jsnexusp.online
```

Paste this configuration:

```nginx
server {
    listen 80;
    listen [::]:80;
    server_name dash.jsnexusp.online;
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name dash.jsnexusp.online;

    ssl_certificate /etc/letsencrypt/live/panel.jsnexusp.online/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/panel.jsnexusp.online/privkey.pem;

    location / {
        proxy_pass http://127.0.0.1:4832;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

Save (Ctrl+X, Y, Enter).

### Create symlink and reload
```bash
ln -s /etc/nginx/sites-available/dash.jsnexusp.online /etc/nginx/sites-enabled/dash.jsnexusp.online
nginx -t
systemctl reload nginx
```

## Step 7: Get SSL Certificate

Since port 80 is used by Pterodactyl, we'll reuse its existing SSL certificate:

```bash
# Update nginx config to point to the cert (if needed)
sed -i "s|ssl_certificate.*;|ssl_certificate /etc/letsencrypt/live/panel.jsnexusp.online/fullchain.pem;|" /etc/nginx/sites-available/dash.jsnexusp.online
sed -i "s|ssl_certificate_key.*;|ssl_certificate_key /etc/letsencrypt/live/panel.jsnexusp.online/privkey.pem;|" /etc/nginx/sites-available/dash.jsnexusp.online

# Reload nginx
nginx -t && systemctl reload nginx
```

## Step 8: Point Your DNS

In your DNS provider (Cloudflare, Namecheap, etc.), create an **A record**:

```
Type: A
Name: dash
Value: <your VPS IP address>
TTL: Automatic
```

Wait 5-10 minutes for DNS propagation.

## Step 9: Verify Deployment

```bash
# Check containers are running
docker compose ps

# Test the site
curl -I https://dash.jsnexusp.online
```

Visit `https://dash.jsnexusp.online` in your browser - it should load the NEXUS dashboard!

## Step 10: Make First User Admin

After you sign up for the first time, run this in your Supabase SQL Editor:

```sql
update user_profiles set role = 'admin' where id = 'YOUR_USER_UUID';
```

(Find your UUID in Supabase → Authentication → Users)

## Directory Structure After Deployment

```
/opt/nexus/
├── NEXUS-P/                    # Project source
│   └── nexus/                  # Application code
│       ├── .env                # Environment variables (created above)
│       ├── docker-compose.yml  # Docker compose (from project)
│       ├── deploy/             # Deployment scripts
│       └── supabase/           # Database schema
└── docker-compose.yml          # Root compose (if using alternative method)
```

## Troubleshooting

- **502 Bad Gateway**: Check container logs with `docker compose logs -f`
- **SSL not working**: Verify `/etc/letsencrypt/live/panel.jsnexusp.online/` exists
- **Slow first load**: Containers may take 10-15 seconds to warm up on first request
- **Database connection errors**: Verify your SUPABASE_URL and keys in `.env`

## Maintenance

- View logs: `docker compose logs -f [service-name]`
- Update code: `git pull` then `docker compose up -d --build`
- Stop containers: `docker compose down`
- Restart containers: `docker compose restart`

Your NEXUS Dashboard is now live at https://dash.jsnexusp.online!