#!/usr/bin/env bash
set -euo pipefail

TAG="$(date -u +%Y.%m.%d)"
if git rev-parse -q --verify "refs/tags/${TAG}" >/dev/null; then
  echo "Release tag ${TAG} already exists; only one release per UTC day is supported." >&2
  exit 1
fi

LAST_TAG="$(git for-each-ref refs/tags --sort=-creatordate --format='%(refname:short)' | head -n 1 || true)"

{
  echo "tag=$TAG"
  echo "release_name=$TAG"
  echo "last_tag=$LAST_TAG"
} >> "${GITHUB_OUTPUT:?GITHUB_OUTPUT is required}"

echo "Release tag: $TAG"
echo "Previous tag: ${LAST_TAG:-none}"
