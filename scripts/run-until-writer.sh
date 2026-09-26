#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="${REPO_ROOT:-$HOME/Documents/LUKMAN/software-engineering-lab}"
AGENT_ROOT="${AGENT_ROOT:-$HOME/research-agent}"

PROMPTS="$AGENT_ROOT/prompts"
TASK_ROOT="$AGENT_ROOT/tasks"
TASK_PENDING="$TASK_ROOT/pending"
TASK_DONE="$TASK_ROOT/done"
LOG_ROOT="$AGENT_ROOT/logs"

MAX_RESEARCH_REVISIONS="${MAX_RESEARCH_REVISIONS:-6}"
MAX_ENGINEERING_REVISIONS="${MAX_ENGINEERING_REVISIONS:-6}"
MAX_WRITER_REVISIONS="${MAX_WRITER_REVISIONS:-6}"

MODEL_DEFAULT="${MODEL_DEFAULT:-9router/free-combo}"
MODEL_CRITICAL="${MODEL_CRITICAL:-9router/ag-combo}"
MODEL_FALLBACK="${MODEL_FALLBACK:-9router/ag-combo}"
MODEL_AUDITOR_OS="${MODEL_AUDITOR_OS:-9router/auditor-opensource-combo}"

mkdir -p "$TASK_PENDING" "$TASK_DONE" "$LOG_ROOT"

required_prompts=(
  "researcher.md"
  "research-auditor.md"
  "research-reviser.md"
  "engineer.md"
  "engineering-auditor.md"
  "engineering-reviser.md"
  "technical-writer.md"
  "technical-writer-auditor.md"
  "technical-writer-reviser.md"
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
  local model="${5:-$MODEL_DEFAULT}"
  local logfile
  logfile="$(log_file_for "$lab" "$stage")"

  local cmd_args=()
  if [[ -n "$model" ]]; then
    cmd_args+=("-m" "$model")
  fi

  echo
  echo "[$(timestamp)] START: $stage"
  echo "LAB: $lab"
  echo "Log: $logfile"
  echo "Model: $model"

  set +e
  (
    cd "$REPO_ROOT"
    opencode run "${cmd_args[@]}" "$(cat "$prompt_file")

---

$instruction"
  ) 2>&1 | tee "$logfile"

  local rc=${PIPESTATUS[0]}

  if [[ $rc -ne 0 && "$model" != "$MODEL_FALLBACK" ]]; then
    echo "WARNING: Model $model gagal (exit $rc). Fallback ke $MODEL_FALLBACK..."
    
    cmd_args=("-m" "$MODEL_FALLBACK")
    (
      cd "$REPO_ROOT"
      opencode run "${cmd_args[@]}" "$(cat "$prompt_file")

---

$instruction"
    ) 2>&1 | tee -a "$logfile"
    rc=${PIPESTATUS[0]}
  fi

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
  $lab/research-audit/07-verdict.md" "$MODEL_CRITICAL"
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
      APPROVED)
        return 0
        ;;
      NEEDS_REVISION|APPROVED_WITH_WARNINGS)
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
- Do not generate publication content." "$MODEL_CRITICAL"
}

run_engineering_audit() {
  local lab="$1"
  local round="$2"
  local dir_suffix="${3:-}"
  local model="${4:-}"
  local audit_dir="$lab/engineering-audit${dir_suffix}"

  rm -rf "$audit_dir"

  run_opencode \
    "$lab" \
    "05-engineering-audit${dir_suffix}-r${round}" \
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
  $audit_dir/
- Final verdict:
  $audit_dir/06-verdict.md" "$model"
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
and
$lab/engineering-audit-opensource/ (if available)

Requirements:
- Fix only confirmed engineering issues.
- Add/fix tests where required.
- Re-run build/tests/race detector/demo.
- Do not perform new research.
- Do not generate publication content.
- Write revision records to:
  $lab/engineering-revision/
- Finish with READY_FOR_ENGINEERING_REAUDIT." "$MODEL_CRITICAL"
}

engineering_pipeline() {
  local lab="$1"

  run_engineer "$lab"

  local revision_count=0

  while true; do
    local audit_round=$((revision_count + 1))
    
    run_engineering_audit "$lab" "$audit_round" "" "$MODEL_CRITICAL" &
    local pid1=$!
    run_engineering_audit "$lab" "$audit_round" "-opensource" "$MODEL_AUDITOR_OS" &
    local pid2=$!

    wait $pid1
    wait $pid2

    local verdict1
    verdict1="$(parse_verdict "$lab/engineering-audit/06-verdict.md")"
    
    local verdict2
    verdict2="$(parse_verdict "$lab/engineering-audit-opensource/06-verdict.md")"

    echo "Engineering verdict 1 (default): $verdict1"
    echo "Engineering verdict 2 (opensource): $verdict2"

    local final_verdict="NEEDS_REVISION"
    if [[ "$verdict1" == "APPROVED" && "$verdict2" == "APPROVED" ]]; then
       final_verdict="APPROVED"
    elif [[ "$verdict1" == "REJECTED" || "$verdict2" == "REJECTED" ]]; then
       final_verdict="REJECTED"
    elif [[ "$verdict1" == "APPROVED_WITH_WARNINGS" || "$verdict2" == "APPROVED_WITH_WARNINGS" ]]; then
       final_verdict="APPROVED_WITH_WARNINGS"
    fi

    case "$final_verdict" in
      APPROVED)
        return 0
        ;;
      NEEDS_REVISION|APPROVED_WITH_WARNINGS)
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
        echo "BLOCKED: engineering final verdict tidak valid."
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

run_writer_audit() {
  local lab="$1"
  local round="$2"

  rm -rf "$lab/content-audit"

  run_opencode \
    "$lab" \
    "08-content-audit-r${round}" \
    "$PROMPTS/technical-writer-auditor.md" \
    "Audit the technical content for target lab:

$lab

Use approved research and implementation as the specification.

Requirements:
- Inspect the generated content in: $lab/content/
- Verify accuracy, completeness, and formatting against the research and engineering outputs.
- Do not modify the content files directly.
- Write all audit output to:
  $lab/content-audit/
- Final verdict must be written to:
  $lab/content-audit/09-verdict.md"
}

run_writer_revision() {
  local lab="$1"
  local round="$2"

  rm -rf "$lab/content-revision"

  run_opencode \
    "$lab" \
    "09-content-revision-r${round}" \
    "$PROMPTS/technical-writer-reviser.md" \
    "Revise the technical content for target lab:

$lab

Use content audit findings in:
$lab/content-audit/

Requirements:
- Fix only confirmed content issues.
- Do not perform new research or modify code.
- Write revision records to:
  $lab/content-revision/
- Finish with READY_FOR_CONTENT_REAUDIT."
}

writer_pipeline() {
  local lab="$1"

  run_writer "$lab"

  local revision_count=0

  while true; do
    local audit_round=$((revision_count + 1))
    run_writer_audit "$lab" "$audit_round"

    local verdict
    verdict="$(parse_verdict "$lab/content-audit/09-verdict.md")"

    echo "Writer verdict: $verdict"

    case "$verdict" in
      APPROVED)
        return 0
        ;;
      NEEDS_REVISION|APPROVED_WITH_WARNINGS)
        if (( revision_count >= MAX_WRITER_REVISIONS )); then
          echo "BLOCKED: maximum writer revisions tercapai."
          return 1
        fi
        revision_count=$((revision_count + 1))
        run_writer_revision "$lab" "$revision_count"
        ;;
      REJECTED)
        echo "BLOCKED: content ditolak auditor."
        return 1
        ;;
      *)
        echo "BLOCKED: content verdict tidak valid: $verdict"
        return 1
        ;;
    esac
  done
}

process_task() {
  local task_file="$1"

  # Atomic lock untuk multi-instance
  if ! mkdir "$TASK_PENDING/.lock_${task_file}" 2>/dev/null; then
    echo "SKIP: $task_file sedang dikerjakan instance lain (terkunci)."
    return 0
  fi

  local task
  task="$(task_path "$task_file")"

  [[ -n "$task" ]] || {
    echo "ERROR: task tidak ditemukan: $task_file"
    rm -rf "$TASK_PENDING/.lock_${task_file}"
    return 1
  }

  local lab
  lab="$(resolve_lab_from_task "$task")"

  if [[ -z "$lab" ]]; then
    echo "ERROR: tidak bisa menentukan target lab dari task: $task"
    echo "Tambahkan explicit path seperti:"
    echo '  `labs/15-load-testing`'
    rm -rf "$TASK_PENDING/.lock_${task_file}"
    return 1
  fi

  echo
  echo "============================================================"
  echo "TASK: $task_file"
  echo "LAB : $lab"
  echo "============================================================"

  research_pipeline "$lab" "$task_file"
  engineering_pipeline "$lab"
  writer_pipeline "$lab"
  archive_task_done "$task_file"

  rm -rf "$TASK_PENDING/.lock_${task_file}"

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
    process_task "$(basename "$task")" &
  done

  echo "Menunggu semua task (async) selesai..."
  wait
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
