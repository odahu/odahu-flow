import subprocess
from pathlib import Path


def zip_dir(source: Path, target: Path):
    if not target.exists():
        target.touch()

    subprocess.run([
        "tar", "--exclude", target,
        "-cv", "--use-compress-program=pigz", "-f", target, ".",
    ], cwd=str(source))


def unzip(source: Path, target: Path):
    subprocess.run(["tar", "-xvf", source, "-C", "."], cwd=str(target))
