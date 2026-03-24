.PHONY: format lint type typecheck test tests integration_test integration_tests test_watch help
.DEFAULT_GOAL := help
.EXPORT_ALL_VARIABLES:
UV_FROZEN = true

TEST_FILE ?= tests/unit_tests/
PYTEST_EXTRA ?=
PYTHON_FILES=.
integration_test integration_tests: TEST_FILE=tests/integration_tests/

test tests:
	uv run --group test pytest -vvv $(PYTEST_EXTRA) --disable-socket --allow-unix-socket $(TEST_FILE)

integration_test integration_tests:
	uv run --group test pytest -vvv --timeout 30 $(TEST_FILE)

test_watch:
	uv run --group test ptw --now . -- -vv $(TEST_FILE)

lint lint_diff:
	[ "$(PYTHON_FILES)" = "" ] || uv run --all-groups ruff check $(PYTHON_FILES)
	[ "$(PYTHON_FILES)" = "" ] || uv run --all-groups ruff format $(PYTHON_FILES) --diff
	$(MAKE) type

type typecheck:
	uv run --all-groups ty check deepagents_oracle

format format_diff:
	[ "$(PYTHON_FILES)" = "" ] || uv run --all-groups ruff format $(PYTHON_FILES)
	[ "$(PYTHON_FILES)" = "" ] || uv run --all-groups ruff check --fix $(PYTHON_FILES)

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
