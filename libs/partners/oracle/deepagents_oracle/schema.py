"""Oracle schema DDL and auto-initialization."""

from __future__ import annotations

import logging

import oracledb

logger = logging.getLogger(__name__)

DA_FILES_DDL = """
CREATE TABLE da_files (
    namespace   VARCHAR2(512)  NOT NULL,
    file_path   VARCHAR2(1024) NOT NULL,
    content     CLOB,
    encoding    VARCHAR2(10)   DEFAULT 'utf-8',
    created_at  TIMESTAMP      DEFAULT SYSTIMESTAMP,
    modified_at TIMESTAMP      DEFAULT SYSTIMESTAMP,
    CONSTRAINT da_files_pk PRIMARY KEY (namespace, file_path)
)
"""

VECTOR_DDL = """
ALTER TABLE da_files ADD (
    embedding VECTOR(1536, FLOAT64)
)
"""

VECTOR_INDEX_DDL = """
CREATE VECTOR INDEX da_files_vec_idx ON da_files(embedding)
    ORGANIZATION NEIGHBOR PARTITIONS
    DISTANCE COSINE
"""


def _execute_ignoring_exists(cursor: oracledb.Cursor, sql: str) -> None:
    """Execute DDL, silently skipping ORA-00955 (already exists) and ORA-01430 (column already added)."""
    try:
        cursor.execute(sql)
    except oracledb.DatabaseError as exc:
        error_msg = str(exc)
        if "ORA-00955" in error_msg or "ORA-01430" in error_msg:
            logger.debug("Schema object already exists, skipping: %s", error_msg)
        else:
            raise


def init_schema(conn: oracledb.Connection, *, vector: bool = False) -> None:
    """Create the da_files table (and optionally vector column + index).

    Safe to call multiple times; existing objects are silently skipped.

    Args:
        conn: An active Oracle connection.
        vector: If True, also add the embedding VECTOR column and index.
    """
    with conn.cursor() as cur:
        _execute_ignoring_exists(cur, DA_FILES_DDL)
        logger.info("da_files table ready")

        if vector:
            _execute_ignoring_exists(cur, VECTOR_DDL)
            _execute_ignoring_exists(cur, VECTOR_INDEX_DDL)
            logger.info("Vector column and index ready")

    conn.commit()
