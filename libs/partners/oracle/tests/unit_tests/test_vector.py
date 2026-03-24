"""Tests for OracleVectorBackend."""

from __future__ import annotations

from contextlib import contextmanager
from unittest.mock import MagicMock, patch

import pytest

from deepagents_oracle.vector import OracleVectorBackend

# ---------------------------------------------------------------------------
# Fixtures
# ---------------------------------------------------------------------------


@pytest.fixture
def mock_embeddings():
    """Return a MagicMock embeddings model that produces deterministic vectors."""
    emb = MagicMock()
    # embed_query returns a list of floats (dimension 1536 in production, 3 here for brevity)
    emb.embed_query.return_value = [0.1, 0.2, 0.3]
    return emb


@pytest.fixture
def mock_cursor():
    """Return a MagicMock cursor that tracks execute/fetchone/fetchall calls."""
    cursor = MagicMock()
    cursor.description = None
    cursor.fetchone.return_value = None
    cursor.fetchall.return_value = []
    return cursor


@pytest.fixture
def mock_conn(mock_cursor):
    """Return a MagicMock connection whose cursor() context manager yields *mock_cursor*."""
    conn = MagicMock()

    @contextmanager
    def _cursor_ctx():
        yield mock_cursor

    conn.cursor = _cursor_ctx
    return conn


@pytest.fixture
def mock_manager(mock_conn):
    """Return a MagicMock OracleConnectionManager whose get_connection yields *mock_conn*."""
    from deepagents_oracle.connection import OracleConnectionManager  # noqa: PLC0415

    manager = MagicMock(spec=OracleConnectionManager)

    @contextmanager
    def _conn_ctx():
        yield mock_conn

    manager.get_connection = _conn_ctx
    return manager


@pytest.fixture
def vector_backend(mock_manager, mock_embeddings):
    """Return an OracleVectorBackend wired to mock manager and embeddings."""
    return OracleVectorBackend(connection_manager=mock_manager, namespace="test_ns", embeddings=mock_embeddings)


# -- write ------------------------------------------------------------------


class TestVectorWrite:
    def test_write_calls_embed_query(self, vector_backend, mock_embeddings, mock_cursor):
        """write() should call embeddings.embed_query() with the content."""
        result = vector_backend.write("/app/hello.txt", "hello world")

        assert result.error is None
        assert result.path == "/app/hello.txt"
        mock_embeddings.embed_query.assert_called_once_with("hello world")

    def test_write_stores_embedding_vector(self, vector_backend, mock_embeddings, mock_cursor):
        """write() should execute SQL that includes the embedding vector."""
        vector_backend.write("/app/hello.txt", "hello world")

        # The MERGE SQL should have been called with embedding parameter
        assert mock_cursor.execute.call_count >= 1
        last_call_args = mock_cursor.execute.call_args
        params = last_call_args[0][1] if len(last_call_args[0]) > 1 else last_call_args[1]
        assert "embedding" in params
        assert params["embedding"] == "[0.1, 0.2, 0.3]"

    def test_write_returns_write_result(self, vector_backend, mock_cursor):
        """write() should return a WriteResult with the correct path."""
        result = vector_backend.write("/app/data.py", "import os")

        assert result.error is None
        assert result.path == "/app/data.py"
        assert result.files_update is None


# -- semantic_grep ----------------------------------------------------------


class TestSemanticGrep:
    def test_semantic_grep_embeds_query(self, vector_backend, mock_embeddings, mock_cursor):
        """semantic_grep() should embed the query string."""
        mock_cursor.fetchall.return_value = []
        vector_backend.semantic_grep("find similar content")

        mock_embeddings.embed_query.assert_called_once_with("find similar content")

    def test_semantic_grep_executes_vector_distance_sql(self, vector_backend, mock_embeddings, mock_cursor):
        """semantic_grep() should use VECTOR_DISTANCE in the SQL query."""
        mock_cursor.fetchall.return_value = []
        vector_backend.semantic_grep("find similar")

        assert mock_cursor.execute.call_count >= 1
        sql = mock_cursor.execute.call_args[0][0]
        assert "VECTOR_DISTANCE" in sql
        assert "COSINE" in sql

    def test_semantic_grep_returns_grep_result(self, vector_backend, mock_embeddings, mock_cursor):
        """semantic_grep() should return a GrepResult with matches."""
        mock_cursor.fetchall.return_value = [
            ("/app/hello.txt", "hello world\nsecond line", 0.05),
            ("/app/data.py", "import os\nimport sys", 0.15),
        ]

        result = vector_backend.semantic_grep("hello")

        assert result.error is None
        assert result.matches is not None
        assert len(result.matches) == 2
        assert result.matches[0]["path"] == "/app/hello.txt"

    def test_semantic_grep_with_path_filter(self, vector_backend, mock_embeddings, mock_cursor):
        """semantic_grep() with path should restrict results to that prefix."""
        mock_cursor.fetchall.return_value = [
            ("/app/src/hello.txt", "hello world", 0.05),
        ]

        vector_backend.semantic_grep("hello", path="/app/src")

        assert mock_cursor.execute.call_count >= 1
        sql = mock_cursor.execute.call_args[0][0]
        params = mock_cursor.execute.call_args[0][1]
        assert "file_path LIKE" in sql
        assert params["path_prefix"] == "/app/src%"

    def test_semantic_grep_top_k(self, vector_backend, mock_embeddings, mock_cursor):
        """semantic_grep() should respect the top_k parameter."""
        mock_cursor.fetchall.return_value = []
        vector_backend.semantic_grep("query", top_k=5)

        params = mock_cursor.execute.call_args[0][1]
        assert params["top_k"] == 5

    def test_semantic_grep_default_top_k(self, vector_backend, mock_embeddings, mock_cursor):
        """semantic_grep() should default to top_k=10."""
        mock_cursor.fetchall.return_value = []
        vector_backend.semantic_grep("query")

        params = mock_cursor.execute.call_args[0][1]
        assert params["top_k"] == 10


# -- _ensure_initialized (via init_schema) ----------------------------------


class TestEnsureInitialized:
    @patch("deepagents_oracle.vector.init_schema")
    def test_ensure_initialized_uses_vector_true(self, mock_init_schema, mock_manager, mock_embeddings, mock_conn):
        """_ensure_initialized() should call init_schema with vector=True."""
        backend = OracleVectorBackend(connection_manager=mock_manager, namespace="test_ns", embeddings=mock_embeddings)
        backend._ensure_initialized()

        mock_init_schema.assert_called_once_with(mock_conn, vector=True)

    @patch("deepagents_oracle.vector.init_schema")
    def test_ensure_initialized_only_runs_once(self, mock_init_schema, mock_manager, mock_embeddings, mock_conn):
        """_ensure_initialized() should be idempotent, only running schema init once."""
        backend = OracleVectorBackend(connection_manager=mock_manager, namespace="test_ns", embeddings=mock_embeddings)
        backend._ensure_initialized()
        backend._ensure_initialized()

        mock_init_schema.assert_called_once()


# -- inherited grep still works ---------------------------------------------


class TestInheritedGrep:
    def test_grep_still_works(self, vector_backend, mock_cursor):
        """Regular grep() should still work via parent class."""
        mock_cursor.fetchall.return_value = [
            ("/app/hello.py", "utf-8", "import os\nprint('hello')\n", "2025-01-01T00:00:00", "2025-01-01T00:00:00"),
        ]

        result = vector_backend.grep("import")

        assert result.error is None
        assert result.matches is not None
        assert len(result.matches) == 1
        assert result.matches[0]["path"] == "/app/hello.py"
