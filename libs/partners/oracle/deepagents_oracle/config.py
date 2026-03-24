"""Oracle Database configuration via environment variables."""

from typing import Literal

from pydantic_settings import BaseSettings


class OracleConfig(BaseSettings):
    """Configuration for Oracle AI Database connection.

    Reads from environment variables with ``DEEPAGENTS_`` prefix.
    Supports two modes: ``freepdb`` (local 26ai Free container) and
    ``adb`` (Oracle Autonomous Database in the cloud).
    """

    model_config = {"env_prefix": "DEEPAGENTS_"}

    oracle_mode: Literal["freepdb", "adb"] = "freepdb"
    oracle_user: str = "deepagents"
    oracle_password: str = "DeepAgents2024"
    oracle_host: str = "localhost"
    oracle_port: int = 1521
    oracle_service: str = "FREEPDB1"
    oracle_dsn: str | None = None
    oracle_wallet_path: str | None = None
    oracle_pool_min: int = 2
    oracle_pool_max: int = 10

    @property
    def is_adb(self) -> bool:
        """Whether this config targets Autonomous Database."""
        return self.oracle_mode == "adb"

    @property
    def uses_wallet(self) -> bool:
        """Whether a wallet path is configured for mTLS."""
        return self.oracle_wallet_path is not None

    def get_dsn(self) -> str:
        """Build the connection DSN string.

        For ``freepdb`` mode, constructs ``host:port/service``.
        For ``adb`` mode, returns the ``oracle_dsn`` field directly.

        Raises:
            ValueError: If ``adb`` mode but no ``oracle_dsn`` provided.
        """
        if self.is_adb:
            if self.oracle_dsn is None:
                msg = "oracle_dsn is required when oracle_mode='adb'"
                raise ValueError(msg)
            return self.oracle_dsn
        return f"{self.oracle_host}:{self.oracle_port}/{self.oracle_service}"
