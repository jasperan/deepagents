---
name: oracle-knowledge-ingest
description: Ingest documents into Oracle as a persistent, searchable knowledge base. Use when the user wants to save content for later retrieval, build a knowledge base, or remember important information.
allowed-tools: read_file, write_file, ls, glob
---

# Oracle Knowledge Ingest

Build a persistent, searchable knowledge base in Oracle by ingesting documents, notes, and content.

## When to Use

- User says "ingest this document" or "add to knowledge base"
- User says "remember this" or "save this for later"
- User provides a large block of text to store
- User wants to import content from a file

## Workflow

### Step 1: Identify the content source

- **Pasted text**: User provides content directly in the chat
- **File path**: User references a file to ingest (read it first)
- **Multiple files**: User wants to ingest a directory of files

### Step 2: Determine the topic

Ask the user or infer a topic name from the content. This becomes the directory name:
`/memory/kb/<topic>/`

### Step 3: Chunk the content

Break content into logical chunks for better searchability:

- **By heading**: Split at `#`, `##`, `###` headings (preferred for markdown)
- **By paragraph**: Split at double newlines (for plain text)
- **By size**: Split at ~500 words if no natural boundaries exist

Each chunk should be self-contained and meaningful on its own.

### Step 4: Write chunks to Oracle

For each chunk, write a file:

```
write_file(
    path="/memory/kb/<topic>/chunk-01.md",
    content="---\nsource: <original source>\ndate: <today>\ntopic: <topic>\nchunk: 1 of N\n---\n\n<chunk content>"
)
```

Number chunks sequentially: `chunk-01.md`, `chunk-02.md`, etc.

### Step 5: Write an index file

Create a summary index:

```
write_file(
    path="/memory/kb/<topic>/index.md",
    content="# Knowledge Base: <topic>\n\nSource: <source>\nIngested: <date>\nChunks: N\n\n## Contents\n- chunk-01.md: <first line summary>\n- chunk-02.md: <first line summary>\n..."
)
```

### Step 6: Confirm ingestion

Tell the user:
- How many chunks were created
- Where they're stored (`/memory/kb/<topic>/`)
- How to search them later ("ask me to search memory for <topic>")

## Tips

- Keep chunk sizes between 200-800 words for optimal search relevance
- Include metadata headers in every chunk for context
- The index file helps the agent find content without searching every chunk
- Embeddings are auto-generated on write if the vector backend is active
