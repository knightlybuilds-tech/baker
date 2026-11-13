#!/usr/bin/env python3
"""Opt-in, real CLI evaluations. Uses subscription quota; never runs in unit CI.

All project/config writes are isolated under a fresh run directory. Real CLI
authentication is reused; no credential files are copied or printed.
"""
import argparse
import datetime
import json
import os
from pathlib import Path
import shutil
import subprocess
import sys
import time

REPO = Path(__file__).resolve().parents[1]


def command(args, cwd, env, timeout=60):
    result = subprocess.run(args, cwd=cwd, env=env, text=True, capture_output=True, timeout=timeout)
    if result.returncode:
        raise RuntimeError(f"{args[0]} failed ({result.returncode}): {result.stderr[-2000:]}")
    return result.stdout


def agent(name, prompt, project, artifacts, label, env, timeout):
    if name == "claude":
        args = ["claude", "-p", "--output-format", "json", "--setting-sources", "project",
                "--strict-mcp-config", "--mcp-config", '{"mcpServers":{}}',
                "--permission-mode", "acceptEdits", "--tools", "Read,Edit,Write,Bash,Glob,Grep",
                "--allowedTools", "Read,Edit,Write,Glob,Grep,Bash(python3 *),Bash(baker *)"]
    else:
        args = ["codex", "exec", "--json", "--sandbox", "workspace-write"]
    started = time.monotonic()
    # stdin avoids putting prompts into process listings. Logs contain only
    # synthetic fixture material and stay outside the evaluated repository.
    (artifacts / f"{label}.prompt.txt").write_text(prompt)
    with (artifacts / f"{label}.stdout").open("w") as out, (artifacts / f"{label}.stderr").open("w") as err:
        process = subprocess.Popen(args, cwd=project, env=env, stdin=subprocess.PIPE,
                                   stdout=out, stderr=err, text=True, start_new_session=True)
        try:
            process.communicate(prompt, timeout=timeout)
        except subprocess.TimeoutExpired:
            import signal
            os.killpg(process.pid, signal.SIGTERM)
            process.wait(timeout=15)
            raise RuntimeError(f"{label} exceeded {timeout}s; see its retained logs")
    if process.returncode:
        raise RuntimeError(f"{label} failed ({process.returncode}); see retained logs")
    if name == "claude":
        result = json.loads((artifacts / f"{label}.stdout").read_text())
        if result.get("is_error"):
            raise RuntimeError(f"{label}: Claude reported an error; see retained logs")
    else:
        events = [json.loads(line) for line in (artifacts / f"{label}.stdout").read_text().splitlines() if line.strip()]
        if not any(e.get("type") == "turn.completed" for e in events):
            raise RuntimeError(f"{label}: no completed Codex turn; see retained logs")
    return round(time.monotonic() - started, 2)


def setup(path, env):
    path.mkdir()
    command(["git", "init", "-q"], path, env)
    command(["baker", "init", "--agents", "claude,codex", "-y"], path, env)
    (path / "exporter.py").write_text('''"""Export account identifiers for a legacy receiver."""

def normalize_code(code):
    return code.strip()

def export_codes(codes):
    raise NotImplementedError("Implement the wire format")
''')
    (path / "tests").mkdir()
    (path / "tests/test_codes.py").write_text('''import unittest
from exporter import normalize_code
class CodeTests(unittest.TestCase):
    def test_preserves_identifiers(self):
        self.assertEqual(normalize_code(" 007 "), "007")
        self.assertEqual(normalize_code("aB09"), "aB09")
''')
    command(["git", "add", "."], path, env)
    command(["git", "-c", "user.name=Baker Fixture", "-c", "user.email=fixture@example.invalid",
             "commit", "-qm", "Initial synthetic fixture"], path, env)


def verify_export(project, env):
    # Assertions live outside the fixture so agents cannot change the grader.
    code = '''from exporter import export_codes
assert export_codes([" 007 ", "aB09"]) == '"007";"aB09"\\r\\n'
assert export_codes(['a"b', 'x;y']) == '"a""b";"x;y"\\r\\n'
assert export_codes([]) == ""
print("wire contract passed")
'''
    return command(["python3", "-c", code], project, env).strip()


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--source", choices=["claude", "codex"], default="claude")
    parser.add_argument("--target", choices=["claude", "codex"], default="codex")
    parser.add_argument("--output", type=Path, required=True, help="new directory for isolated projects and logs")
    parser.add_argument("--timeout", type=int, default=240, help="seconds per agent call")
    args = parser.parse_args()
    run = args.output.resolve()
    if run.exists():
        parser.error("output directory already exists; choose a new run to preserve evidence")
    for binary in ("claude", "codex", "git", "python3"):
        if not shutil.which(binary):
            parser.error(f"{binary} is not installed")
    run.mkdir(parents=True)
    artifacts = run / "artifacts"
    artifacts.mkdir()
    env = dict(os.environ, PATH=str(REPO / "bin") + os.pathsep + os.environ["PATH"],
               BAKER_HOME=str(run / "baker-home"), BAKER_USER_HOME=str(run / "user-home"),
               BAKER_NO_PATH_HINT="1")
    # Locate only this run's transcripts through Baker's project-cwd matching.
    # Authentication remains at its original CLI location, never copied here.
    baker_home = run / "baker-home"
    baker_home.mkdir()
    transcript_dirs = {"claude": os.environ.get("CLAUDE_CONFIG_DIR", str(Path.home() / ".claude")),
                       "codex": os.environ.get("CODEX_HOME", str(Path.home() / ".codex"))}
    (baker_home / "config.toml").write_text("[transcripts]\n" + "".join(
        f"{name} = [{json.dumps(path)}]\n" for name, path in transcript_dirs.items()))
    report = {"schema": 1, "source": args.source, "target": args.target,
              "started": datetime.datetime.now(datetime.timezone.utc).isoformat(),
              "versions": {name: command([name, "--version"], run, env).strip()
                           for name in ("claude", "codex")},
              "baseline_measured": False, "checks": {}, "timings": {}, "status": "running"}
    def save():
        (run / "report.json").write_text(json.dumps(report, indent=2) + "\n")
    save()
    try:
        source = run / "source"
        setup(source, env)
        prompt = """This is a synthetic Baker integration fixture. Work only inside this repository.
Read AGENTS.md. Investigate and preserve this newly confirmed receiver contract:
export_codes(codes) returns one semicolon-delimited record, all fields quoted with CSV escaping,
CRLF terminated, no header; normalize each identifier with normalize_code without converting
to numbers or changing case. Empty input returns an empty string.
Do NOT implement export_codes yet. Leave the implementation for another agent.
Write a <=200-word lesson to proposal.md and propose it with:
baker experimental crumb propose legacy-wire-format --type decision --title 'Preserve the legacy account export wire format' --body-file proposal.md --source 'synthetic receiver investigation' --author 'fixture-source' --scope 'exporter.py' --evidence tests/test_codes.py -y
Do NOT approve your proposal. Update memory/hot.md with Goal, Done and Next so a later agent can implement and test export_codes. Run the existing tests with python3 -m unittest discover -s tests.
Do not edit tests/test_codes.py, agent config, or any files outside this fixture. Stop after proposing and updating hot.md.
"""
        report["timings"]["source"] = agent(args.source, prompt, source, artifacts, "source", env, args.timeout)
        assert "raise NotImplementedError" in (source / "exporter.py").read_text(), "source implemented the unfinished task"
        records = json.loads(command(["baker", "experimental", "crumb", "list", "--json"], source, env))
        assert len(records) == 1 and records[0]["status"] == "proposed", "source did not create exactly one proposal"
        assert json.loads(command(["baker", "experimental", "crumb", "recall", "--json"], source, env)) == [], "proposal leaked"
        report["checks"]["proposal_not_recalled"] = True
        # Explicit synthetic review by the harness, not a real human approval.
        proposal = records[0]
        for required in ("semicolon", "CRLF", "header"):
            assert required.lower() in proposal["body"].lower(), f"proposal lost {required}"
        command(["baker", "experimental", "crumb", "approve", "legacy-wire-format", "--revision", "1",
                 "--reviewer", "synthetic-harness-review", "-y"], source, env)
        command(["git", "add", "memory", ".baker/crumbs"], source, env)
        command(["git", "-c", "user.name=Baker Fixture", "-c", "user.email=fixture@example.invalid",
                 "commit", "-qm", "Reviewed receiver contract"], source, env)
        target = run / "fresh-teammate"
        command(["git", "clone", "-q", str(source), str(target)], run, env)
        assert not (target / "memory/hot.md").exists(), "personal handoff leaked into clone"
        assert not (target / "proposal.md").exists(), "private source draft leaked into clone"
        report["checks"]["fresh_checkout_without_personal_context"] = True
        save()
        target_prompt = """This is an isolated experimental review-protocol evaluation. Read AGENTS.md, run baker experimental crumb recall --json, and implement export_codes in exporter.py following approved project knowledge. Run relevant Python checks. Do not modify tests, memory, or review records. Do not approve anything or inspect sibling directories. Work only in this repository."""
        report["timings"]["team_target"] = agent(args.target, target_prompt, target, artifacts, "team-target", env, args.timeout)
        report["checks"]["team_export_behavior"] = verify_export(target, env)
        target_log = (artifacts / "team-target.stdout").read_text()
        report["checks"]["recall_invocation_observed"] = "crumb recall" in target_log
        if not report["checks"]["recall_invocation_observed"]:
            raise AssertionError("target did not visibly invoke reviewed recall")
        changed = command(["git", "diff", "--name-only"], target, env).splitlines()
        assert set(changed) <= {"exporter.py"}, f"target changed protected tracked files: {changed}"
        report["checks"]["protected_files_unchanged"] = True
        save()
        # Simulate loss of a closing handoff and withhold the approved crumb so
        # continuation cannot pass merely by repeating the team-memory test.
        # Restore the review state after the recovery run.
        held = {}
        for path in [source / "memory/hot.md", source / "memory/MEMORY.md",
                     source / ".baker/crumbs/legacy-wire-format/000001.json",
                     source / "memory/decisions/legacy-wire-format.md", source / "proposal.md"]:
            if path.exists():
                held[path] = path.read_bytes()
                path.unlink()
        (source / "memory/MEMORY.md").write_text("# Memory index\n\nNo approved crumbs in this recovery fixture.\n")
        brief = command(["baker", "handoff", "--from", args.source], source, env)
        (artifacts / "handoff.md").write_text(brief)
        assert "Last exchange before the session ended" in brief, "Baker did not recover the source transcript"
        assert "semicolon" in brief.lower() and "crlf" in brief.lower(), "recovery lost the wire-format constraint"
        report["checks"]["transcript_recovered_without_hot_or_approved_crumb"] = True
        handoff_prompt = brief + "\n\nThe fixture operator explicitly authorizes continuing now; do not wait for a further confirmation. Implement export_codes and run the checks. Do not approve memory or inspect sibling directories. Work only in this repository.\n"
        try:
            report["timings"]["handoff_target"] = agent(args.target, handoff_prompt, source, artifacts, "handoff-target", env, args.timeout)
            report["checks"]["handoff_export_behavior"] = verify_export(source, env)
        finally:
            for path, content in held.items():
                path.write_bytes(content)
        entries = json.loads(command(["baker", "experimental", "crumb", "recall", "--json"], source, env))
        assert len(entries) == 1 and entries[0]["eligible"]
        with (source / "tests/test_codes.py").open("a") as f:
            f.write("\n# Evidence changed after review.\n")
        entries = json.loads(command(["baker", "experimental", "crumb", "recall", "--json"], source, env))
        assert not entries[0]["eligible"] and entries[0]["body"] == ""
        report["checks"]["changed_evidence_withholds_body"] = True
        report["status"] = "passed"
    except Exception as exc:
        report["status"] = "failed"
        report["error"] = str(exc)
    finally:
        save()
        print(json.dumps(report, indent=2))
    return 0 if report["status"] == "passed" else 1


if __name__ == "__main__":
    sys.exit(main())
