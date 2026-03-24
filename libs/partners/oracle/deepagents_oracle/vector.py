"""OracleVectorBackend: extends OracleStoreBackend with vector embedding support.

On write, generates an embedding via a LangChain Embeddings model and stores
it in the ``embedding`` VECTOR column.  Provides ``semantic_grep()`` for
similarity search using ``VECTOR_DISTANCE(COSINE)``.
"""

from __future__ import annotations

import logging
from typing import TYPE_CHECKING

from deepagents.backends.protocol import GrepMatch, GrepResult

from deepagents_oracle.backend import OracleStoreBackend, WriteResult, _read_clob, create_file_data
from deepagents_oracle.schema import init_schema

if TYPE_CHECKING:
    from langchain_core.embeddings import Embeddings

    from deepagents_oracle.connection import OracleConnectionManager

logger = logging.getLogger(__name__)

# -- SQL statements ----------------------------------------------------------

_MERGE_WITH_EMBEDDING_SQL = """
MERGE INTO da_files tgt
USING (SELECT :namespace AS namespace, :file_path AS file_path FROM DUAL) src
ON (tgt.namespace = src.namespace AND tgt.file_path = src.file_path)
WHEN MATCHED THEN
    UPDATE SET content     = :content,
               encoding    = :encoding,
               embedding   = :embedding,
               modified_at = SYSTIMESTAMP
WHEN NOT MATCHED THEN
    INSERT (namespace, file_path, content, encoding, embedding, created_at, modified_at)
    VALUES (:namespace, :file_path, :content, :encoding, :embedding, SYSTIMESTAMP, SYSTIMESTAMP)
"""

_SEMANTIC_SEARCH_SQL = """
SELECT file_path, content, VECTOR_DISTANCE(embedding, :qvec, COSINE) AS dist
FROM da_files
WHERE namespace = :namespace
ORDER BY dist ASC
FETCH FIRST :top_k ROWS ONLY
"""

_SEMANTIC_SEARCH_WITH_PATH_SQL = """
SELECT file_path, content, VECTOR_DISTANCE(embedding, :qvec, COSINE) AS dist
FROM da_files
WHERE namespace = :namespace AND file_path LIKE :path_prefix
ORDER BY dist ASC
FETCH FIRST :top_k ROWS ONLY
"""


class OracleVectorBackend(OracleStoreBackend):
    """Backend that extends ``OracleStoreBackend`` with vector embedding support.

    On every ``write()``, the content is embedded via the provided LangChain
    ``Embeddings`` model and stored in the ``embedding`` VECTOR column.  The
    ``semantic_grep()`` method embeds a query and ranks stored files by cosine
    similarity using Oracle's ``VECTOR_DISTANCE`` function.

    Args:
        connection_manager: An already-connected ``OracleConnectionManager``.
        namespace: Namespace string used to partition rows in ``da_files``.
        embeddings: A LangChain ``Embeddings`` model for generating vectors.
    """

    def __init__(
        self,
        connection_manager: OracleConnectionManager,
        namespace: str = "default",
        *,
        embeddings: Embeddings,
    ) -> None:
        """Initialize with a connection manager, namespace, and embeddings model."""
        super().__init__(connection_manager=connection_manager, namespace=namespace)
        self._embeddings = embeddings
        self._initialized = False

    # -- schema initialization -----------------------------------------------

    def _ensure_initialized(self) -> None:
        """Create the da_files table with the vector column and index.

        Safe to call multiple times; sets a flag after the first successful run.
        """
        if self._initialized:
            return
        with self._cm.get_connection() as conn:
            init_schema(conn, vector=True)
        self._initialized = True

    # -- overridden write ----------------------------------------------------

    def write(self, file_path: str, content: str) -> WriteResult:
        """MERGE (upsert) a file into ``da_files`` with its embedding vector.

        Generates an embedding from *content* via the configured Embeddings
        model and stores it alongside the file data.
        """
        file_data = create_file_data(content)
        embedding_vector = self._embeddings.embed_query(content)
        embedding_str = "[" + ", ".join(str(v) for v in embedding_vector) + "]"

        with self._cm.get_connection() as conn:
            with conn.cursor() as cur:
                cur.execute(
                    _MERGE_WITH_EMBEDDING_SQL,
                    {
                        "namespace": self._namespace,
                        "file_path": file_path,
                        "content": file_data["content"],
                        "encoding": file_data["encoding"],
                        "embedding": embedding_str,
                    },
                )
            conn.commit()
        return WriteResult(path=file_path, files_update=None)

    # -- semantic search -----------------------------------------------------

    def semantic_grep(
        self,
        query: str,
        path: str | None = None,
        *,
        top_k: int = 10,
    ) -> GrepResult:
        """Search files by semantic similarity using ``VECTOR_DISTANCE(COSINE)``.

        Embeds *query* via the configured Embeddings model and ranks stored
        files by cosine distance to the query vector.

        Args:
            query: Natural-language search string to embed.
            path: Optional path prefix to restrict the search scope.
            top_k: Maximum number of results to return (default 10).

        Returns:
            A ``GrepResult`` whose matches are ordered by similarity (closest first).
            Each match contains the file path and the first line of content as text.
        """
        query_vector = self._embeddings.embed_query(query)
        qvec_str = "[" + ", ".join(str(v) for v in query_vector) + "]"

        if path is not None:
            sql = _SEMANTIC_SEARCH_WITH_PATH_SQL
            params: dict = {
                "namespace": self._namespace,
                "qvec": qvec_str,
                "path_prefix": path + "%",
                "top_k": top_k,
            }
        else:
            sql = _SEMANTIC_SEARCH_SQL
            params = {
                "namespace": self._namespace,
                "qvec": qvec_str,
                "top_k": top_k,
            }

        matches: list[GrepMatch] = []
        with self._cm.get_connection() as conn, conn.cursor() as cur:
            cur.execute(sql, params)
            for row in cur.fetchall():
                fp, content_raw, _dist = row
                content = _read_clob(content_raw)
                first_line = content.split("\n")[0] if content else ""
                matches.append(GrepMatch(path=fp, line=1, text=first_line))

        return GrepResult(matches=matches)
