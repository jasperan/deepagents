# Oracle Agent Example

A deep agent with persistent memory powered by Oracle AI Database.

## Quick Start

1. Start Oracle 26ai Free:
   ```bash
   docker compose -f oracle/docker-compose.yml --profile freepdb up -d
   bash oracle/scripts/setup-oracle.sh
   ```

2. Install dependencies:
   ```bash
   pip install deepagents deepagents-oracle
   ```

3. Run the agent:
   ```bash
   python agent.py
   ```

4. Stop the agent, restart it. Memory persists in Oracle.

## What This Demonstrates

- Agent memory stored in Oracle `da_files` table
- Cross-session persistence (restart the agent, memory is still there)
- CompositeBackend routing: `/memory/*` goes to Oracle, temp files stay in-memory
