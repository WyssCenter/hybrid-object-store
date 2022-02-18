import shutil
from pathlib import Path

CONFIG_FILES = ['.env', 'traefik.yaml']
CONFIG_DIRS = ['core', 'auth', 'sync', 'opensearch']


def backup_files(config_dir: Path, output_dir: Path) -> None:
    """Method to backup config files

    Args:
        config_dir: root configuration directory (e.g. ~/.hoss)
        output_dir: destination directory for configuration data

    Returns:
        None
    """
    for f in CONFIG_FILES:
        src = Path(config_dir, f)
        dst = Path(output_dir, f)
        shutil.copyfile(src, dst)

    for d in CONFIG_DIRS:
        src = Path(config_dir, d)
        dst = Path(output_dir, d)
        shutil.copytree(src, dst)


def restore_files(config_dir: Path, input_dir: Path) -> None:
    """Method to restore config files

    Args:
        config_dir: root configuration directory (e.g. ~/.hoss)
        input_dir: source directory for configuration data

    Returns:
        None
    """
    for f in CONFIG_FILES:
        src = Path(input_dir, f)
        dst = Path(config_dir, f)
        shutil.copyfile(src, dst)

    for d in CONFIG_DIRS:
        src = Path(input_dir, d)
        dst = Path(config_dir, d)
        if dst.exists():
            shutil.rmtree(dst)
        shutil.copytree(src, dst)
