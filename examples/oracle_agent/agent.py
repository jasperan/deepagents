"""Example: Deep Agent with Oracle AI Database persistence.

Demonstrates memory that survives across sessions. Start the agent,
give it some research tasks, stop it, restart it, and watch it
remember what it learned.

Usage:
    # Start Oracle (from repo root)
    docker compose -f oracle/docker-compose.yml --profile freepdb up -d
    bash oracle/scripts/setup-oracle.sh

    # Run the agent
    python agent.py
"""

from deepagents_oracle import create_oracle_deep_agent

agent = create_oracle_deep_agent(
    system_prompt="You are a research assistant. Store important findings in /memory/ files so you remember them across sessions.",
)

if __name__ == "__main__":
    result = agent.invoke(
        {"messages": [{"role": "user", "content": "What are the key features of Oracle AI Vector Search? Save your findings to /memory/oracle-vectors.md"}]},
    )
    for msg in result["messages"]:
        if hasattr(msg, "content") and msg.content:
            print(msg.content)
