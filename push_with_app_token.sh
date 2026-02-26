#!/bin/bash
set -euo pipefail

APP_ID=2862166
INSTALLATION_ID=110026444
PRIVATE_KEY_PATH="$HOME/.openclaw/credentials/github.pem"
REPO_URL="https://github.com/sirdeggen/go-authsocket.git"

REPO_ROOT=$(pwd)
if [ ! -d ".git" ]; then
  echo "Error: This script must be run from the repo root that has a Git history." >&2
  exit 1
fi

# Helpers: base64url
b64url() {
  openssl base64 -e -A | tr '+/' '-_' | tr -d '='
}

# JWT header/payload — GitHub requires exp <= now + 600s
NOW=$(date +%s)
HEADER='{"alg":"RS256","typ":"JWT"}'
PAYLOAD="{\"iat\":${NOW},\"exp\":$((NOW+600)),\"iss\":${APP_ID}}"

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

# Extract token — use python3 with proper stdin piping
TOKEN=$(printf '%s' "$RESPONSE" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('token',''))")

if [ -z "$TOKEN" ]; then
  echo "Failed to obtain installation token. Response:"
  printf '%s\n' "$RESPONSE"
  exit 1
fi

echo "Obtained installation token. Pushing to main..."
REMOTE_URL="https://x-access-token:${TOKEN}@${REPO_URL#https://}"
GIT_TERMINAL_PROMPT=0 git -c credential.helper= push "$REMOTE_URL" main
