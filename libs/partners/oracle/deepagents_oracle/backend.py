"""OracleStoreBackend: BackendProtocol implementation backed by Oracle Database.

All file operations translate to SQL against the ``da_files`` table, keyed
by ``(namespace, file_path)``.  Uses shared utility functions from
``deepagents.backends.utils`` wherever possible to stay consistent with the
reference ``StoreBackend``.
"""

from __future__ import annotations

import base64
import logging
from typing import Any

from deepagents.backends.protocol import (
    BackendProtocol,
    EditResult,
    FileData,
    FileDownloadResponse,
    FileInfo,
    FileUploadResponse,
    GlobResult,
    GrepResult,
    LsResult,
    ReadResult,
    WriteResult,
)
from deepagents.backends.utils import (
    _get_file_type,
    _glob_search_files,
    create_file_data,
    file_data_to_string,
    grep_matches_from_files,
    perform_string_replacement,
    slice_read_response,
    update_file_data,
)

from deepagents_oracle.connection import OracleConnectionManager

logger = logging.getLogger(__name__)

# -- SQL statements ----------------------------------------------------------

_MERGE_SQL = """
MERGE INTO da_files tgt
USING (SELECT :namespace AS namespace, :file_path AS file_path FROM DUAL) src
ON (tgt.namespace = src.namespace AND tgt.file_path = src.file_path)
WHEN MATCHED THEN
    UPDATE SET content     = :content,
               encoding    = :encoding,
               modified_at = SYSTIMESTAMP
WHEN NOT MATCHED THEN
    INSERT (namespace, file_path, content, encoding, created_at, modified_at)
    VALUES (:namespace, :file_path, :content, :encoding, SYSTIMESTAMP, SYSTIMESTAMP)
"""

_SELECT_ONE_SQL = """
SELECT content, encoding, TO_CHAR(created_at, 'YYYY-MM-DD"T"HH24:MI:SS'), TO_CHAR(modified_at, 'YYYY-MM-DD"T"HH24:MI:SS')
FROM da_files
WHERE namespace = :namespace AND file_path = :file_path
"""

_SELECT_ALL_SQL = """
SELECT file_path, encoding, content, TO_CHAR(created_at, 'YYYY-MM-DD"T"HH24:MI:SS'), TO_CHAR(modified_at, 'YYYY-MM-DD"T"HH24:MI:SS')
FROM da_files
WHERE namespace = :namespace
"""

_UPDATE_CONTENT_SQL = """
UPDATE da_files
SET content = :content, encoding = :encoding, modified_at = SYSTIMESTAMP
WHERE namespace = :namespace AND file_path = :file_path
"""


def _read_clob(value: Any) -> str:
    """Read a CLOB/LOB value, handling both str and LOB objects."""
    if isinstance(value, str):
        return value
    # oracledb LOB objects expose a .read() method
    if hasattr(value, "read"):
        return value.read()
    return str(value)


class OracleStoreBackend(BackendProtocol):
    """Backend that stores files in an Oracle Database ``da_files`` table.

    Every file is identified by ``(namespace, file_path)``.  The namespace
    isolates different agents or tenants within the same table.

    Args:
        connection_manager: An already-connected ``OracleConnectionManager``.
        namespace: Namespace string used to partition rows in ``da_files``.
    """

    def __init__(
        self,
        connection_manager: OracleConnectionManager,
        namespace: str = "default",
    ) -> None:
        self._cm = connection_manager
        self._namespace = namespace

    # -- internal helpers ----------------------------------------------------

    def _get_all_files(self) -> dict[str, FileData]:
        """Load every file in the namespace as ``{path: FileData}``."""
        files: dict[str, FileData] = {}
        with self._cm.get_connection() as conn:
            with conn.cursor() as cur:
                cur.execute(_SELECT_ALL_SQL, {"namespace": self._namespace})
                for row in cur.fetchall():
                    file_path, encoding, content_raw, created_at, modified_at = row
                    files[file_path] = FileData(
                        content=_read_clob(content_raw),
                        encoding=encoding or "utf-8",
                        created_at=created_at or "",
                        modified_at=modified_at or "",
                    )
        return files

    def _fetch_one(self, file_path: str) -> FileData | None:
        """Fetch a single file's data, or ``None`` if missing."""
        with self._cm.get_connection() as conn:
            with conn.cursor() as cur:
                cur.execute(_SELECT_ONE_SQL, {"namespace": self._namespace, "file_path": file_path})
                row = cur.fetchone()
                if row is None:
                    return None
                content_raw, encoding, created_at, modified_at = row
                return FileData(
                    content=_read_clob(content_raw),
                    encoding=encoding or "utf-8",
                    created_at=created_at or "",
                    modified_at=modified_at or "",
                )

    # -- BackendProtocol methods ---------------------------------------------

    def write(self, file_path: str, content: str) -> WriteResult:
        """MERGE (upsert) a file into ``da_files``.

        Unlike the reference ``StoreBackend`` which rejects writes to existing
        files, Oracle backend uses MERGE for idempotent upserts.
        """
        file_data = create_file_data(content)
        with self._cm.get_connection() as conn:
            with conn.cursor() as cur:
                cur.execute(
                    _MERGE_SQL,
                    {
                        "namespace": self._namespace,
                        "file_path": file_path,
                        "content": file_data["content"],
                        "encoding": file_data["encoding"],
                    },
                )
            conn.commit()
        return WriteResult(path=file_path, files_update=None)

    def read(self, file_path: str, offset: int = 0, limit: int = 2000) -> ReadResult:
        """SELECT content from ``da_files``, then slice with ``slice_read_response()``."""
        with self._cm.get_connection() as conn:
            with conn.cursor() as cur:
                cur.execute(_SELECT_ONE_SQL, {"namespace": self._namespace, "file_path": file_path})
                row = cur.fetchone()

        if row is None:
            return ReadResult(error=f"File '{file_path}' not found")

        content_raw, encoding, created_at, modified_at = row
        file_data = FileData(
            content=_read_clob(content_raw),
            encoding=encoding or "utf-8",
            created_at=created_at or "",
            modified_at=modified_at or "",
        )

        if _get_file_type(file_path) != "text":
            return ReadResult(file_data=file_data)

        sliced = slice_read_response(file_data, offset, limit)
        if isinstance(sliced, ReadResult):
            return sliced
        return ReadResult(
            file_data=FileData(
                content=sliced,
                encoding=file_data.get("encoding", "utf-8"),
                created_at=file_data.get("created_at", ""),
                modified_at=file_data.get("modified_at", ""),
            )
        )

    def edit(self, file_path: str, old_string: str, new_string: str, replace_all: bool = False) -> EditResult:  # noqa: FBT001, FBT002
        """SELECT, apply ``perform_string_replacement()``, then UPDATE."""
        with self._cm.get_connection() as conn:
            with conn.cursor() as cur:
                cur.execute(_SELECT_ONE_SQL, {"namespace": self._namespace, "file_path": file_path})
                row = cur.fetchone()

                if row is None:
                    return EditResult(error=f"Error: File '{file_path}' not found")

                content_raw, encoding, created_at, modified_at = row
                file_data = FileData(
                    content=_read_clob(content_raw),
                    encoding=encoding or "utf-8",
                    created_at=created_at or "",
                    modified_at=modified_at or "",
                )

                content = file_data_to_string(file_data)
                result = perform_string_replacement(content, old_string, new_string, replace_all)

                if isinstance(result, str):
                    return EditResult(error=result)

                new_content, occurrences = result
                new_file_data = update_file_data(file_data, new_content)

                cur.execute(
                    _UPDATE_CONTENT_SQL,
                    {
                        "namespace": self._namespace,
                        "file_path": file_path,
                        "content": new_file_data["content"],
                        "encoding": new_file_data["encoding"],
                    },
                )
            conn.commit()

        return EditResult(path=file_path, files_update=None, occurrences=int(occurrences))

    def ls(self, path: str) -> LsResult:
        """List files and directories under *path* (non-recursive)."""
        files = self._get_all_files()
        infos: list[FileInfo] = []
        subdirs: set[str] = set()

        normalized_path = path if path.endswith("/") else path + "/"

        for fp, fd in files.items():
            if not fp.startswith(normalized_path):
                continue

            relative = fp[len(normalized_path):]

            if "/" in relative:
                subdir_name = relative.split("/")[0]
                subdirs.add(normalized_path + subdir_name + "/")
                continue

            size = len(fd.get("content", ""))
            infos.append(
                {
                    "path": fp,
                    "is_dir": False,
                    "size": int(size),
                    "modified_at": fd.get("modified_at", ""),
                }
            )

        infos.extend(FileInfo(path=subdir, is_dir=True, size=0, modified_at="") for subdir in sorted(subdirs))
        infos.sort(key=lambda x: x.get("path", ""))
        return LsResult(entries=infos)

    def glob(self, pattern: str, path: str = "/") -> GlobResult:
        """Find files matching a glob pattern in Oracle storage."""
        files = self._get_all_files()
        result = _glob_search_files(files, pattern, path)

        if result == "No files found":
            return GlobResult(matches=[])

        paths = result.split("\n")
        infos: list[FileInfo] = []
        for p in paths:
            fd = files.get(p)
            if fd:
                size = len(fd.get("content", ""))
            else:
                size = 0
            infos.append(
                {
                    "path": p,
                    "is_dir": False,
                    "size": int(size),
                    "modified_at": fd.get("modified_at", "") if fd else "",
                }
            )
        return GlobResult(matches=infos)

    def grep(self, pattern: str, path: str | None = None, glob: str | None = None) -> GrepResult:
        """Search stored files for a literal text pattern."""
        files = self._get_all_files()
        return grep_matches_from_files(files, pattern, path, glob)

    def upload_files(self, files: list[tuple[str, bytes]]) -> list[FileUploadResponse]:
        """Upload multiple files, encoding binary content as base64."""
        responses: list[FileUploadResponse] = []

        for fpath, raw_bytes in files:
            try:
                content_str = raw_bytes.decode("utf-8")
                encoding = "utf-8"
            except UnicodeDecodeError:
                content_str = base64.standard_b64encode(raw_bytes).decode("ascii")
                encoding = "base64"

            file_data = create_file_data(content_str, encoding=encoding)
            with self._cm.get_connection() as conn:
                with conn.cursor() as cur:
                    cur.execute(
                        _MERGE_SQL,
                        {
                            "namespace": self._namespace,
                            "file_path": fpath,
                            "content": file_data["content"],
                            "encoding": file_data["encoding"],
                        },
                    )
                conn.commit()
            responses.append(FileUploadResponse(path=fpath, error=None))

        return responses

    def download_files(self, paths: list[str]) -> list[FileDownloadResponse]:
        """Download multiple files as bytes."""
        responses: list[FileDownloadResponse] = []

        for fpath in paths:
            file_data = self._fetch_one(fpath)

            if file_data is None:
                responses.append(FileDownloadResponse(path=fpath, content=None, error="file_not_found"))
                continue

            content_str = file_data_to_string(file_data)
            encoding = file_data["encoding"]
            content_bytes = base64.standard_b64decode(content_str) if encoding == "base64" else content_str.encode("utf-8")
            responses.append(FileDownloadResponse(path=fpath, content=content_bytes, error=None))

        return responses
