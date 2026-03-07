#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SITE_DIR="${ROOT_DIR}/site"
CONTENT_DIR="${SITE_DIR}/content"
STATIC_DIR="${SITE_DIR}/static"
DOCS_STATIC_DIR="${SITE_DIR}/docs/static"
RECORDS_FEED_FILE="${SITE_DIR}/records/feed.md"

rm -rf "${CONTENT_DIR}"
mkdir -p "${CONTENT_DIR}" "${STATIC_DIR}" "${DOCS_STATIC_DIR}" "${SITE_DIR}/records"

trim_title() {
  local t="$1"
  t="${t#\# }"
  t="${t#\## }"
  printf '%s' "${t}"
}

infer_title_from_file() {
  local src="$1"
  local first_header
  first_header="$(grep -m1 -E '^# ' "${src}" || true)"
  if [[ -n "${first_header}" ]]; then
    trim_title "${first_header}"
    return
  fi
  basename "${src}" .md
}

infer_summary_from_file() {
  local src="$1"
  local summary
  summary="$(grep -m1 -E '^[A-Za-z0-9].+' "${src}" || true)"
  printf '%s' "${summary}"
}

render_markdown_document() {
  local src="$1"
  local rel_path="$2"
  local out_path="${CONTENT_DIR}/${rel_path}"
  local permalink_path="/content/${rel_path%.md}/"
  local title
  local summary

  title="$(infer_title_from_file "${src}")"
  summary="$(infer_summary_from_file "${src}")"

  mkdir -p "$(dirname "${out_path}")"
  {
    echo "---"
    echo "layout: page"
    echo "title: ${title}"
    echo "permalink: ${permalink_path}"
    echo "---"
    echo
    echo "<div class=\"doc-meta\">"
    echo "  <p><strong>Document Class:</strong> Governance or technical archive record</p>"
    if [[ -n "${summary}" ]]; then
      echo "  <p><strong>Summary:</strong> ${summary}</p>"
    fi
    echo "</div>"
    echo
    awk 'NR==1 && /^# / {next} {print}' "${src}"
    echo
  } > "${out_path}"
}

render_example_document() {
  local src="$1"
  local rel_path="$2"
  local out_path="${CONTENT_DIR}/${rel_path}.md"
  local permalink_path="/content/${rel_path}/"
  local title

  title="$(basename "${src}")"

  mkdir -p "$(dirname "${out_path}")"
  {
    echo "---"
    echo "layout: page"
    echo "title: ${title}"
    echo "permalink: ${permalink_path}"
    echo "---"
    echo
    echo "<div class=\"doc-meta\">"
    echo "  <p><strong>Source Path:</strong> <code>${rel_path}</code></p>"
    echo "  <p><strong>Document Type:</strong> Jave executable example</p>"
    echo "</div>"
    echo
    echo '```jave'
    cat "${src}"
    echo
    echo '```'
  } > "${out_path}"
}

while IFS= read -r src; do
  rel="${src#${ROOT_DIR}/}"
  render_markdown_document "${src}" "${rel}"
done < <(find "${ROOT_DIR}/governance" "${ROOT_DIR}/docs" "${ROOT_DIR}/specs" -type f -name '*.md' | sort)

while IFS= read -r src; do
  rel="${src#${ROOT_DIR}/}"
  rel_no_ext="${rel%.jave}"
  render_example_document "${src}" "${rel_no_ext}"
done < <(find "${ROOT_DIR}/examples" -type f -name '*.jave' | sort)

if [[ -f "${ROOT_DIR}/docs/static/logo-trnsprnt.png" ]]; then
  cp "${ROOT_DIR}/docs/static/logo-trnsprnt.png" "${STATIC_DIR}/logo-trnsprnt.png"
  cp "${ROOT_DIR}/docs/static/logo-trnsprnt.png" "${DOCS_STATIC_DIR}/logo-trnsprnt.png"
fi

cat > "${RECORDS_FEED_FILE}" <<'EOF'
---
layout: page
title: Records Feed
permalink: /records/feed/
---

This feed is generated from governance record archives on each merge to `main`.

EOF

append_record_section() {
  local source_dir="$1"
  local section_title="$2"

  {
    echo "## ${section_title}"
    echo
  } >> "${RECORDS_FEED_FILE}"

  mapfile -t files < <(find "${source_dir}" -maxdepth 1 -type f -name '*.md' | sort)

  if [[ ${#files[@]} -eq 0 ]]; then
    echo "- No records published yet." >> "${RECORDS_FEED_FILE}"
    echo >> "${RECORDS_FEED_FILE}"
    return
  fi

  for file in "${files[@]}"; do
    local rel_path="${file#${ROOT_DIR}/}"
    local display_path="${rel_path%.md}"
    local title
    title="$(infer_title_from_file "${file}")"
    echo "- [${title}]({{ '/content/${display_path}/' | relative_url }})" >> "${RECORDS_FEED_FILE}"
  done

  echo >> "${RECORDS_FEED_FILE}"
}

append_record_section "${ROOT_DIR}/governance/v0.1/lore" "v0.1 Governance Records"
append_record_section "${ROOT_DIR}/governance/v0.2/lore" "v0.2 Governance Records"
