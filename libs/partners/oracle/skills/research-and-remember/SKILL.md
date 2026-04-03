---
name: research-and-remember
description: Multi-turn research loop that persists findings to Oracle. Use when the user wants to research a topic thoroughly, investigate something, or build up knowledge over multiple sessions.
allowed-tools: read_file, write_file, ls, glob, grep
---

# Research and Remember

Conduct structured research on a topic and persist all findings to Oracle for future sessions.

## When to Use

- User says "research X and save findings"
- User says "deep dive into Y"
- User says "investigate Z and remember what you find"
- Any open-ended research task that should persist

## Workflow

### Step 1: Check for existing research

Before starting, check if you've already researched this topic:

```
ls(path="/memory/research/")
grep(query="<topic>", path="/memory/research/")
```

If prior research exists, read it first and build on it rather than starting over.

### Step 2: Break the topic into sub-questions

Decompose the research topic into 3-7 specific sub-questions. Write the research plan:

```
write_file(
    path="/memory/research/<topic>/plan.md",
    content="# Research Plan: <topic>\n\nDate started: <today>\n\n## Sub-questions\n1. <question 1>\n2. <question 2>\n3. <question 3>\n\n## Status\n- [ ] Question 1\n- [ ] Question 2\n- [ ] Question 3"
)
```

### Step 3: Research each sub-question

For each sub-question, use available tools to find answers:
- Read relevant files in the codebase
- Search memory for related prior findings
- Use any available search or analysis tools

### Step 4: Save each finding

Write individual finding files:

```
write_file(
    path="/memory/research/<topic>/finding-01.md",
    content="# Finding: <sub-question answered>\n\nDate: <today>\nConfidence: high/medium/low\nSource: <where you found this>\n\n## Answer\n<detailed finding>\n\n## Evidence\n<supporting details>"
)
```

Number findings sequentially. Include confidence level and source attribution.

### Step 5: Update the plan

After each finding, update the research plan to mark questions as complete:

```
edit_file(path="/memory/research/<topic>/plan.md", ...)
```

### Step 6: Write a summary

When all sub-questions are answered (or the user is satisfied):

```
write_file(
    path="/memory/research/<topic>/summary.md",
    content="# Research Summary: <topic>\n\nDate: <today>\nFindings: N\n\n## Key Takeaways\n1. <takeaway 1>\n2. <takeaway 2>\n3. <takeaway 3>\n\n## Detailed Findings\n- finding-01.md: <one-line summary>\n- finding-02.md: <one-line summary>\n\n## Open Questions\n<anything still unknown>"
)
```

### Step 7: Present results

Give the user a concise summary of what you found, noting:
- Key takeaways
- Confidence levels
- Open questions that need more investigation
- Where everything is stored for future reference

## Tips

- Always check for existing research first to avoid duplicating work
- Use consistent naming: `/memory/research/<topic>/` as the base directory
- Include confidence levels (high/medium/low) in every finding
- The summary file should be readable standalone without opening individual findings
- Future sessions can pick up where you left off by reading the plan and checking status
