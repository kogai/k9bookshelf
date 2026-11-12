#!/bin/bash
# Fails when tracked files (or, with --staged, the staged diff) contain something that looks like a secret.
# This repository is public: never commit credentials, tokens, keys or .env files.
set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

# Shopify (shpat_/shppa_/shpss_/shpca_), AWS, GitHub, Slack, Google API keys, OpenAI/Anthropic style keys and private keys.
pattern='shp(at|pa|ss|ca)_[0-9a-fA-F]{20,}|AKIA[0-9A-Z]{16}|gh[pousr]_[0-9A-Za-z]{30,}|github_pat_[0-9A-Za-z_]{20,}|xox[abpr]-[0-9A-Za-z-]{10,}|AIza[0-9A-Za-z_-]{35}|sk-(ant-)?[A-Za-z0-9_-]{32,}|-----BEGIN [A-Z ]*PRIVATE KEY-----|"private_key": *"'
# Files which must never be committed.
forbidden='(^|/)\.env(\..*)?$|\.pem$|\.p12$|\.key$|(^|/)id_(rsa|ed25519)$|(^|/)secrets\.json$|service[-_]?account.*\.json$'

status=0
if [ "${1:-}" = "--staged" ]; then
  files=$(git diff --cached --name-only --diff-filter=ACMR)
  if hits=$(git diff --cached -U0 --no-color | grep -E '^\+' | grep -nE "$pattern"); then
    echo "check-secrets: staged changes look like they contain a secret:" >&2
    echo "$hits" | cut -c1-120 >&2
    status=1
  fi
else
  files=$(git ls-files)
  if hits=$(git grep -nIE "$pattern" -- . ':!yarn.lock'); then
    echo "check-secrets: tracked files look like they contain a secret:" >&2
    echo "$hits" | cut -c1-120 >&2
    status=1
  fi
fi

if bad=$(echo "$files" | grep -E "$forbidden"); then
  echo "check-secrets: these files must not be committed:" >&2
  echo "$bad" >&2
  status=1
fi

exit $status
