#!/bin/bash
# PreToolUse hook: blocks `git commit` when the staged changes look like they contain a secret.
# This repository is public.
set -uo pipefail

input=$(cat)
if ! printf '%s' "$input" | grep -qE 'git[^"]*commit'; then
  exit 0
fi

cd "${CLAUDE_PROJECT_DIR:-.}" || exit 0
# NOTE: `git commit -a` commits tracked changes which aren't staged yet, so check the whole tree too.
if ! out=$(./scripts/check-secrets.sh --staged 2>&1 && ./scripts/check-secrets.sh 2>&1); then
  echo "$out" >&2
  echo "Commit blocked: this repository is public. Remove the secret (and rotate it if it was ever pushed)." >&2
  exit 2
fi
exit 0
