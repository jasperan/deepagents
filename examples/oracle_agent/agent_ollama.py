"""Deep Agent with qwopus3.5 via Ollama + Oracle AI Database persistence.

Uses a local Ollama model instead of a cloud API. Memory persists in Oracle.

Usage:
    # Ensure Ollama is running with the model:
    #   ollama pull qwopus3.5:9b-v3
    #
    # Ensure Oracle is running:
    #   docker compose -f oracle/docker-compose.yml --profile freepdb up -d
    #   bash oracle/scripts/setup-oracle.sh
    #
    # Set env vars (or source oracle/.env):
    #   export DEEPAGENTS_ORACLE_MODE=freepdb
    #   export DEEPAGENTS_ORACLE_USER=deepagents
    #   export DEEPAGENTS_ORACLE_PASSWORD=DeepAgents2024
    #   export DEEPAGENTS_ORACLE_HOST=localhost
    #   export DEEPAGENTS_ORACLE_PORT=1521
    #   export DEEPAGENTS_ORACLE_SERVICE=FREEPDB1
    #
    # Run:
    #   python agent_ollama.py
"""

from deepagents_oracle import create_oracle_deep_agent
from deepagents_oracle.agent import DEFAULT_MODEL

# No model= needed: create_oracle_deep_agent defaults to qwopus3.5 via Ollama.
# Override with model="anthropic:claude-sonnet-4-6" for cloud inference.
agent = create_oracle_deep_agent(
    system_prompt=(
        "You are a helpful research assistant running locally on qwopus3.5 via Ollama. "
        "You have persistent memory in Oracle Database. "
        "Store important findings in /memory/ files so you remember them across sessions. "
        "Keep responses concise."
    ),
)

if __name__ == "__main__":
    print(f"Deep Agent ready with model: {DEFAULT_MODEL}")
    print("Oracle persistence enabled. Memory survives restarts.\n")

    prompts = [
        "What is 2+2? Save the answer to /memory/math.md",
        "Read /memory/math.md and tell me what you saved earlier.",
    ]

    for prompt in prompts:
        print(f">>> {prompt}")
        result = agent.invoke(
            {"messages": [{"role": "user", "content": prompt}]},
        )
        for msg in result["messages"]:
            if hasattr(msg, "content") and msg.content:
                print(msg.content)
        print()
