---
name: oracle-session-resume
description: Resume a previous conversation by loading session history from Oracle. Use when the user wants to continue where they left off, recall previous work, or review past sessions.
allowed-tools: read_file, ls, glob
---

# Oracle Session Resume

Load and review previous session history stored in Oracle to continue past work.

## When to Use

- User says "resume session" or "continue where we left off"
- User asks "what were we working on last time"
- User wants to review past conversations or decisions

## Workflow

### Step 1: List available sessions

```
ls(path="/history/")
```

This shows all saved session files. Sessions are typically named with timestamps or IDs.

### Step 2: Read the most recent session

Read the latest session file (sorted by name/date):

```
read_file(path="/history/<most-recent-session>")
```

### Step 3: Summarize and offer to continue

Present the user with:
1. **When**: When the session occurred
2. **What**: Key topics discussed or tasks worked on
3. **Status**: What was completed vs. in progress
4. **Next steps**: What was planned but not yet done

Ask the user which thread they'd like to continue.

### Step 4: Load relevant memory

Based on the session context, proactively load related memory files:

```
ls(path="/memory/")
grep(query="<topic from session>", path="/memory/")
```

This gives you full context to pick up seamlessly.

## Saving Sessions

To save the current session for future resumption, write a summary at the end:

```
write_file(
    path="/history/<timestamp>-<topic>.md",
    content="# Session: <topic>\n\n## Date: <date>\n\n## Summary\n<what was discussed>\n\n## Decisions\n<decisions made>\n\n## Next Steps\n<what to do next>"
)
```

Always save sessions before ending if important work was discussed.
