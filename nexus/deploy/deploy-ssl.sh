#!/bin/bash
# Install SSL certificate for dash.jsnexusp.online using existing Pterodactyl cert
# Run as root on your VPS

set -e

echo "=== SSL Certificate Setup ==="

# Check if certbot exists
if ! command -v certbot &> /dev/null; then
  echo "Installing certbot..."
  apt install certbot python3-certbot-nginx -y
fi

# Check if Pterodactyl cert exists (reuse it)
PT_CERT="/etc/letsencrypt/live/panel.jsnexusp.online/fullchain.pem"
if [ -f "$PT_CERT" ]; then
  echo "Reusing existing Pterodactyl SSL certificate..."
  echo "  Cert: $PT_CERT"
  echo "  Key: /etc/letsencrypt/live/panel.jsnexusp.online/privkey.pem"

  # Update nginx config to use the cert
  sed -i "s|ssl_certificate.*;|ssl_certificate $PT_CERT;|" /etc/nginx/sites-available/nexus
  sed -i "s|ssl_certificate_key.*;|ssl_certificate_key /etc/letsencrypt/live/panel.jsnexusp.online/privkey.pem;|" /etc/nginx/sites-available/nexus

  nginx -t && systemctl reload nginx
  echo ""
  echo "SSL configured using existing Pterodactyl certificate."
  echo "Visit: https://dash.jsnexusp.online"
  exit 0
fi

# Otherwise, get a new cert for dash.jsnexusp.online
echo "Requesting new SSL certificate for dash.jsnexusp.online..."
certbot --nginx -d dash.jsnexusp.online --non-interactive --agree-tos --email Junayedsheikh749@gmail.com

nginx -t && systemctl reload nginx
echo ""
echo "SSL installed!"
echo "Visit: https://dash.jsnexusp.online"
