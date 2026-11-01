#!/usr/bin/env bash
set -euo pipefail

: "${GITHUB_OUTPUT:?GITHUB_OUTPUT is required}"
: "${RELEASE_NAME:?RELEASE_NAME is required}"

OUTPUT_FILE="$(mktemp)"
FIX_FILE="$(mktemp)"
FEAT_FILE="$(mktemp)"
REFACTOR_FILE="$(mktemp)"
PERF_FILE="$(mktemp)"
STYLE_FILE="$(mktemp)"
CHORE_FILE="$(mktemp)"
DOCS_FILE="$(mktemp)"
CI_FILE="$(mktemp)"
TEST_FILE="$(mktemp)"
OTHER_FILE="$(mktemp)"
trap 'rm -f "$OUTPUT_FILE" "$FIX_FILE" "$FEAT_FILE" "$REFACTOR_FILE" "$PERF_FILE" "$STYLE_FILE" "$CHORE_FILE" "$DOCS_FILE" "$CI_FILE" "$TEST_FILE" "$OTHER_FILE"' EXIT

{
  echo "# 🪟 Build Name: ${RELEASE_NAME}"
  echo
  echo "## 🌟 What's Changed"
  echo
} > "$OUTPUT_FILE"

COMMIT_REGEX='^([a-zA-Z]+)(\(([^)]+)\))?!?: (.+)$'

if [[ -n "${LAST_TAG:-}" ]] && git rev-parse -q --verify "${LAST_TAG}^{commit}" >/dev/null; then
  COMMITS="$(git log "${LAST_TAG}..HEAD" --pretty=tformat:'%H|%s' --no-merges)"
else
  COMMITS="$(git log HEAD --pretty=tformat:'%H|%s' --no-merges)"
fi

while IFS='|' read -r hash message; do
  [[ -n "$hash" ]] || continue
  short_hash="${hash:0:7}"

  if [[ "$message" =~ $COMMIT_REGEX ]]; then
    type="${BASH_REMATCH[1]}"
    scope="${BASH_REMATCH[3]}"
    subject="${BASH_REMATCH[4]}"
    if [[ -n "$scope" ]]; then
      line="- ${type}(${scope}): ${subject} (\`${short_hash}\`)"
    else
      line="- ${type}: ${subject} (\`${short_hash}\`)"
    fi

    case "${type,,}" in
      feat|remove) printf '%s\n' "$line" >> "$FEAT_FILE" ;;
      fix)         printf '%s\n' "$line" >> "$FIX_FILE" ;;
      chore)       printf '%s\n' "$line" >> "$CHORE_FILE" ;;
      docs)        printf '%s\n' "$line" >> "$DOCS_FILE" ;;
      style)       printf '%s\n' "$line" >> "$STYLE_FILE" ;;
      refactor)    printf '%s\n' "$line" >> "$REFACTOR_FILE" ;;
      perf)        printf '%s\n' "$line" >> "$PERF_FILE" ;;
      ci)          printf '%s\n' "$line" >> "$CI_FILE" ;;
      test)        printf '%s\n' "$line" >> "$TEST_FILE" ;;
      *)           printf '%s\n' "$line" >> "$OTHER_FILE" ;;
    esac
  else
    printf '%s\n' "- ${message} (\`${short_hash}\`)" >> "$OTHER_FILE"
  fi
done <<< "$COMMITS"

append_section() {
  local file="$1"
  local title="$2"
  if [[ -s "$file" ]]; then
    printf '### %s\n\n' "$title" >> "$OUTPUT_FILE"
    cat "$file" >> "$OUTPUT_FILE"
    printf '\n' >> "$OUTPUT_FILE"
  fi
}

append_section "$FIX_FILE" "🐛 Fixed"
append_section "$FEAT_FILE" "➕ Added or Removed"
append_section "$REFACTOR_FILE" "♻️ Refactored"
append_section "$PERF_FILE" "⚡ Performance"
append_section "$STYLE_FILE" "💄 Style"
append_section "$CHORE_FILE" "🔧 Chores"
append_section "$DOCS_FILE" "📚 Documentation"
append_section "$CI_FILE" "🔄 CI"
append_section "$TEST_FILE" "✅ Tests"
append_section "$OTHER_FILE" "📦 Other"

DELIMITER="RELEASE_NOTES_$(date +%s)_$$"
{
  echo "body<<${DELIMITER}"
  cat "$OUTPUT_FILE"
  echo "$DELIMITER"
} >> "$GITHUB_OUTPUT"
