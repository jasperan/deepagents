"""Tests for bundled Oracle skills."""

from pathlib import Path

import pytest
import yaml

from deepagents_oracle import DEFAULT_MODEL, SKILLS_DIR

EXPECTED_SKILLS = [
    "oracle-memory-search",
    "oracle-session-resume",
    "oracle-knowledge-ingest",
    "oracle-data-analyst",
    "oracle-multi-agent-memory",
    "research-and-remember",
    "oracle-health-monitor",
]


@pytest.fixture
def skills_path():
    return Path(SKILLS_DIR)


class TestSkillsDirectory:
    def test_skills_dir_exists(self, skills_path):
        assert skills_path.is_dir(), f"Skills directory not found at {skills_path}"

    def test_all_skills_present(self, skills_path):
        actual = sorted(d.name for d in skills_path.iterdir() if d.is_dir())
        assert actual == sorted(EXPECTED_SKILLS)

    @pytest.mark.parametrize("skill_name", EXPECTED_SKILLS)
    def test_skill_has_skill_md(self, skills_path, skill_name):
        skill_file = skills_path / skill_name / "SKILL.md"
        assert skill_file.is_file(), f"Missing {skill_file}"


class TestSkillMetadata:
    @pytest.mark.parametrize("skill_name", EXPECTED_SKILLS)
    def test_frontmatter_parses(self, skills_path, skill_name):
        skill_file = skills_path / skill_name / "SKILL.md"
        content = skill_file.read_text()

        assert content.startswith("---"), "SKILL.md must start with YAML frontmatter"
        _, fm_text, _ = content.split("---", 2)
        metadata = yaml.safe_load(fm_text)

        assert metadata is not None
        assert "name" in metadata
        assert "description" in metadata

    @pytest.mark.parametrize("skill_name", EXPECTED_SKILLS)
    def test_name_matches_directory(self, skills_path, skill_name):
        skill_file = skills_path / skill_name / "SKILL.md"
        content = skill_file.read_text()
        _, fm_text, _ = content.split("---", 2)
        metadata = yaml.safe_load(fm_text)

        assert metadata["name"] == skill_name, (
            f"Frontmatter name '{metadata['name']}' doesn't match directory '{skill_name}'"
        )

    @pytest.mark.parametrize("skill_name", EXPECTED_SKILLS)
    def test_description_not_empty(self, skills_path, skill_name):
        skill_file = skills_path / skill_name / "SKILL.md"
        content = skill_file.read_text()
        _, fm_text, _ = content.split("---", 2)
        metadata = yaml.safe_load(fm_text)

        assert len(metadata["description"]) > 10, "Description too short"

    @pytest.mark.parametrize("skill_name", EXPECTED_SKILLS)
    def test_allowed_tools_present(self, skills_path, skill_name):
        skill_file = skills_path / skill_name / "SKILL.md"
        content = skill_file.read_text()
        _, fm_text, _ = content.split("---", 2)
        metadata = yaml.safe_load(fm_text)

        assert "allowed-tools" in metadata, "Missing allowed-tools in frontmatter"


class TestFactorySkillsIntegration:
    def test_skills_dir_constant(self):
        assert Path(SKILLS_DIR).is_dir()

    def test_default_model_exported(self):
        assert DEFAULT_MODEL == "ollama:qwopus3.5:9b-v3"

    def test_skills_dir_exported(self):
        assert str(Path(SKILLS_DIR)) == SKILLS_DIR
