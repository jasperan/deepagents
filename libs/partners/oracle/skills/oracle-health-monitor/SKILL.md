---
name: oracle-health-monitor
description: Agent self-diagnostics for the Oracle backend. Use when the user wants to check if Oracle is working, diagnose connection issues, or review backend health and statistics.
allowed-tools: read_file, write_file, execute, ls
---

# Oracle Health Monitor

Run diagnostics on the Oracle backend and report health status.

## When to Use

- User says "check health" or "is Oracle working"
- User says "diagnose connection" or "backend status"
- Something seems slow or broken
- User wants to see memory usage stats

## Workflow

### Step 1: Run health check script

Execute the following diagnostic script:

```python
import oracledb
import os
import json
from datetime import datetime

results = {"timestamp": datetime.now().isoformat(), "checks": []}

try:
    conn = oracledb.connect(
        user=os.environ.get("DEEPAGENTS_ORACLE_USER", "deepagents"),
        password=os.environ.get("DEEPAGENTS_ORACLE_PASSWORD", "DeepAgents2024"),
        dsn=os.environ.get("DEEPAGENTS_ORACLE_DSN",
            f"{os.environ.get('DEEPAGENTS_ORACLE_HOST', 'localhost')}:"
            f"{os.environ.get('DEEPAGENTS_ORACLE_PORT', '1521')}/"
            f"{os.environ.get('DEEPAGENTS_ORACLE_SERVICE', 'FREEPDB1')}")
    )
    results["checks"].append({"name": "connection", "status": "OK"})

    cur = conn.cursor()

    # Basic connectivity
    cur.execute("SELECT 1 FROM DUAL")
    results["checks"].append({"name": "query_execution", "status": "OK"})

    # Table existence
    cur.execute("SELECT COUNT(*) FROM user_tables WHERE table_name = 'DA_FILES'")
    table_exists = cur.fetchone()[0] > 0
    results["checks"].append({
        "name": "da_files_table",
        "status": "OK" if table_exists else "MISSING"
    })

    if table_exists:
        # Total files
        cur.execute("SELECT COUNT(*) FROM da_files")
        total = cur.fetchone()[0]
        results["checks"].append({"name": "total_files", "value": total})

        # Files per namespace
        cur.execute(
            "SELECT namespace, COUNT(*) as cnt FROM da_files "
            "GROUP BY namespace ORDER BY cnt DESC FETCH FIRST 10 ROWS ONLY"
        )
        namespaces = [{"namespace": r[0], "files": r[1]} for r in cur]
        results["checks"].append({"name": "namespaces", "value": namespaces})

        # Recent activity
        cur.execute(
            "SELECT file_path, modified_at FROM da_files "
            "ORDER BY modified_at DESC FETCH FIRST 5 ROWS ONLY"
        )
        recent = [{"path": r[0], "modified": str(r[1])} for r in cur]
        results["checks"].append({"name": "recent_files", "value": recent})

        # Data size estimate
        cur.execute(
            "SELECT SUM(DBMS_LOB.GETLENGTH(content)) FROM da_files"
        )
        size = cur.fetchone()[0] or 0
        results["checks"].append({
            "name": "total_content_bytes",
            "value": size,
            "human": f"{size / 1024:.1f} KB" if size < 1048576 else f"{size / 1048576:.1f} MB"
        })

    # Oracle version
    cur.execute("SELECT banner FROM v$version FETCH FIRST 1 ROWS ONLY")
    row = cur.fetchone()
    if row:
        results["checks"].append({"name": "oracle_version", "value": row[0]})

    conn.close()
    results["overall"] = "HEALTHY"

except Exception as e:
    results["overall"] = "UNHEALTHY"
    results["error"] = str(e)

print(json.dumps(results, indent=2, default=str))
```

### Step 2: Interpret results

Read the JSON output and present a clear health report:

- **HEALTHY**: All checks passed. Report stats.
- **UNHEALTHY**: Something failed. Diagnose the specific issue.

### Step 3: Report to user

Format the results as a readable report:

```
## Oracle Health Report

Status: HEALTHY / UNHEALTHY
Timestamp: <when>

### Connectivity
- Connection: OK / FAILED
- Query execution: OK / FAILED

### Data
- Total files: N
- Content size: X KB/MB
- Namespaces: list with file counts

### Recent Activity
- <file1> (modified <when>)
- <file2> (modified <when>)

### Oracle Version
<version string>
```

### Step 4: Save report (optional)

If the user wants to track health over time:

```
write_file(
    path="/memory/health/<date>-health.md",
    content="<formatted report>"
)
```

## Common Issues

| Symptom | Likely Cause | Fix |
|---------|-------------|-----|
| Connection refused | Oracle container not running | `docker start oracle-free` |
| ORA-01017 | Wrong credentials | Check DEEPAGENTS_ORACLE_PASSWORD |
| DA_FILES missing | Schema not initialized | Run `setup-oracle.sh` |
| Slow queries | Too many files | Check total file count |
