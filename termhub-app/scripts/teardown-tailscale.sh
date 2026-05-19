#!/usr/bin/env bash
# Uninstall TermHub Tailscale Funnel deployment
set -euo pipefail

echo "Stopping Tailscale Funnel..."
tailscale funnel 443 off 2>/dev/null || true
tailscale serve --https=443 off 2>/dev/null || true

echo "Stopping Docker containers..."
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$(dirname "$SCRIPT_DIR")"
docker compose down

echo ""
echo "✅ TermHub stopped and Funnel disabled."
echo ""
echo "To also remove data volumes:  docker compose down -v"
echo "To remove Tailscale node:     sudo tailscale logout"
