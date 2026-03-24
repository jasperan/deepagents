"""Tests for create_oracle_deep_agent factory."""

from unittest.mock import MagicMock, patch

from deepagents_oracle.agent import create_oracle_deep_agent
from deepagents_oracle.config import OracleConfig


class TestCreateOracleDeepAgent:
    @patch("deepagents_oracle.agent.create_deep_agent")
    @patch("deepagents_oracle.agent.OracleStoreBackend")
    @patch("deepagents_oracle.agent.OracleConnectionManager")
    def test_creates_agent_with_oracle_backend(self, mock_conn_cls, mock_backend_cls, mock_create):
        mock_create.return_value = MagicMock()
        mock_backend_cls.return_value = MagicMock()
        mock_mgr = MagicMock()
        mock_conn_cls.return_value = mock_mgr

        create_oracle_deep_agent()

        mock_conn_cls.assert_called_once()
        mock_mgr.connect.assert_called_once()
        mock_backend_cls.assert_called_once()
        mock_create.assert_called_once()
        call_kwargs = mock_create.call_args.kwargs
        assert "backend" in call_kwargs
        assert "memory" in call_kwargs
        assert "skills" in call_kwargs

    @patch("deepagents_oracle.agent.create_deep_agent")
    @patch("deepagents_oracle.agent.OracleStoreBackend")
    @patch("deepagents_oracle.agent.OracleConnectionManager")
    def test_passes_through_model_and_tools(self, mock_conn_cls, mock_backend_cls, mock_create):
        mock_create.return_value = MagicMock()
        mock_backend_cls.return_value = MagicMock()
        mock_conn_cls.return_value = MagicMock()

        my_tool = MagicMock()
        create_oracle_deep_agent(model="openai:gpt-4o", tools=[my_tool])

        call_kwargs = mock_create.call_args.kwargs
        assert call_kwargs["model"] == "openai:gpt-4o"
        assert my_tool in call_kwargs["tools"]

    @patch("deepagents_oracle.agent.create_deep_agent")
    @patch("deepagents_oracle.agent.OracleStoreBackend")
    @patch("deepagents_oracle.agent.OracleConnectionManager")
    def test_accepts_custom_oracle_config(self, mock_conn_cls, mock_backend_cls, mock_create):
        mock_create.return_value = MagicMock()
        mock_backend_cls.return_value = MagicMock()
        mock_conn_cls.return_value = MagicMock()

        config = OracleConfig(oracle_mode="adb", oracle_dsn="(description=...)")
        create_oracle_deep_agent(oracle_config=config)

        mock_conn_cls.assert_called_once_with(config)

    @patch("deepagents_oracle.agent.create_deep_agent")
    @patch("deepagents_oracle.agent.OracleStoreBackend")
    @patch("deepagents_oracle.agent.OracleConnectionManager")
    def test_custom_namespace(self, mock_conn_cls, mock_backend_cls, mock_create):
        mock_create.return_value = MagicMock()
        mock_backend_cls.return_value = MagicMock()
        mock_conn_cls.return_value = MagicMock()

        create_oracle_deep_agent(namespace="my-agent")

        backend_kwargs = mock_backend_cls.call_args.kwargs
        assert backend_kwargs["namespace"] == "my-agent"

    @patch("deepagents_oracle.agent.create_deep_agent")
    @patch("deepagents_oracle.agent.OracleStoreBackend")
    @patch("deepagents_oracle.agent.OracleConnectionManager")
    def test_passes_extra_kwargs_to_create_deep_agent(self, mock_conn_cls, mock_backend_cls, mock_create):
        mock_create.return_value = MagicMock()
        mock_backend_cls.return_value = MagicMock()
        mock_conn_cls.return_value = MagicMock()

        create_oracle_deep_agent(system_prompt="You are helpful", debug=True)

        call_kwargs = mock_create.call_args.kwargs
        assert call_kwargs["system_prompt"] == "You are helpful"
        assert call_kwargs["debug"] is True
