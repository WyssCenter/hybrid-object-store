from typing import Tuple, Optional

import requests
from pathlib import Path
import json


def get_server_version(endpoint: str = "http://localhost") -> dict:
    """Function to fetch the discover endpoint from the local server

    Args:
        endpoint: The server root endpoint, including the protocol without a trailing slash.

    Returns:
        a dict with the version and build hash
    """
    resp = requests.get(f"{endpoint}/core/v1/discover")
    if resp.status_code != 200:
        raise Exception("Failed to load version information from server")

    data = resp.json()

    return {"version": data['version'], "build": data['build']}


def save_server_version(current_dir: str, endpoint: str = "http://localhost") -> None:
    """Function to save the current server version information to the current in-progress backup

    Args:
        current_dir: directory for the current in-progress backup
        endpoint: The server root endpoint, including the protocol without a trailing slash.

    Returns:
        None
    """
    version = get_server_version(endpoint)

    with open(Path(current_dir, "version.json"), "wt") as f:
        json.dump(version, f)


def server_version_is_supported(backup_dir: str, endpoint: str = "http://localhost") -> Tuple[bool, dict, dict]:
    """Function to check if the current server version is supported to complete a restore

    Args:
        backup_dir: Directory containing and unpacked backup archive that is being restored
        endpoint: The server root endpoint, including the protocol without a trailing slash.

    Returns:
        A tuple with the a boolean indicating if the backup is compatible with the running server, the backup's version
        information, the server's version information
    """
    with open(Path(backup_dir, "version.json"), "rt") as f:
        backup_version = json.load(f)

    server_version = get_server_version(endpoint)

    supported = False
    backup_parts = backup_version['version'].split('.')
    server_parts = server_version['version'].split('.')
    if (backup_parts[0] == server_parts[0]) and (backup_parts[1] == server_parts[1]):
        supported = True

    return supported, backup_version, server_version
