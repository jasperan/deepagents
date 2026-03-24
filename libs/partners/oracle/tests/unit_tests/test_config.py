"""Tests for OracleConfig."""

from deepagents_oracle.config import OracleConfig


class TestOracleConfigDefaults:
    def test_default_mode_is_freepdb(self):
        config = OracleConfig()
        assert config.oracle_mode == "freepdb"

    def test_default_connection_params(self):
        config = OracleConfig()
        assert config.oracle_user == "deepagents"
        assert config.oracle_host == "localhost"
        assert config.oracle_port == 1521
        assert config.oracle_service == "FREEPDB1"

    def test_is_adb_false_for_freepdb(self):
        config = OracleConfig()
        assert config.is_adb is False

    def test_is_adb_true_for_adb_mode(self):
        config = OracleConfig(oracle_mode="adb")
        assert config.is_adb is True

    def test_get_dsn_freepdb(self):
        config = OracleConfig()
        dsn = config.get_dsn()
        assert "localhost" in dsn
        assert "1521" in dsn
        assert "FREEPDB1" in dsn

    def test_get_dsn_adb_uses_oracle_dsn_field(self):
        descriptor = "(description=(address=(protocol=tcps)(host=adb.us-ashburn-1.oraclecloud.com)(port=1522))(connect_data=(service_name=abc_high)))"
        config = OracleConfig(oracle_mode="adb", oracle_dsn=descriptor)
        assert config.get_dsn() == descriptor

    def test_uses_wallet(self):
        config = OracleConfig(oracle_mode="adb", oracle_wallet_path="/path/to/wallet")
        assert config.uses_wallet is True

    def test_no_wallet_by_default(self):
        config = OracleConfig()
        assert config.uses_wallet is False


class TestOracleConfigFromEnv:
    def test_reads_from_env_with_prefix(self, monkeypatch):
        monkeypatch.setenv("DEEPAGENTS_ORACLE_MODE", "adb")
        monkeypatch.setenv("DEEPAGENTS_ORACLE_USER", "myuser")
        monkeypatch.setenv("DEEPAGENTS_ORACLE_PASSWORD", "mypass")
        config = OracleConfig()
        assert config.oracle_mode == "adb"
        assert config.oracle_user == "myuser"
        assert config.oracle_password == "mypass"  # noqa: S105
