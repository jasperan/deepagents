---
name: oracle-memory-search
description: Semantic search over agent memory stored in Oracle AI Database. Use when the user asks to find, search, or recall previously saved notes, findings, or knowledge.
allowed-tools: read_file, grep, glob, ls
---

# Oracle Memory Search

Search your persistent memory in Oracle using semantic similarity or keyword matching.

## When to Use

- User asks "find my notes about X"
- User asks "what did I save about Y"
- User asks "search memory for Z"
- User wants to recall something from a previous session

## Workflow

### Step 1: Determine search type

- If the user's query is conceptual (e.g., "anything about machine learning"), use **semantic search** (Step 2a)
- If the user's query is a specific string (e.g., "the function called parse_config"), use **keyword search** (Step 2b)

### Step 2a: Semantic search

Use the `grep` tool with the user's query against the `/memory/` directory:

```
grep(query="<user's search terms>", path="/memory/")
```

This searches file contents. Review the matches and read the most relevant files.

### Step 2b: Keyword search

Use the `grep` tool for exact text matching:

```
grep(query="<exact keyword>", path="/memory/")
```

### Step 3: Read full content

For each relevant match, read the full file:

```
read_file(path="/memory/<matched-file>")
```

### Step 4: Synthesize results

- Summarize what was found across all matching files
- Quote specific passages that are most relevant
- If nothing was found, tell the user and suggest they save information first

## Tips

- Search `/memory/research/` for research findings
- Search `/memory/kb/` for ingested knowledge base content
- Search `/memory/insights/` for data analysis results
- Use `ls(path="/memory/")` first to see what's available if the search is broad
