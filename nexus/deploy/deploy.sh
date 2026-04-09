#!/bin/bash
# NEXUS Dashboard - VPS Deploy Script
# Run this script on your VPS after copying the project
# Usage: bash deploy.sh

set -e

echo "=== NEXUS Dashboard Deploy ==="

# Check if running as root
if [ "$EUID" -ne 0 ]; then
  echo "Please run as root"
  exit 1
fi

# 1. Install Docker if not present
if ! command -v docker &> /dev/null; then
  echo "Installing Docker..."
  curl -fsSL https://get.docker.com | sh
fi

# 2. Install Docker Compose if not present
if ! docker compose version &> /dev/null; then
  echo "Installing Docker Compose..."
  apt install docker-compose-plugin -y
fi

# 3. Check .env
if [ ! -f .env ]; then
  echo ".env file not found. Copy from .env.example and configure first."
  echo "  cp .env.example .env"
  echo "  nano .env"
  exit 1
fi

# 4. Build and start
echo "Building and starting containers..."
docker compose down 2>/dev/null || true
docker compose up -d --build

# 5. Show status
echo ""
echo "=== Containers ==="
docker compose ps

echo ""
echo "=== Frontend: http://127.0.0.1:4832 ==="
echo ""
echo "Next: Configure nginx reverse proxy"
echo "  Run: bash deploy-nginx.sh"
echo ""
echo "Or manually: ln -s deploy/nginx.conf /etc/nginx/sites-enabled/nexus && nginx -t && systemctl reload nginx"
