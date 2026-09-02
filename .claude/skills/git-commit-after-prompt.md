# git-commit-after-prompt

Skill to automatically commit changes after prompts.

## What it does

This skill provides git commit capability that runs after Claude Code makes file changes. It can be triggered automatically via hooks or manually.

## Usage

### Automatic Mode (Hook-based)

The skill comes with a configured `PostToolUse` hook that triggers after `Write` or `Edit` operations. To enable:

1. The hook is already configured in `.claude/settings.local.json`
2. It checks for uncommitted changes after edits
3. Auto-generate commit message: `git add <file> && git commit -m "Commit Claude Code: <filename>; <HH:MM:SS>"`

To customize the commit message, edit the hint in the hook configuration.

### Manual Usage

Run `/git-commit-after-prompt` slash command to manually trigger a commit after changes.

## Configuration Options

To modify behavior, edit `.claude/settings.local.json`:

```json
{
  "hooks": {
    "PostToolUse": [{
      "matcher": "Write|Edit",
      "hooks": [{
        "type": "agent",
        "prompt": "Custom prompt here",
        "timeout": 30
      }]
    }]
  }
}
```

## When to Use

- **Automatic**: For projects where you want commits after every file edit
- **Manual**: For selective commits when you need control over timing

## Notes

- Requires git repository initialized
- Will only commit if there are actual changes (staged or unstaged)
- Commit message can be customized in the hook configuration