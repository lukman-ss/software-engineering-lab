#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="${REPO_ROOT:-$HOME/Documents/LUKMAN/software-engineering-lab}"
AGENT_ROOT="${AGENT_ROOT:-$HOME/research-agent}"

PROMPTS="$AGENT_ROOT/prompts"
TASK_ROOT="$AGENT_ROOT/tasks"
TASK_PENDING="$TASK_ROOT/pending"
TASK_DONE="$TASK_ROOT/done"
LOG_ROOT="$AGENT_ROOT/logs"

MAX_RESEARCH_REVISIONS="${MAX_RESEARCH_REVISIONS:-3}"
MAX_ENGINEERING_REVISIONS="${MAX_ENGINEERING_REVISIONS:-3}"

mkdir -p "$TASK_PENDING" "$TASK_DONE" "$LOG_ROOT"

required_prompts=(
  "researcher.md"
  "research-auditor.md"
  "research-reviser.md"
  "engineer.md"
  "engineering-auditor.md"
  "engineering-reviser.md"
  "technical-writer.md"
)

for prompt in "${required_prompts[@]}"; do
  [[ -f "$PROMPTS/$prompt" ]] || {
    echo "ERROR: prompt tidak ditemukan: $PROMPTS/$prompt"
    exit 1
  }
done

[[ -d "$REPO_ROOT" ]] || {
  echo "ERROR: repository tidak ditemukan: $REPO_ROOT"
  exit 1
}

cd "$REPO_ROOT"

git rev-parse --show-toplevel >/dev/null 2>&1 || {
  echo "ERROR: $REPO_ROOT bukan Git repository."
  exit 1
}

REPO_ROOT="$(git rev-parse --show-toplevel)"
cd "$REPO_ROOT"

timestamp() {
  date +"%Y-%m-%d %H:%M:%S"
}

task_path() {
  local task_file="$1"

  if [[ -f "$TASK_PENDING/$task_file" ]]; then
    echo "$TASK_PENDING/$task_file"
    return 0
  fi

  if [[ -f "$TASK_ROOT/$task_file" ]]; then
    echo "$TASK_ROOT/$task_file"
    return 0
  fi

  echo ""
}

resolve_lab_from_task() {
  local task="$1"
  local target=""

  # Prefer explicit Target Lab from markdown.
  target="$(
    grep -Eo 'labs/[0-9]+-[A-Za-z0-9._-]+' "$task" \
      | head -n 1 \
      || true
  )"

  if [[ -n "$target" ]]; then
    echo "$target"
    return 0
  fi

  # Fallback from filename: lab-15-topic.md -> labs/15-*
  local base number matches
  base="$(basename "$task")"

  if [[ "$base" =~ ^lab-([0-9]+)-topic\.md$ ]]; then
    number="${BASH_REMATCH[1]}"
    matches=( "$REPO_ROOT"/labs/"$number"-* )

    if [[ -e "${matches[0]}" ]]; then
      local rel="${matches[0]#"$REPO_ROOT"/}"
      echo "$rel"
      return 0
    fi
  fi

  echo ""
}

archive_task_done() {
  local task_file="$1"
  local source
  source="$(task_path "$task_file")"

  [[ -n "$source" && -f "$source" ]] || {
    echo "WARNING: task tidak ditemukan untuk dipindahkan ke done: $task_file"
    return 0
  }

  local target="$TASK_DONE/$task_file"

  if [[ -e "$target" ]]; then
    local stamp
    stamp="$(date +"%Y%m%d-%H%M%S")"
    target="$TASK_DONE/${task_file%.md}-$stamp.md"
  fi

  mv "$source" "$target"
  echo "TASK DONE: $target"
}

log_file_for() {
  local lab="$1"
  local stage="$2"
  local slug
  slug="$(basename "$lab")"

  mkdir -p "$LOG_ROOT/$slug"
  echo "$LOG_ROOT/$slug/${stage}.log"
}

run_opencode() {
  local lab="$1"
  local stage="$2"
  local prompt_file="$3"
  local instruction="$4"
  local logfile
  logfile="$(log_file_for "$lab" "$stage")"

  echo
  echo "[$(timestamp)] START: $stage"
  echo "LAB: $lab"
  echo "Log: $logfile"

  set +e
  (
    cd "$REPO_ROOT"
    opencode run "$(cat "$prompt_file")

---

$instruction"
  ) 2>&1 | tee "$logfile"

  local rc=${PIPESTATUS[0]}
  set -e

  if [[ $rc -ne 0 ]]; then
    echo "ERROR: stage $stage gagal dengan exit code $rc"
    return "$rc"
  fi

  echo "[$(timestamp)] DONE: $stage"
}

parse_verdict() {
  local file="$1"

  [[ -f "$file" ]] || {
    echo "MISSING"
    return
  }

  local result
  result="$(
    grep -Eo 'APPROVED_WITH_WARNINGS|NEEDS_REVISION|REJECTED|APPROVED' "$file" \
      | tail -n 1 \
      || true
  )"

  [[ -n "$result" ]] && echo "$result" || echo "UNKNOWN"
}

run_research() {
  local lab="$1"
  local task_file="$2"
  local task
  task="$(task_path "$task_file")"

  [[ -n "$task" && -f "$task" ]] || {
    echo "ERROR: task tidak ditemukan: $task_file"
    return 1
  }

  mkdir -p "$lab"

  run_opencode \
    "$lab" \
    "01-research" \
    "$PROMPTS/researcher.md" \
    "Target lab: $lab

Execute the topic specification below for the RESEARCH PHASE ONLY.

PIPELINE BOUNDARY:
- Research only.
- Do not implement source code.
- Do not write tests.
- Do not generate publication content.
- Engineering is handled by a separate Engineer Agent.
- Write research output only inside: $lab/research/

TOPIC SPECIFICATION:

$(cat "$task")"
}

run_research_audit() {
  local lab="$1"
  local round="$2"

  rm -rf "$lab/research-audit"

  run_opencode \
    "$lab" \
    "02-research-audit-r${round}" \
    "$PROMPTS/research-auditor.md" \
    "Audit the research for target lab:

$lab

PIPELINE OVERRIDE:
- Audit research only.
- Do not audit implementation/code in this stage.
- Do not modify research files.
- Write all audit output to: $lab/research-audit/
- Final verdict must be written to:
  $lab/research-audit/07-verdict.md"
}

run_research_revision() {
  local lab="$1"
  local round="$2"

  rm -rf "$lab/research-revision"

  run_opencode \
    "$lab" \
    "03-research-revision-r${round}" \
    "$PROMPTS/research-reviser.md" \
    "Revise the research for target lab:

$lab

Use findings from:
$lab/research-audit/

PIPELINE OVERRIDE:
- Fix research issues only.
- Do not implement code.
- Do not generate publication content.
- Modify research files only where required by the audit.
- Write revision records to:
  $lab/research-revision/
- Finish with READY_FOR_REAUDIT."
}

research_pipeline() {
  local lab="$1"
  local task_file="$2"

  run_research "$lab" "$task_file"

  local revision_count=0

  while true; do
    local audit_round=$((revision_count + 1))
    run_research_audit "$lab" "$audit_round"

    local verdict
    verdict="$(parse_verdict "$lab/research-audit/07-verdict.md")"

    echo "Research verdict: $verdict"

    case "$verdict" in
      APPROVED|APPROVED_WITH_WARNINGS)
        return 0
        ;;
      NEEDS_REVISION)
        if (( revision_count >= MAX_RESEARCH_REVISIONS )); then
          echo "BLOCKED: maximum research revisions tercapai."
          return 1
        fi
        revision_count=$((revision_count + 1))
        run_research_revision "$lab" "$revision_count"
        ;;
      REJECTED)
        echo "BLOCKED: research ditolak auditor."
        return 1
        ;;
      *)
        echo "BLOCKED: research verdict tidak valid: $verdict"
        return 1
        ;;
    esac
  done
}

run_engineer() {
  local lab="$1"

  run_opencode \
    "$lab" \
    "04-engineering" \
    "$PROMPTS/engineer.md" \
    "Implement the approved technical lab:

$lab

Use ONLY the approved research in:
$lab/research/

Use the latest approved research audit in:
$lab/research-audit/

Requirements:
- Design the runnable lab.
- Write actual implementation code.
- Write automated tests.
- Execute tests.
- Run race detector when relevant.
- Run the demo.
- Fix implementation failures.
- Record real execution results.
- Do not generate publication content."
}

run_engineering_audit() {
  local lab="$1"
  local round="$2"

  rm -rf "$lab/engineering-audit"

  run_opencode \
    "$lab" \
    "05-engineering-audit-r${round}" \
    "$PROMPTS/engineering-auditor.md" \
    "Audit the implementation for target lab:

$lab

Use approved research as the specification.

Requirements:
- Inspect code and tests.
- Execute required build/test/demo commands yourself.
- Verify README against actual behavior.
- Do not modify implementation.
- Write all audit output to:
  $lab/engineering-audit/
- Final verdict:
  $lab/engineering-audit/06-verdict.md"
}

run_engineering_revision() {
  local lab="$1"
  local round="$2"

  rm -rf "$lab/engineering-revision"

  run_opencode \
    "$lab" \
    "06-engineering-revision-r${round}" \
    "$PROMPTS/engineering-reviser.md" \
    "Revise the implementation for target lab:

$lab

Use engineering audit findings in:
$lab/engineering-audit/

Requirements:
- Fix only confirmed engineering issues.
- Add/fix tests where required.
- Re-run build/tests/race detector/demo.
- Do not perform new research.
- Do not generate publication content.
- Write revision records to:
  $lab/engineering-revision/
- Finish with READY_FOR_ENGINEERING_REAUDIT."
}

engineering_pipeline() {
  local lab="$1"

  run_engineer "$lab"

  local revision_count=0

  while true; do
    local audit_round=$((revision_count + 1))
    run_engineering_audit "$lab" "$audit_round"

    local verdict
    verdict="$(parse_verdict "$lab/engineering-audit/06-verdict.md")"

    echo "Engineering verdict: $verdict"

    case "$verdict" in
      APPROVED|APPROVED_WITH_WARNINGS)
        return 0
        ;;
      NEEDS_REVISION)
        if (( revision_count >= MAX_ENGINEERING_REVISIONS )); then
          echo "BLOCKED: maximum engineering revisions tercapai."
          return 1
        fi
        revision_count=$((revision_count + 1))
        run_engineering_revision "$lab" "$revision_count"
        ;;
      REJECTED)
        echo "BLOCKED: engineering ditolak auditor."
        return 1
        ;;
      *)
        echo "BLOCKED: engineering verdict tidak valid: $verdict"
        return 1
        ;;
    esac
  done
}

run_writer() {
  local lab="$1"

  rm -rf "$lab/content"

  run_opencode \
    "$lab" \
    "07-technical-writer" \
    "$PROMPTS/technical-writer.md" \
    "Create the canonical technical master content for:

$lab

Use ONLY:
- approved research
- approved research audit
- approved implementation
- approved engineering audit
- verified tests/demo
- revision records where applicable

Do not:
- perform new research
- modify code
- add external facts
- create platform-specific content

Write output to:
$lab/content/"
}

process_task() {
  local task_file="$1"
  local task
  task="$(task_path "$task_file")"

  [[ -n "$task" ]] || {
    echo "ERROR: task tidak ditemukan: $task_file"
    return 1
  }

  local lab
  lab="$(resolve_lab_from_task "$task")"

  if [[ -z "$lab" ]]; then
    echo "ERROR: tidak bisa menentukan target lab dari task: $task"
    echo "Tambahkan explicit path seperti:"
    echo '  `labs/15-load-testing`'
    return 1
  fi

  echo
  echo "============================================================"
  echo "TASK: $task_file"
  echo "LAB : $lab"
  echo "============================================================"

  research_pipeline "$lab" "$task_file"
  engineering_pipeline "$lab"
  run_writer "$lab"
  archive_task_done "$task_file"

  echo
  echo "COMPLETED: $task_file -> $lab"
}

run_all_pending() {
  shopt -s nullglob
  local tasks=( "$TASK_PENDING"/lab-*-topic.md )
  shopt -u nullglob

  if (( ${#tasks[@]} == 0 )); then
    echo "Tidak ada task pending."
    echo "Folder: $TASK_PENDING"
    return 0
  fi

  echo "Pending tasks: ${#tasks[@]}"
  printf ' - %s\n' "${tasks[@]##*/}"

  for task in "${tasks[@]}"; do
    process_task "$(basename "$task")"
  done
}

usage() {
  cat <<'EOF'
Usage:

  ./scripts/run-until-writer.sh all
  ./scripts/run-until-writer.sh lab-15-topic.md
  ./scripts/run-until-writer.sh lab-16-topic.md

"all" hanya menjalankan file yang ada di:

  ~/research-agent/tasks/pending/

Task yang selesai otomatis dipindahkan ke:

  ~/research-agent/tasks/done/
EOF
}

case "${1:-all}" in
  all)
    run_all_pending
    ;;
  lab-*-topic.md)
    process_task "$1"
    ;;
  *)
    usage
    exit 1
    ;;
esac
