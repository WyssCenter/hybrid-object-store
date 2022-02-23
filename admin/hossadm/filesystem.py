from typing import Optional
from pathlib import Path
import shutil
import os
import stat
import time


def initialize_current_dir(current_dir: Path) -> None:
    current_dir.mkdir()

    current_config_dir = Path(current_dir, 'config')
    current_config_dir.mkdir()

    Path(current_dir, 'volumes').mkdir()
    Path(current_dir, 'db').mkdir()


def _recreate_dir(directory: Path) -> None:
    if directory.exists():
        shutil.rmtree(directory)

    directory.mkdir(parents=True)


def initialize_config_dir(config_dir: Path) -> None:
    config_dir.mkdir(exist_ok=True)

    data_db = Path(config_dir, 'data', 'db')
    _recreate_dir(data_db)

    data_opensearch = Path(config_dir, 'data', 'opensearch')
    _recreate_dir(data_opensearch)
    os.chmod(data_opensearch.as_posix(), stat.S_IRWXU | stat.S_IRWXG)

    data_events = Path(config_dir, 'data', 'events')
    _recreate_dir(data_events)

    # This creates a "default" bucket based on our current default config files
    Path(config_dir, 'data', 'nas', 'data').mkdir(parents=True, exist_ok=True)

    Path(config_dir, 'ui').mkdir(parents=True, exist_ok=True)
    Path(config_dir, 'auth', 'certificates').mkdir(parents=True, exist_ok=True)


def prepare_for_restore(backup_root: Path) -> None:
    backup_db = Path(backup_root, '.db')
    _recreate_dir(backup_db)

    backup_opensearch = Path(backup_root, '.opensearch')
    _recreate_dir(backup_opensearch)

    os.chmod(backup_opensearch.as_posix(), stat.S_IRWXU | stat.S_IRWXG)


def wait_for_group(path: Path, gid: int) -> None:
    for _ in range(30*3):
        s = os.stat(path.as_posix())
        if s.st_gid == gid:
            return
        else:
            time.sleep(2)

    raise Exception(f"Failed to detect proper ownership of {path} after 3 minutes. Try again.")
