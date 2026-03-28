"""Integration tests against a live Oracle 26ai Free instance.

These tests require a running Oracle container with:
- Host: localhost:1521/FREEPDB1
- User: deepagents / DeepAgents2024

Run: cd libs/partners/oracle && uv run --group test pytest tests/integration_tests/ -v --timeout 30
"""

import logging
import time

import pytest

from deepagents_oracle.backend import OracleStoreBackend
from deepagents_oracle.config import OracleConfig
from deepagents_oracle.connection import OracleConnectionManager
from deepagents_oracle.schema import init_schema

# Enable logging to verify it works
logging.basicConfig(level=logging.DEBUG)
logger = logging.getLogger(__name__)

# ── Fixtures ──────────────────────────────────────────────────────────

ORACLE_PASSWORD = "DeepAgents2024"  # noqa: S105


@pytest.fixture(scope="session")
def oracle_config():
    return OracleConfig(
        oracle_mode="freepdb",
        oracle_user="deepagents",
        oracle_password=ORACLE_PASSWORD,
        oracle_host="localhost",
        oracle_port=1521,
        oracle_service="FREEPDB1",
    )


@pytest.fixture(scope="session")
def connection_manager(oracle_config):
    mgr = OracleConnectionManager(oracle_config)
    mgr.connect()
    yield mgr
    mgr.close()


@pytest.fixture(scope="session")
def schema_ready(connection_manager):
    """Initialize schema once for all tests."""
    with connection_manager.get_connection() as conn:
        init_schema(conn, vector=False)
    return True


@pytest.fixture
def backend(connection_manager, schema_ready):  # noqa: ARG001
    """Fresh backend with a unique namespace per test."""
    ns = f"test_{int(time.time() * 1000)}"
    be = OracleStoreBackend(connection_manager=connection_manager, namespace=ns)
    yield be
    # Cleanup: delete all rows in this namespace
    with connection_manager.get_connection() as conn:
        with conn.cursor() as cur:
            cur.execute("DELETE FROM da_files WHERE namespace = :ns", {"ns": ns})
        conn.commit()


# ── Connection & Health ──────────────────────────────────────────────


class TestConnection:
    def test_health_check(self, connection_manager):
        assert connection_manager.health_check() is True

    def test_dsn(self, oracle_config):
        assert oracle_config.get_dsn() == "localhost:1521/FREEPDB1"

    def test_pool_is_active(self, connection_manager):
        with connection_manager.get_connection() as conn, conn.cursor() as cur:
            cur.execute("SELECT 1 FROM DUAL")
            assert cur.fetchone()[0] == 1


# ── Schema ───────────────────────────────────────────────────────────


class TestSchema:
    def test_table_exists(self, connection_manager, schema_ready):
        with connection_manager.get_connection() as conn, conn.cursor() as cur:
            cur.execute("""
                    SELECT COUNT(*) FROM user_tables WHERE table_name = 'DA_FILES'
                """)
            assert cur.fetchone()[0] == 1

    def test_columns_exist(self, connection_manager, schema_ready):
        with connection_manager.get_connection() as conn, conn.cursor() as cur:
            cur.execute("""
                    SELECT column_name FROM user_tab_columns
                    WHERE table_name = 'DA_FILES'
                    ORDER BY column_id
                """)
            columns = [row[0] for row in cur.fetchall()]
            assert "NAMESPACE" in columns
            assert "FILE_PATH" in columns
            assert "CONTENT" in columns
            assert "ENCODING" in columns
            assert "CREATED_AT" in columns
            assert "MODIFIED_AT" in columns

    def test_idempotent_init(self, connection_manager):
        """Calling init_schema twice should not raise."""
        with connection_manager.get_connection() as conn:
            init_schema(conn, vector=False)
            init_schema(conn, vector=False)


# ── Write ────────────────────────────────────────────────────────────


class TestWrite:
    def test_write_new_file(self, backend):
        result = backend.write("/memory/test.md", "Hello Oracle!")
        assert result.path == "/memory/test.md"
        assert result.files_update is None

    def test_write_overwrites_existing(self, backend):
        backend.write("/memory/test.md", "version 1")
        backend.write("/memory/test.md", "version 2")
        read_result = backend.read("/memory/test.md")
        content = read_result.file_data["content"]
        assert "version 2" in content

    def test_write_unicode(self, backend):
        backend.write("/memory/unicode.md", "Hello 世界 🌍 émojis café")
        read_result = backend.read("/memory/unicode.md")
        content = read_result.file_data["content"]
        assert "世界" in content
        assert "café" in content

    def test_write_large_content(self, backend):
        large = "x" * 100_000  # 100KB
        backend.write("/memory/large.md", large)
        read_result = backend.read("/memory/large.md")
        content = read_result.file_data["content"]
        assert len(content) == 100_000

    def test_write_empty_content(self, backend):
        backend.write("/memory/empty.md", "")
        read_result = backend.read("/memory/empty.md")
        assert read_result.file_data is not None

    def test_write_multiline(self, backend):
        content = "line 1\nline 2\nline 3\n"
        backend.write("/memory/multiline.md", content)
        read_result = backend.read("/memory/multiline.md")
        assert "line 2" in read_result.file_data["content"]


# ── Read ─────────────────────────────────────────────────────────────


class TestRead:
    def test_read_existing(self, backend):
        backend.write("/memory/read_test.md", "read me")
        result = backend.read("/memory/read_test.md")
        assert result.error is None
        assert "read me" in result.file_data["content"]

    def test_read_missing_file(self, backend):
        result = backend.read("/memory/nonexistent.md")
        assert result.error is not None
        assert "not found" in result.error.lower()

    def test_read_preserves_encoding(self, backend):
        backend.write("/memory/enc.md", "utf-8 content")
        result = backend.read("/memory/enc.md")
        assert result.file_data["encoding"] == "utf-8"


# ── Edit ─────────────────────────────────────────────────────────────


class TestEdit:
    def test_edit_replaces_content(self, backend):
        backend.write("/memory/edit.md", "hello world")
        result = backend.edit("/memory/edit.md", "hello", "goodbye")
        assert result.error is None
        # Verify the change stuck
        read_result = backend.read("/memory/edit.md")
        assert "goodbye world" in read_result.file_data["content"]

    def test_edit_missing_file(self, backend):
        result = backend.edit("/memory/nope.md", "a", "b")
        assert result.error is not None

    def test_edit_string_not_found(self, backend):
        backend.write("/memory/edit2.md", "hello world")
        result = backend.edit("/memory/edit2.md", "nonexistent", "replacement")
        assert result.error is not None

    def test_edit_multiple_occurrences(self, backend):
        backend.write("/memory/multi.md", "aaa bbb aaa")
        result = backend.edit("/memory/multi.md", "aaa", "ccc", replace_all=True)
        assert result.error is None
        read_result = backend.read("/memory/multi.md")
        assert "ccc bbb ccc" in read_result.file_data["content"]


# ── Ls ───────────────────────────────────────────────────────────────


class TestLs:
    def test_ls_directory(self, backend):
        backend.write("/memory/a.md", "a")
        backend.write("/memory/b.md", "b")
        backend.write("/memory/sub/c.md", "c")
        result = backend.ls("/memory/")
        paths = [e.get("path", "") for e in result.entries]
        assert any("a.md" in p for p in paths)
        assert any("b.md" in p for p in paths)
        # sub/ should appear as a directory
        assert any("sub" in p for p in paths)

    def test_ls_empty_directory(self, backend):
        result = backend.ls("/nonexistent/")
        assert len(result.entries) == 0


# ── Glob ─────────────────────────────────────────────────────────────


class TestGlob:
    def test_glob_matches(self, backend):
        backend.write("/memory/notes.md", "notes")
        backend.write("/memory/todo.txt", "todo")
        result = backend.glob("*.md", "/memory/")
        paths = [m.get("path", "") for m in result.matches]
        assert any("notes.md" in p for p in paths)
        assert not any("todo.txt" in p for p in paths)

    def test_glob_no_matches(self, backend):
        result = backend.glob("*.xyz", "/memory/")
        assert len(result.matches) == 0


# ── Grep ─────────────────────────────────────────────────────────────


class TestGrep:
    def test_grep_finds_content(self, backend):
        backend.write("/memory/searchme.md", "the quick brown fox")
        result = backend.grep("quick", "/memory/")
        assert len(result.matches) >= 1

    def test_grep_no_matches(self, backend):
        backend.write("/memory/nope.md", "nothing here")
        result = backend.grep("zzzznotfound", "/memory/")
        assert len(result.matches) == 0


# ── Upload / Download ────────────────────────────────────────────────


class TestUploadDownload:
    def test_upload_and_download(self, backend):
        backend.upload_files([("/memory/uploaded.md", b"uploaded content")])
        results = backend.download_files(["/memory/uploaded.md"])
        assert len(results) == 1
        assert results[0].error is None
        assert results[0].content == b"uploaded content"

    def test_download_missing(self, backend):
        results = backend.download_files(["/memory/nope.md"])
        assert results[0].error == "file_not_found"

    def test_upload_multiple(self, backend):
        files = [
            ("/memory/f1.md", b"file one"),
            ("/memory/f2.md", b"file two"),
            ("/memory/f3.md", b"file three"),
        ]
        results = backend.upload_files(files)
        assert all(r.error is None for r in results)
        dl = backend.download_files(["/memory/f1.md", "/memory/f2.md", "/memory/f3.md"])
        assert all(d.error is None for d in dl)


# ── Namespace Isolation ──────────────────────────────────────────────


class TestNamespaceIsolation:
    def test_different_namespaces_are_isolated(self, connection_manager, schema_ready):
        ns1 = f"iso_a_{int(time.time() * 1000)}"
        ns2 = f"iso_b_{int(time.time() * 1000)}"
        be1 = OracleStoreBackend(connection_manager=connection_manager, namespace=ns1)
        be2 = OracleStoreBackend(connection_manager=connection_manager, namespace=ns2)

        be1.write("/memory/secret.md", "namespace 1 only")
        be2.write("/memory/secret.md", "namespace 2 only")

        r1 = be1.read("/memory/secret.md")
        r2 = be2.read("/memory/secret.md")

        assert "namespace 1" in r1.file_data["content"]
        assert "namespace 2" in r2.file_data["content"]

        # Cleanup
        with connection_manager.get_connection() as conn:
            with conn.cursor() as cur:
                cur.execute("DELETE FROM da_files WHERE namespace IN (:a, :b)", {"a": ns1, "b": ns2})
            conn.commit()


# ── Repeated Cycles (Stress) ─────────────────────────────────────────


class TestStress:
    def test_repeated_write_read_cycles(self, backend):
        """Write and read 50 files rapidly."""
        for i in range(50):
            path = f"/memory/stress_{i}.md"
            content = f"Stress test file {i} with content {i * 42}"
            backend.write(path, content)
            result = backend.read(path)
            assert result.error is None
            assert str(i * 42) in result.file_data["content"]

    def test_rapid_edit_cycles(self, backend):
        """Edit the same file 20 times."""
        backend.write("/memory/rapid.md", "version_0")
        for i in range(1, 21):
            backend.edit("/memory/rapid.md", f"version_{i - 1}", f"version_{i}")
        result = backend.read("/memory/rapid.md")
        assert "version_20" in result.file_data["content"]

    def test_concurrent_namespaces(self, connection_manager, schema_ready):
        """10 different namespaces writing simultaneously."""
        backends = []
        namespaces = []
        for i in range(10):
            ns = f"conc_{i}_{int(time.time() * 1000)}"
            namespaces.append(ns)
            backends.append(OracleStoreBackend(connection_manager=connection_manager, namespace=ns))

        # Write to all
        for i, be in enumerate(backends):
            be.write("/memory/test.md", f"namespace {i}")

        # Read from all and verify isolation
        for i, be in enumerate(backends):
            result = be.read("/memory/test.md")
            assert f"namespace {i}" in result.file_data["content"]

        # Cleanup
        with connection_manager.get_connection() as conn:
            with conn.cursor() as cur:
                for ns in namespaces:
                    cur.execute("DELETE FROM da_files WHERE namespace = :ns", {"ns": ns})
            conn.commit()


# ── Logging Verification ─────────────────────────────────────────────


class TestLogging:
    def test_schema_init_logs(self, connection_manager, caplog):
        with caplog.at_level(logging.INFO, logger="deepagents_oracle.schema"):
            with connection_manager.get_connection() as conn:
                init_schema(conn, vector=False)
            assert any("da_files" in r.message.lower() or "ready" in r.message.lower() for r in caplog.records)

    def test_connection_pool_logs(self, oracle_config, caplog):
        with caplog.at_level(logging.INFO, logger="deepagents_oracle.connection"):
            mgr = OracleConnectionManager(oracle_config)
            mgr.connect()
            mgr.close()
        assert any("pool" in r.message.lower() for r in caplog.records)
