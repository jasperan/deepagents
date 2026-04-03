---
name: oracle-multi-agent-memory
description: Enable multiple agents to share memory through Oracle namespaces. Use when agents need to collaborate, share findings, or publish results for other agents to consume.
allowed-tools: read_file, write_file, ls, glob, grep
---

# Oracle Multi-Agent Memory

Share memory between multiple agents using Oracle namespace conventions.

## When to Use

- User wants agents to collaborate on a task
- User says "share this with other agents" or "publish findings"
- User sets up a multi-agent workflow
- User wants one agent to read another agent's output

## Concepts

### Namespaces

Each agent has its own namespace (set at creation time). Namespaces isolate data so agents don't accidentally overwrite each other's files.

### Shared paths

By convention, paths under `/memory/shared/` are meant to be read by all agents. Each agent writes to its own sub-directory:

```
/memory/shared/<agent-name>/<topic>.md
```

### Private paths

Everything else under `/memory/` is private to the agent's namespace.

## Workflow

### Publishing findings (writer agent)

When you have results other agents should see:

```
write_file(
    path="/memory/shared/<your-name>/<topic>.md",
    content="# <Topic>\n\nAuthor: <your name>\nDate: <today>\nNamespace: <your namespace>\n\n## Findings\n<content>"
)
```

### Reading shared findings (reader agent)

To see what other agents have published:

```
ls(path="/memory/shared/")
```

Then read specific findings:

```
read_file(path="/memory/shared/<other-agent>/<topic>.md")
```

### Collaboration pattern

1. **Research agent** writes findings to `/memory/shared/researcher/`
2. **Analyst agent** reads researcher's findings, writes analysis to `/memory/shared/analyst/`
3. **Writer agent** reads both, produces final output at `/memory/shared/writer/`

### Setting up shared memory

When creating multiple agents that need to share memory, use the same Oracle backend but different namespaces:

```python
from deepagents_oracle import create_oracle_deep_agent

# Both agents share the same Oracle backend but have separate namespaces
researcher = create_oracle_deep_agent(namespace="researcher")
analyst = create_oracle_deep_agent(namespace="analyst")
```

To read across namespaces, agents use the shared path convention above. The shared directory exists in each agent's namespace, but the convention makes intent clear.

## Tips

- Always include author and date metadata in shared files
- Use descriptive topic names so other agents can find content
- List `/memory/shared/` before reading to see what's available
- Keep shared files self-contained (don't reference private paths)
