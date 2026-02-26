#!/bin/bash
set -euo pipefail

APP_ID=2862166
INSTALLATION_ID=110026444
PRIVATE_KEY_PATH="$HOME/.openclaw/credentials/github.pem"
REPO_URL="https://github.com/sirdeggen/go-authsocket.git"

# Ensure you’re in the repo root
REPO_ROOT=$(pwd)
if [ ! -d ".git" ]; then
  echo "Error: This script must be run from the repo root that has a Git history." >&2
  exit 1
fi

# Helpers: base64url
b64url() {
  openssl base64 -e -A | tr '+/' '-_' | tr -d '='
}

# JWT header/payload
NOW=$(date +%s)
HEADER='{"alg":"RS256","typ":"JWT"}'
PAYLOAD="{\"iat\":${NOW},\"exp\":$((${NOW}+6000)),\"iss\":${APP_ID}}"

HEADER_B64=$(printf "%s" "$HEADER" | b64url)
PAYLOAD_B64=$(printf "%s" "$PAYLOAD" | b64url)

DATA="$HEADER_B64.$PAYLOAD_B64"

# Sign with private key (RS256)
SIGNATURE=$(printf "%s" "$DATA" | openssl dgst -sha256 -sign "$PRIVATE_KEY_PATH" -binary | openssl base64 -A | tr '+/' '-_' | tr -d '=')

JWT="$DATA.$SIGNATURE"

# Get installation access token
RESPONSE=$(curl -sS -X POST \
  -H "Authorization: Bearer $JWT" \
  -H "Accept: application/vnd.github.machine-man-preview+json" \
  "https://api.github.com/app/installations/${INSTALLATION_ID}/access_tokens")

TOKEN=$(echo "$RESPONSE" | python3 - << 'PY'
import sys, json
d = json.load(sys.stdin)
print(d.get("token", "")) if "token" in d else print("")
PY
)

if [ -z "$TOKEN" ]; then
  echo "Failed to obtain installation token. Response:"
  echo "$RESPONSE" | sed -n '1,200p'
  exit 1
fi

echo "Obtained installation token. Pushing to main..."
REMOTE_URL="https://${TOKEN}@${REPO_URL#https://}"
git remote set-url origin "$REMOTE_URL" >/dev/null 2>&1 || true
git push -u origin main
