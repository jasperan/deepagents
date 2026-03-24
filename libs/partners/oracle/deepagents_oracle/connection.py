"""Oracle Database connection pool management."""

from __future__ import annotations

import contextlib
import logging
from collections.abc import Generator
from typing import Any

import oracledb

from deepagents_oracle.config import OracleConfig

logger = logging.getLogger(__name__)


class OracleConnectionManager:
    """Manages an oracledb connection pool for FreePDB or ADB.

    Args:
        config: Oracle configuration with connection parameters.
    """

    def __init__(self, config: OracleConfig) -> None:
        self.config = config
        self.dsn = config.get_dsn()
        self._pool: oracledb.ConnectionPool | None = None

    def connect(self) -> None:
        """Create the connection pool.

        For ADB with wallet, passes ``wallet_location`` and
        ``wallet_password`` to the pool constructor.
        """
        if self._pool is not None:
            return

        kwargs: dict[str, Any] = {
            "user": self.config.oracle_user,
            "password": self.config.oracle_password,
            "dsn": self.dsn,
            "min": self.config.oracle_pool_min,
            "max": self.config.oracle_pool_max,
        }

        if self.config.uses_wallet:
            kwargs["wallet_location"] = self.config.oracle_wallet_path
            kwargs["wallet_password"] = self.config.oracle_password

        self._pool = oracledb.create_pool(**kwargs)
        logger.info(
            "Oracle connection pool created (min=%d, max=%d)",
            self.config.oracle_pool_min,
            self.config.oracle_pool_max,
        )

    def close(self) -> None:
        """Shut down the connection pool."""
        if self._pool is not None:
            self._pool.close()
            self._pool = None
            logger.info("Oracle connection pool closed")

    @contextlib.contextmanager
    def get_connection(self) -> Generator[oracledb.Connection, None, None]:
        """Acquire a connection from the pool, release on exit.

        Raises:
            RuntimeError: If the pool has not been created via ``connect()``.
        """
        if self._pool is None:
            msg = "Connection pool not initialized. Call connect() first."
            raise RuntimeError(msg)
        conn = self._pool.acquire()
        try:
            yield conn
        finally:
            self._pool.release(conn)

    def health_check(self) -> bool:
        """Run SELECT 1 FROM DUAL to verify connectivity."""
        try:
            with self.get_connection() as conn:
                with conn.cursor() as cur:
                    cur.execute("SELECT 1 FROM DUAL")
                    row = cur.fetchone()
                    return row is not None and row[0] == 1
        except Exception:
            logger.exception("Oracle health check failed")
            return False
