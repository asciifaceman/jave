#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SITE_DIR="${ROOT_DIR}/site"
REF_DIR="${SITE_DIR}/reference"
STATIC_DIR="${SITE_DIR}/static"
LORE_FEED_FILE="${SITE_DIR}/lore/feed.md"

rm -rf "${REF_DIR}"
mkdir -p "${REF_DIR}" "${STATIC_DIR}" "${SITE_DIR}/lore"

cp -R "${ROOT_DIR}/governance" "${REF_DIR}/"
cp -R "${ROOT_DIR}/docs" "${REF_DIR}/"
cp -R "${ROOT_DIR}/specs" "${REF_DIR}/"
cp -R "${ROOT_DIR}/examples" "${REF_DIR}/"

if [[ -f "${ROOT_DIR}/docs/static/logo-trnsprnt.png" ]]; then
  cp "${ROOT_DIR}/docs/static/logo-trnsprnt.png" "${STATIC_DIR}/logo-trnsprnt.png"
fi

cat > "${LORE_FEED_FILE}" <<'EOF'
---
layout: page
title: Lore Feed
permalink: /lore/feed/
---

This feed is generated from governance lore records on each merge to `main`.

EOF

append_lore_section() {
  local source_dir="$1"
  local section_title="$2"

  {
    echo "## ${section_title}"
    echo
  } >> "${LORE_FEED_FILE}"

  mapfile -t files < <(find "${source_dir}" -maxdepth 1 -type f -name '*.md' | sort)

  if [[ ${#files[@]} -eq 0 ]]; then
    echo "- No records published yet." >> "${LORE_FEED_FILE}"
    echo >> "${LORE_FEED_FILE}"
    return
  fi

  for file in "${files[@]}"; do
    local rel_path="${file#${ROOT_DIR}/}"
    local label
    label="$(basename "${file}" .md)"
    echo "- [${label}]({{ '/reference/${rel_path}' | relative_url }})" >> "${LORE_FEED_FILE}"
  done

  echo >> "${LORE_FEED_FILE}"
}

append_lore_section "${ROOT_DIR}/governance/v0.1/lore" "v0.1 Lore Records"
append_lore_section "${ROOT_DIR}/governance/v0.2/lore" "v0.2 Lore Records"
