"""Oracle AI Database backend for Deep Agents."""

from deepagents_oracle._version import __version__
from deepagents_oracle.agent import DEFAULT_MODEL, SKILLS_DIR, create_oracle_deep_agent
from deepagents_oracle.backend import OracleStoreBackend
from deepagents_oracle.config import OracleConfig
from deepagents_oracle.connection import OracleConnectionManager

__all__ = ["DEFAULT_MODEL", "SKILLS_DIR", "OracleConfig", "OracleConnectionManager", "OracleStoreBackend", "__version__", "create_oracle_deep_agent"]

try:
    from deepagents_oracle.vector import OracleVectorBackend
except ImportError:
    pass
else:
    __all__ = [*__all__, "OracleVectorBackend"]  # noqa: PLE0604
