"""Convenience factory for creating Oracle-backed deep agents."""

from __future__ import annotations

from typing import TYPE_CHECKING

from deepagents import create_deep_agent
from deepagents.backends import CompositeBackend, StateBackend

from deepagents_oracle.backend import OracleStoreBackend
from deepagents_oracle.config import OracleConfig
from deepagents_oracle.connection import OracleConnectionManager

if TYPE_CHECKING:
    from collections.abc import Callable, Sequence
    from typing import Any

    from deepagents.backends.protocol import BackendProtocol
    from langchain_core.language_models import BaseChatModel
    from langchain_core.tools import BaseTool

DEFAULT_MODEL = "ollama:qwopus3.5:9b-v3"
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
    **kwargs: object,
) -> object:
    """Create a deep agent with Oracle AI Database persistence.

    Sets up a ``CompositeBackend`` that routes persistent paths
    (memory, history, skills) to Oracle and ephemeral paths to
    ``StateBackend``.

    Args:
        model: The model to use. Defaults to ``ollama:qwopus3.5:9b-v3``.
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
        connection_manager=conn_manager,
        namespace=namespace,
    )

    routes: dict[str, BackendProtocol] = dict.fromkeys(persistent_routes, oracle_backend)

    def _backend_factory(runtime: object) -> CompositeBackend:
        return CompositeBackend(
            default=StateBackend(runtime),  # type: ignore[arg-type]
            routes=routes,
        )

    if model is None:
        model = DEFAULT_MODEL
    if memory is None:
        memory = ["/memory/AGENTS.md"]
    if skills is None:
        skills = ["/skills/"]

    return create_deep_agent(
        model=model,
        tools=tools,
        backend=_backend_factory,
        memory=memory,
        skills=skills,
        **kwargs,  # type: ignore[arg-type]
    )
