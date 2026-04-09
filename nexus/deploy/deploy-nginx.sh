#!/bin/bash
# Setup Nginx Reverse Proxy for dash.jsnexusp.online
# Run as root on your VPS

set -e

echo "=== Nginx Reverse Proxy Setup ==="

# Check nginx
if ! command -v nginx &> /dev/null; then
  echo "Nginx is not installed!"
  exit 1
fi

# Remove broken existing configs
for f in /etc/nginx/sites-enabled/mythicaldash; do
  if [ -f "$f" ] || [ -L "$f" ]; then
    echo "Removing broken config: $f"
    rm -f "$f"
  fi
done

# Create symlink
NGINX_SRC="$(pwd)/deploy/nginx.conf"
NGINX_DST="/etc/nginx/sites-available/nexus"
SITES_ENABLED="/etc/nginx/sites-enabled"

cp "$NGINX_SRC" "$NGINX_DST"
rm -f "$SITES_ENABLED/nexus"
ln -s "$NGINX_DST" "$SITES_ENABLED/nexus"

# Test and reload
echo "Testing nginx config..."
nginx -t

echo "Reloading nginx..."
systemctl reload nginx

echo ""
echo "Done! Nginx is now proxying dash.jsnexusp.online to the dashboard."
echo ""
echo "To add SSL, run: bash deploy/deploy-ssl.sh"
