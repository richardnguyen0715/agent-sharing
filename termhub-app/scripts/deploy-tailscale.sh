#!/usr/bin/env bash
# TermHub — Tailscale Funnel Deployment Script
# Installs Tailscale, deploys the Docker stack, and exposes it to the public internet.
#
# Usage:
#   chmod +x scripts/deploy-tailscale.sh
#   sudo ./scripts/deploy-tailscale.sh
#
# Prerequisites:
#   - Ubuntu/Debian Linux server
#   - Docker Engine + Docker Compose plugin installed
#   - A Tailscale account (https://tailscale.com)

set -uo pipefail
# Note: not using -e — we handle errors explicitly to avoid silent exits.

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

log()  { echo -e "${GREEN}[✓]${NC} $*"; }
warn() { echo -e "${YELLOW}[!]${NC} $*"; }
err()  { echo -e "${RED}[✗]${NC} $*" >&2; }
info() { echo -e "${CYAN}[i]${NC} $*"; }

# ── Preflight checks ───────────────────────────────────────────

if [[ $EUID -ne 0 ]]; then
    err "This script must be run as root (use sudo)"
    exit 1
fi

if ! command -v docker &>/dev/null; then
    err "Docker is not installed. Install it first:"
    echo "  curl -fsSL https://get.docker.com | sh"
    exit 1
fi

if ! docker compose version &>/dev/null; then
    err "Docker Compose plugin is not installed."
    echo "  sudo apt-get install docker-compose-plugin"
    exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_DIR"

if [[ ! -f docker-compose.yml ]]; then
    err "docker-compose.yml not found in $PROJECT_DIR"
    exit 1
fi

# ── Step 1: Install Tailscale ───────────────────────────────────

if ! command -v tailscale &>/dev/null; then
    info "Installing Tailscale..."
    curl -fsSL https://tailscale.com/install.sh | sh
    log "Tailscale installed"
else
    log "Tailscale already installed ($(tailscale version | head -1))"
fi

# ── Step 2: Authenticate Tailscale ──────────────────────────────

TS_STATUS=$(tailscale status --json 2>/dev/null | grep -o '"BackendState":"[^"]*"' | cut -d'"' -f4 || true)

if [[ "$TS_STATUS" != "Running" ]]; then
    info "Tailscale is not logged in. Starting authentication..."
    echo ""
    echo "  A login URL will appear below. Open it in your browser to authenticate."
    echo ""
    tailscale up --hostname=termhub || true
    # Wait for Tailscale to fully connect
    sleep 3
    log "Tailscale authenticated"
else
    log "Tailscale is already running"
fi

# Get the Tailscale FQDN
TS_FQDN=$(tailscale status --json 2>/dev/null | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    me = data.get('Self', {})
    dns = me.get('DNSName', '')
    print(dns.rstrip('.'))
except:
    print('')
" 2>/dev/null || true)

if [[ -z "$TS_FQDN" ]]; then
    warn "Could not auto-detect Tailscale FQDN. Trying hostname-based fallback..."
    TS_FQDN=$(tailscale whois --json $(tailscale ip -4) 2>/dev/null | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    print(data.get('Node', {}).get('Name', '').rstrip('.'))
except:
    print('')
" 2>/dev/null || true)
fi

if [[ -z "$TS_FQDN" ]]; then
    err "Could not determine Tailscale FQDN"
    echo "  Run 'tailscale status' to check your connection."
    echo "  Then set CORS_ORIGINS manually in .env and re-run."
    exit 1
fi

log "Tailscale hostname: ${CYAN}${TS_FQDN}${NC}"

# ── Step 3: Configure environment ──────────────────────────────

if [[ ! -f .env ]]; then
    cp .env.example .env
    warn "Created .env from template — generating secure secrets..."
fi

# Generate secure secrets if still using defaults
if grep -q "change-me" .env 2>/dev/null; then
    JWT_SECRET=$(openssl rand -base64 32)
    sed -i "s|JWT_SECRET=.*|JWT_SECRET=${JWT_SECRET}|" .env
    log "Generated secure JWT_SECRET"
fi

if grep -q "POSTGRES_PASSWORD=change-me\|POSTGRES_PASSWORD=termhub" .env 2>/dev/null; then
    PG_PASS=$(openssl rand -base64 16 | tr -d '=+/')
    sed -i "s|POSTGRES_PASSWORD=.*|POSTGRES_PASSWORD=${PG_PASS}|" .env
    log "Generated secure POSTGRES_PASSWORD"
fi

# Set CORS to Tailscale FQDN
if ! grep -q "^CORS_ORIGINS=https://${TS_FQDN}" .env 2>/dev/null; then
    if grep -q "^CORS_ORIGINS=" .env; then
        sed -i "s|^CORS_ORIGINS=.*|CORS_ORIGINS=https://${TS_FQDN}|" .env
    elif grep -q "^# *CORS_ORIGINS=" .env; then
        sed -i "s|^# *CORS_ORIGINS=.*|CORS_ORIGINS=https://${TS_FQDN}|" .env
    else
        echo "CORS_ORIGINS=https://${TS_FQDN}" >> .env
    fi
    log "Set CORS_ORIGINS to https://${TS_FQDN}"
fi

# ── Step 4: Build & start Docker stack ─────────────────────────

info "Building and starting TermHub..."
docker compose up --build -d

# Wait for health
info "Waiting for services to be healthy..."
for i in {1..30}; do
    if curl -sf http://localhost:3000/api/health &>/dev/null; then
        log "TermHub is healthy"
        break
    fi
    if [[ $i -eq 30 ]]; then
        err "Timed out waiting for TermHub to start"
        echo "  Check logs: docker compose logs"
        exit 1
    fi
    sleep 2
done

# ── Step 5: Configure Tailscale Funnel ─────────────────────────

info "Configuring Tailscale Serve (reverse proxy to port 3000)..."
tailscale serve --bg --https=443 http://localhost:3000

info "Enabling Tailscale Funnel (public internet access)..."
tailscale funnel --bg 443 on

log "Tailscale Funnel is active"

# ── Done ───────────────────────────────────────────────────────

echo ""
echo "╔══════════════════════════════════════════════════════════════╗"
echo "║                                                              ║"
echo "║   ${GREEN}TermHub is live!${NC}                                          ║"
echo "║                                                              ║"
echo "║   ${CYAN}Public URL:${NC}  https://${TS_FQDN}              "
echo "║                                                              ║"
echo "║   Default login:  admin / admin123                           ║"
echo "║   ${YELLOW}⚠  Change the admin password after first login!${NC}           ║"
echo "║                                                              ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo ""
echo "Useful commands:"
echo "  docker compose logs -f          # View logs"
echo "  docker compose ps               # Container status"
echo "  tailscale status                 # Tailscale status"
echo "  tailscale funnel status          # Funnel status"
echo "  sudo tailscale funnel 443 off    # Disable public access"
echo "  docker compose down              # Stop TermHub"
echo ""
