---
name: oracle-data-analyst
description: Query Oracle database tables directly and store insights. Use when the user wants to run SQL queries, analyze data, explore tables, or generate reports from Oracle.
allowed-tools: read_file, write_file, execute, ls
---

# Oracle Data Analyst

Query Oracle tables directly, analyze results, and store insights in persistent memory.

## When to Use

- User asks "query the database" or "run this SQL"
- User asks "what tables exist" or "describe table X"
- User wants to analyze data stored in Oracle
- User asks for a report or summary from database data

## Workflow

### Step 1: Understand the request

Clarify what the user wants to know. If they haven't specified a table, start with discovery.

### Step 2: Discover available tables (if needed)

Write and execute a discovery script:

```python
# Save to a temp file and execute
import oracledb
import os

conn = oracledb.connect(
    user=os.environ.get("DEEPAGENTS_ORACLE_USER", "deepagents"),
    password=os.environ.get("DEEPAGENTS_ORACLE_PASSWORD", "DeepAgents2024"),
    dsn=os.environ.get("DEEPAGENTS_ORACLE_DSN",
        f"{os.environ.get('DEEPAGENTS_ORACLE_HOST', 'localhost')}:"
        f"{os.environ.get('DEEPAGENTS_ORACLE_PORT', '1521')}/"
        f"{os.environ.get('DEEPAGENTS_ORACLE_SERVICE', 'FREEPDB1')}")
)

cur = conn.cursor()
cur.execute("SELECT table_name FROM user_tables ORDER BY table_name")
for row in cur:
    print(row[0])
conn.close()
```

### Step 3: Write and execute the query

Write a Python script that:
1. Connects to Oracle using env vars
2. Executes the user's SQL (SELECT only by default)
3. Formats and prints results

```python
import oracledb
import os

conn = oracledb.connect(
    user=os.environ.get("DEEPAGENTS_ORACLE_USER", "deepagents"),
    password=os.environ.get("DEEPAGENTS_ORACLE_PASSWORD", "DeepAgents2024"),
    dsn=os.environ.get("DEEPAGENTS_ORACLE_DSN",
        f"{os.environ.get('DEEPAGENTS_ORACLE_HOST', 'localhost')}:"
        f"{os.environ.get('DEEPAGENTS_ORACLE_PORT', '1521')}/"
        f"{os.environ.get('DEEPAGENTS_ORACLE_SERVICE', 'FREEPDB1')}")
)

cur = conn.cursor()
cur.execute("<USER SQL HERE>")

# Print column headers
cols = [d[0] for d in cur.description]
print(" | ".join(cols))
print("-" * 60)

# Print rows
for row in cur:
    print(" | ".join(str(v) for v in row))

conn.close()
```

### Step 4: Analyze results

- Summarize key findings
- Identify patterns, outliers, or notable values
- Answer the user's original question

### Step 5: Save insights (optional)

If the analysis produced valuable findings:

```
write_file(
    path="/memory/insights/<topic>.md",
    content="# Insight: <topic>\n\nDate: <today>\nQuery: <the SQL>\n\n## Findings\n<analysis>\n\n## Raw Data\n<formatted results>"
)
```

## Safety Rules

1. **SELECT only by default.** Never run INSERT, UPDATE, DELETE, DROP, TRUNCATE, or ALTER without explicit user confirmation.
2. **Always show the SQL** to the user before executing.
3. **Limit results** with `FETCH FIRST 100 ROWS ONLY` for large tables.
4. **Never expose passwords** in output or saved files.
