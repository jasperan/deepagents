<div align="center">
  <a href="https://docs.langchain.com/oss/python/deepagents/overview#deep-agents-overview">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset=".github/images/logo-dark.svg">
      <source media="(prefers-color-scheme: light)" srcset=".github/images/logo-light.svg">
      <img alt="Deep Agents Logo" src=".github/images/logo-dark.svg" width="50%">
    </picture>
  </a>
</div>

<div align="center">
  <h3>The batteries-included agent harness + Oracle AI Database</h3>
</div>

<div align="center">
  <a href="https://opensource.org/licenses/MIT" target="_blank"><img src="https://img.shields.io/pypi/l/deepagents" alt="License"></a>
  <a href="https://pypi.org/project/deepagents-oracle/" target="_blank"><img src="https://img.shields.io/pypi/v/deepagents-oracle?label=deepagents-oracle" alt="PyPI"></a>
</div>

<br>

Fork of [LangChain Deep Agents](https://github.com/langchain-ai/deepagents) with **Oracle AI Database** as the persistent storage and memory layer.

## What Oracle Adds

- **ACID-guaranteed agent memory** that survives crashes, restarts, and redeployments
- **Cross-session conversation history** stored in Oracle, not lost when the process dies
- **In-database vector embeddings** for semantic search over agent memory (via Oracle AI Vector Search)
- **Zero external API calls** for storage. Everything lives in the database

## Two Deployment Modes

**Local development:** Oracle Database 26ai Free in a Docker container. One command to start.

**Production:** Oracle Autonomous Database in OCI. Set 3 env vars and you're connected.

## Quickstart

```bash
pip install deepagents deepagents-oracle
```

Start Oracle locally:
```bash
docker compose -f oracle/docker-compose.yml --profile freepdb up -d
bash oracle/scripts/setup-oracle.sh
```

Run an agent with persistent memory:
```python
from deepagents_oracle import create_oracle_deep_agent

agent = create_oracle_deep_agent()
result = agent.invoke(
    {"messages": [{"role": "user", "content": "Research LangGraph and save findings to /memory/langgraph.md"}]}
)
```

Memory, conversation history, and skills persist in Oracle. Restart the agent and everything's still there.

## How It Works

`create_oracle_deep_agent()` wires up a `CompositeBackend` that routes persistent paths to Oracle and keeps ephemeral data in memory:

| Path prefix | Backend | Persistence |
|------------|---------|-------------|
| `/memory/*` | Oracle `da_files` table | Durable, cross-session |
| `/history/*` | Oracle `da_files` table | Durable, cross-session |
| `/skills/*` | Oracle `da_files` table | Durable, cross-session |
| Everything else | `StateBackend` (in-memory) | Ephemeral, per-run |

## Explicit Wiring

For more control, wire the backends yourself:

```python
from deepagents import create_deep_agent
from deepagents.backends import StateBackend, CompositeBackend
from deepagents_oracle import OracleStoreBackend, OracleConfig

config = OracleConfig()  # reads DEEPAGENTS_* env vars
oracle = OracleStoreBackend(config)

backend = CompositeBackend(
    default=StateBackend,
    routes={"/memory/": oracle, "/history/": oracle, "/skills/": oracle},
)

agent = create_deep_agent(backend=backend, memory=["/memory/AGENTS.md"])
```

## Configuration

All config via environment variables with `DEEPAGENTS_` prefix:

```bash
# Mode: freepdb (local) or adb (cloud)
DEEPAGENTS_ORACLE_MODE=freepdb

# Credentials
DEEPAGENTS_ORACLE_USER=deepagents
DEEPAGENTS_ORACLE_PASSWORD=DeepAgents2024

# FreePDB connection
DEEPAGENTS_ORACLE_HOST=localhost
DEEPAGENTS_ORACLE_PORT=1521
DEEPAGENTS_ORACLE_SERVICE=FREEPDB1

# ADB connection (set these instead for cloud)
# DEEPAGENTS_ORACLE_DSN=(description=...)
# DEEPAGENTS_ORACLE_WALLET_PATH=/path/to/wallet
```

See [`oracle/.env.example`](oracle/.env.example) for the full template.

## GitHub Codespaces

[![Open in GitHub Codespaces](https://github.com/codespaces/badge.svg)](https://codespaces.new/jasperan/deepagents-oracle?quickstart=1)

The `.devcontainer/` config spins up Oracle 26ai Free automatically. Open the repo in Codespaces and start coding.

## Project Structure

```
libs/partners/oracle/     # Partner package (pip install deepagents-oracle)
  deepagents_oracle/
    config.py              # OracleConfig (Pydantic BaseSettings)
    connection.py          # Connection pool manager
    backend.py             # OracleStoreBackend (BackendProtocol)
    schema.py              # DDL for da_files table
    agent.py               # create_oracle_deep_agent() factory
oracle/                    # Docker infrastructure
  docker-compose.yml       # FreePDB + ADB profiles
  scripts/setup-oracle.sh  # User creation + schema init
examples/oracle_agent/     # Working example
```

## Full Deep Agents Docs

This fork inherits all upstream Deep Agents capabilities: planning, filesystem tools, sub-agents, context management, and more.

- [Deep Agents SDK Documentation](https://docs.langchain.com/oss/python/deepagents/overview)
- [Deep Agents CLI](https://docs.langchain.com/oss/python/deepagents/cli/overview)
- [Upstream Repository](https://github.com/langchain-ai/deepagents)

## License

MIT. See [LICENSE](LICENSE).
