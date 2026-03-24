"""Convenience factory for creating Oracle-backed deep agents."""

from __future__ import annotations

from collections.abc import Callable, Sequence
from typing import Any

from deepagents import create_deep_agent
from deepagents.backends import CompositeBackend, StateBackend
from langchain_core.language_models import BaseChatModel
from langchain_core.tools import BaseTool

from deepagents_oracle.backend import OracleStoreBackend
from deepagents_oracle.config import OracleConfig
from deepagents_oracle.connection import OracleConnectionManager

DEFAULT_PERSISTENT_ROUTES = ("/memory/", "/history/", "/skills/")


def create_oracle_deep_agent(
    model: str | BaseChatModel | None = None,
    tools: Sequence[BaseTool | Callable | dict[str, Any]] | None = None,
    *,
    oracle_config: OracleConfig | None = None,
    namespace: str = "default",
    persistent_routes: Sequence[str] = DEFAULT_PERSISTENT_ROUTES,
    memory: list[str] | None = None,
    skills: list[str] | None = None,
    **kwargs: Any,
) -> Any:
    """Create a deep agent with Oracle AI Database persistence.

    Sets up a ``CompositeBackend`` that routes persistent paths
    (memory, history, skills) to Oracle and ephemeral paths to
    ``StateBackend``.

    Args:
        model: The model to use. Defaults to claude-sonnet-4-6.
        tools: Custom tools for the agent.
        oracle_config: Oracle connection config. Reads from env if None.
        namespace: Namespace for data isolation between agents.
        persistent_routes: Path prefixes routed to Oracle.
        memory: Memory file paths. Defaults to ``["/memory/AGENTS.md"]``.
        skills: Skill source paths. Defaults to ``["/skills/"]``.
        **kwargs: Additional arguments passed to ``create_deep_agent``.

    Returns:
        A configured deep agent with Oracle persistence.
    """
    if oracle_config is None:
        oracle_config = OracleConfig()

    conn_manager = OracleConnectionManager(oracle_config)
    conn_manager.connect()

    oracle_backend = OracleStoreBackend(
        config=oracle_config,
        connection_manager=conn_manager,
        namespace=namespace,
    )

    routes = {route: oracle_backend for route in persistent_routes}

    backend = CompositeBackend(
        default=StateBackend,
        routes=routes,
    )

    if memory is None:
        memory = ["/memory/AGENTS.md"]
    if skills is None:
        skills = ["/skills/"]

    return create_deep_agent(
        model=model,
        tools=tools,
        backend=backend,
        memory=memory,
        skills=skills,
        **kwargs,
    )
