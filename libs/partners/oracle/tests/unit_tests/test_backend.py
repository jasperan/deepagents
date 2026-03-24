"""Tests for OracleStoreBackend."""

from __future__ import annotations

from contextlib import contextmanager
from unittest.mock import MagicMock, patch

import pytest

from deepagents_oracle.backend import OracleStoreBackend
from deepagents_oracle.config import OracleConfig
from deepagents_oracle.connection import OracleConnectionManager


# ---------------------------------------------------------------------------
# Fixtures
# ---------------------------------------------------------------------------


@pytest.fixture()
def mock_cursor():
    """Return a MagicMock cursor that tracks execute/fetchone/fetchall calls."""
    cursor = MagicMock()
    cursor.description = None
    cursor.fetchone.return_value = None
    cursor.fetchall.return_value = []
    return cursor


@pytest.fixture()
def mock_conn(mock_cursor):
    """Return a MagicMock connection whose cursor() context manager yields *mock_cursor*."""
    conn = MagicMock()

    @contextmanager
    def _cursor_ctx():
        yield mock_cursor

    conn.cursor = _cursor_ctx
    return conn


@pytest.fixture()
def mock_manager(mock_conn):
    """Return a MagicMock OracleConnectionManager whose get_connection yields *mock_conn*."""
    manager = MagicMock(spec=OracleConnectionManager)

    @contextmanager
    def _conn_ctx():
        yield mock_conn

    manager.get_connection = _conn_ctx
    return manager


@pytest.fixture()
def backend(mock_manager):
    """Return an OracleStoreBackend wired to the mock manager."""
    return OracleStoreBackend(connection_manager=mock_manager, namespace="test_ns")


# ---------------------------------------------------------------------------
# write()
# ---------------------------------------------------------------------------


class TestWrite:
    def test_write_returns_correct_result(self, backend, mock_cursor):
        """write() should INSERT and return WriteResult with path set."""
        result = backend.write("/app/hello.txt", "hello world")

        assert result.error is None
        assert result.path == "/app/hello.txt"
        assert result.files_update is None
        mock_cursor.execute.assert_called_once()

    def test_write_existing_file_upserts(self, backend, mock_cursor):
        """write() uses MERGE so even existing files succeed (upsert semantics)."""
        result = backend.write("/app/hello.txt", "updated")

        assert result.error is None
        assert result.path == "/app/hello.txt"


# ---------------------------------------------------------------------------
# read()
# ---------------------------------------------------------------------------


class TestRead:
    def test_read_existing_file(self, backend, mock_cursor):
        """read() should return file content when the row exists."""
        mock_cursor.fetchone.return_value = (
            "line1\nline2\nline3",  # content (str, not LOB)
            "utf-8",
            "2025-01-01T00:00:00",
            "2025-01-01T00:00:00",
        )

        result = backend.read("/app/hello.txt")

        assert result.error is None
        assert result.file_data is not None
        assert "line1" in result.file_data["content"]

    def test_read_missing_file(self, backend, mock_cursor):
        """read() should return error for non-existent file."""
        mock_cursor.fetchone.return_value = None

        result = backend.read("/app/missing.txt")

        assert result.error is not None
        assert "not found" in result.error.lower()

    def test_read_clob_object(self, backend, mock_cursor):
        """read() should call .read() on LOB objects."""
        lob = MagicMock()
        lob.read.return_value = "clob content"

        mock_cursor.fetchone.return_value = (
            lob,
            "utf-8",
            "2025-01-01T00:00:00",
            "2025-01-01T00:00:00",
        )

        result = backend.read("/app/clob.txt")

        assert result.error is None
        lob.read.assert_called_once()
        assert "clob content" in result.file_data["content"]


# ---------------------------------------------------------------------------
# edit()
# ---------------------------------------------------------------------------


class TestEdit:
    def test_edit_replaces_content(self, backend, mock_cursor):
        """edit() should replace old_string with new_string and UPDATE."""
        mock_cursor.fetchone.return_value = (
            "hello world",
            "utf-8",
            "2025-01-01T00:00:00",
            "2025-01-01T00:00:00",
        )

        result = backend.edit("/app/hello.txt", "hello", "goodbye")

        assert result.error is None
        assert result.path == "/app/hello.txt"
        assert result.files_update is None
        assert result.occurrences == 1
        # Verify UPDATE was called (second execute call after SELECT)
        assert mock_cursor.execute.call_count == 2

    def test_edit_missing_file(self, backend, mock_cursor):
        """edit() should return error for non-existent file."""
        mock_cursor.fetchone.return_value = None

        result = backend.edit("/app/missing.txt", "old", "new")

        assert result.error is not None
        assert "not found" in result.error.lower()

    def test_edit_string_not_found(self, backend, mock_cursor):
        """edit() should return error when old_string is not in content."""
        mock_cursor.fetchone.return_value = (
            "hello world",
            "utf-8",
            "2025-01-01T00:00:00",
            "2025-01-01T00:00:00",
        )

        result = backend.edit("/app/hello.txt", "nonexistent", "replacement")

        assert result.error is not None
        assert "not found" in result.error.lower()


# ---------------------------------------------------------------------------
# ls()
# ---------------------------------------------------------------------------


class TestLs:
    def test_ls_lists_directory_entries(self, backend, mock_cursor):
        """ls() should return files and subdirectories in the given path."""
        mock_cursor.fetchall.return_value = [
            ("/app/file1.txt", "utf-8", "content1", "2025-01-01T00:00:00", "2025-01-01T00:00:00"),
            ("/app/file2.py", "utf-8", "content2", "2025-01-01T00:00:00", "2025-01-01T00:00:00"),
            ("/app/sub/deep.txt", "utf-8", "deep", "2025-01-01T00:00:00", "2025-01-01T00:00:00"),
        ]

        result = backend.ls("/app")

        assert result.error is None
        assert result.entries is not None
        paths = [e["path"] for e in result.entries]
        assert "/app/file1.txt" in paths
        assert "/app/file2.py" in paths
        # Subdirectory should appear with trailing slash
        assert any(e.get("is_dir") for e in result.entries)


# ---------------------------------------------------------------------------
# glob()
# ---------------------------------------------------------------------------


class TestGlob:
    def test_glob_matches_pattern(self, backend, mock_cursor):
        """glob() should return files matching the glob pattern."""
        mock_cursor.fetchall.return_value = [
            ("/app/main.py", "utf-8", "code", "2025-01-01T00:00:00", "2025-01-01T00:00:00"),
            ("/app/test.py", "utf-8", "test", "2025-01-01T00:00:00", "2025-01-01T00:00:00"),
            ("/app/readme.md", "utf-8", "docs", "2025-01-01T00:00:00", "2025-01-01T00:00:00"),
        ]

        result = backend.glob("**/*.py", "/")

        assert result.error is None
        assert result.matches is not None
        matched_paths = [m["path"] for m in result.matches]
        assert "/app/main.py" in matched_paths
        assert "/app/test.py" in matched_paths
        assert "/app/readme.md" not in matched_paths


# ---------------------------------------------------------------------------
# grep()
# ---------------------------------------------------------------------------


class TestGrep:
    def test_grep_finds_matches(self, backend, mock_cursor):
        """grep() should find literal pattern matches in file content."""
        mock_cursor.fetchall.return_value = [
            ("/app/hello.py", "utf-8", "import os\nprint('hello')\n", "2025-01-01T00:00:00", "2025-01-01T00:00:00"),
            ("/app/other.py", "utf-8", "x = 1\n", "2025-01-01T00:00:00", "2025-01-01T00:00:00"),
        ]

        result = backend.grep("import")

        assert result.error is None
        assert result.matches is not None
        assert len(result.matches) == 1
        assert result.matches[0]["path"] == "/app/hello.py"
        assert result.matches[0]["text"] == "import os"


# ---------------------------------------------------------------------------
# download_files()
# ---------------------------------------------------------------------------


class TestDownloadFiles:
    def test_download_returns_bytes(self, backend, mock_cursor):
        """download_files() should return file content as bytes."""
        mock_cursor.fetchone.return_value = (
            "hello bytes",
            "utf-8",
            "2025-01-01T00:00:00",
            "2025-01-01T00:00:00",
        )

        responses = backend.download_files(["/app/hello.txt"])

        assert len(responses) == 1
        assert responses[0].error is None
        assert responses[0].content == b"hello bytes"

    def test_download_missing_file(self, backend, mock_cursor):
        """download_files() should return file_not_found error for missing files."""
        mock_cursor.fetchone.return_value = None

        responses = backend.download_files(["/app/missing.txt"])

        assert len(responses) == 1
        assert responses[0].error == "file_not_found"
        assert responses[0].content is None


# ---------------------------------------------------------------------------
# upload_files()
# ---------------------------------------------------------------------------


class TestUploadFiles:
    def test_upload_writes_files(self, backend, mock_cursor):
        """upload_files() should delegate to write and return success."""
        responses = backend.upload_files([("/app/data.txt", b"file content")])

        assert len(responses) == 1
        assert responses[0].error is None
        assert responses[0].path == "/app/data.txt"

    def test_upload_binary_files(self, backend, mock_cursor):
        """upload_files() should handle binary content via base64 encoding."""
        binary_data = bytes(range(256))
        responses = backend.upload_files([("/app/image.png", binary_data)])

        assert len(responses) == 1
        assert responses[0].error is None
