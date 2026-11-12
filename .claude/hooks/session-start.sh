#!/bin/bash
# Installs dependencies so that `make test` and `make check/theme` work in Claude Code cloud sessions.
set -euo pipefail

if [ "${CLAUDE_CODE_REMOTE:-}" != "true" ]; then
  exit 0
fi

cd "${CLAUDE_PROJECT_DIR:-.}"

# Go modules are downloaded outside the repository, so the hash file lives in the module cache.
gohashes="$(go env GOMODCACHE)/.k9bookshelf-session-start.sha256"
if ! sha256sum --check --status "$gohashes" 2>/dev/null; then
  rm -f "$gohashes"
  go mod download
  mkdir -p "$(dirname "$gohashes")"
  sha256sum .claude/hooks/session-start.sh go.mod go.sum > "$gohashes" || rm -f "$gohashes"
fi

# Node.js packages for building and linting the theme.
nodehashes=node_modules/.session-start.sha256
if ! sha256sum --check --status "$nodehashes" 2>/dev/null; then
  rm -f "$nodehashes"
  yarn install --frozen-lockfile --non-interactive
  sha256sum .claude/hooks/session-start.sh package.json yarn.lock > "$nodehashes" || rm -f "$nodehashes"
fi
