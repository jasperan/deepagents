"""Tests for Oracle schema initialization."""

from unittest.mock import MagicMock

import oracledb
import pytest

from deepagents_oracle.schema import DA_FILES_DDL, VECTOR_DDL, VECTOR_INDEX_DDL, init_schema


class TestSchemaConstants:
    def test_da_files_ddl_creates_table(self):
        assert "CREATE TABLE" in DA_FILES_DDL
        assert "da_files" in DA_FILES_DDL
        assert "namespace" in DA_FILES_DDL
        assert "file_path" in DA_FILES_DDL
        assert "content" in DA_FILES_DDL
        assert "CLOB" in DA_FILES_DDL

    def test_vector_ddl_adds_column(self):
        assert "VECTOR" in VECTOR_DDL
        assert "da_files" in VECTOR_DDL

    def test_vector_index_ddl(self):
        assert "VECTOR INDEX" in VECTOR_INDEX_DDL
        assert "COSINE" in VECTOR_INDEX_DDL


class TestInitSchema:
    def _make_mock_conn(self):
        mock_conn = MagicMock()
        mock_cursor = MagicMock()
        mock_conn.cursor.return_value.__enter__ = lambda _s: mock_cursor
        mock_conn.cursor.return_value.__exit__ = MagicMock(return_value=False)
        return mock_conn, mock_cursor

    def test_init_creates_table(self):
        mock_conn, mock_cursor = self._make_mock_conn()
        init_schema(mock_conn, vector=False)
        mock_cursor.execute.assert_called()
        executed_sql = mock_cursor.execute.call_args_list[0][0][0]
        assert "da_files" in executed_sql

    def test_init_with_vector(self):
        mock_conn, mock_cursor = self._make_mock_conn()
        init_schema(mock_conn, vector=True)
        executed_calls = [c[0][0] for c in mock_cursor.execute.call_args_list]
        has_vector = any("VECTOR" in sql for sql in executed_calls)
        assert has_vector

    def test_init_without_vector_skips_vector_ddl(self):
        mock_conn, mock_cursor = self._make_mock_conn()
        init_schema(mock_conn, vector=False)
        executed_calls = [c[0][0] for c in mock_cursor.execute.call_args_list]
        has_vector = any("VECTOR" in sql for sql in executed_calls)
        assert not has_vector

    def test_init_ignores_already_exists_error(self):
        """ORA-00955: name is already used by an existing object."""
        mock_conn, mock_cursor = self._make_mock_conn()
        error = oracledb.DatabaseError("ORA-00955: name is already used")
        mock_cursor.execute.side_effect = error
        # Should not raise
        init_schema(mock_conn, vector=False)

    def test_init_ignores_column_already_added(self):
        """ORA-01430: column being added already exists."""
        mock_conn, mock_cursor = self._make_mock_conn()
        error = oracledb.DatabaseError("ORA-01430: column being added already exists")
        mock_cursor.execute.side_effect = error
        # Should not raise
        init_schema(mock_conn, vector=True)

    def test_init_raises_on_unexpected_error(self):
        mock_conn, mock_cursor = self._make_mock_conn()
        error = oracledb.DatabaseError("ORA-01017: invalid username/password")
        mock_cursor.execute.side_effect = error
        with pytest.raises(oracledb.DatabaseError, match="ORA-01017"):
            init_schema(mock_conn, vector=False)

    def test_init_commits(self):
        mock_conn, _mock_cursor = self._make_mock_conn()
        init_schema(mock_conn, vector=False)
        mock_conn.commit.assert_called_once()
