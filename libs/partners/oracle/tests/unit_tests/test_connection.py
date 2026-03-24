"""Tests for OracleConnectionManager."""

from unittest.mock import MagicMock, patch

import pytest

from deepagents_oracle.config import OracleConfig
from deepagents_oracle.connection import OracleConnectionManager


class TestConnectionManagerInit:
    def test_creates_manager_from_config(self):
        config = OracleConfig()
        manager = OracleConnectionManager(config)
        assert manager.config is config
        assert manager._pool is None

    def test_dsn_from_config(self):
        config = OracleConfig()
        manager = OracleConnectionManager(config)
        assert manager.dsn == "localhost:1521/FREEPDB1"


class TestConnectionManagerPool:
    @patch("deepagents_oracle.connection.oracledb")
    def test_connect_creates_pool(self, mock_oracledb):
        mock_pool = MagicMock()
        mock_oracledb.create_pool.return_value = mock_pool

        config = OracleConfig()
        manager = OracleConnectionManager(config)
        manager.connect()

        mock_oracledb.create_pool.assert_called_once_with(
            user=config.oracle_user,
            password=config.oracle_password,
            dsn=config.get_dsn(),
            min=config.oracle_pool_min,
            max=config.oracle_pool_max,
        )
        assert manager._pool is mock_pool

    @patch("deepagents_oracle.connection.oracledb")
    def test_connect_adb_with_wallet(self, mock_oracledb):
        mock_pool = MagicMock()
        mock_oracledb.create_pool.return_value = mock_pool

        config = OracleConfig(
            oracle_mode="adb",
            oracle_dsn="(description=(address=...))",
            oracle_wallet_path="/wallet",
        )
        manager = OracleConnectionManager(config)
        manager.connect()

        call_kwargs = mock_oracledb.create_pool.call_args.kwargs
        assert call_kwargs["dsn"] == config.oracle_dsn
        assert "wallet_location" in call_kwargs
        assert call_kwargs["wallet_location"] == "/wallet"

    @patch("deepagents_oracle.connection.oracledb")
    def test_close_shuts_down_pool(self, mock_oracledb):
        mock_pool = MagicMock()
        mock_oracledb.create_pool.return_value = mock_pool

        config = OracleConfig()
        manager = OracleConnectionManager(config)
        manager.connect()
        manager.close()

        mock_pool.close.assert_called_once()
        assert manager._pool is None

    @patch("deepagents_oracle.connection.oracledb")
    def test_get_connection_context_manager(self, mock_oracledb):
        mock_pool = MagicMock()
        mock_conn = MagicMock()
        mock_pool.acquire.return_value = mock_conn
        mock_oracledb.create_pool.return_value = mock_pool

        config = OracleConfig()
        manager = OracleConnectionManager(config)
        manager.connect()

        with manager.get_connection() as conn:
            assert conn is mock_conn

    @patch("deepagents_oracle.connection.oracledb")
    def test_get_connection_raises_without_connect(self, mock_oracledb):
        config = OracleConfig()
        manager = OracleConnectionManager(config)

        with pytest.raises(RuntimeError, match="Connection pool not initialized"), manager.get_connection():
            pass

    @patch("deepagents_oracle.connection.oracledb")
    def test_connect_is_idempotent(self, mock_oracledb):
        mock_pool = MagicMock()
        mock_oracledb.create_pool.return_value = mock_pool

        config = OracleConfig()
        manager = OracleConnectionManager(config)
        manager.connect()
        manager.connect()  # second call should be no-op

        mock_oracledb.create_pool.assert_called_once()
