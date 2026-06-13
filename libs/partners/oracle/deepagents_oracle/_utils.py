"""Shared low-level helpers for the Oracle backend modules."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any

if TYPE_CHECKING:
    from collections.abc import Iterable


def read_clob(value: Any) -> str:  # noqa: ANN401
    """Read a CLOB/LOB value, handling both str and LOB objects."""
    if isinstance(value, str):
        return value
    # oracledb LOB objects expose a .read() method
    if hasattr(value, "read"):
        return value.read()
    return str(value)


def vector_literal(vec: Iterable[float]) -> str:
    """Serialize a numeric vector to Oracle's ``[v0, v1, ...]`` VECTOR literal."""
    return "[" + ", ".join(str(v) for v in vec) + "]"
