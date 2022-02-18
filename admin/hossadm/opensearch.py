import time
import json
import datetime
from pathlib import Path
import shutil
import subprocess
from rich.progress import Progress, TaskID

from hossadm import docker


def wait_for_server(progress: Progress, task_id: TaskID) -> None:
    """Wait for the opensearch service to be up and running

    Note: this function will advance the specified TaskID progress bar 4 times

    Args:
        progress: the rich Progress instance used in the main application
        task_id: the rich TaskID to advance while progressing

    Returns:
        None
    """
    search_up = False
    for _ in range(int(3 * 60 / 5)):
        if docker.container_is_running('server_opensearch_1'):
            search_up = True
            break
        else:
            time.sleep(5)

    if not search_up:
        raise Exception("Server not detected within 3 minutes. Did you run 'make up'?")
    progress.advance(task_id)
    time.sleep(5)
    progress.console.print(":hourglass: Server detected, waiting for services to be ready...")
    progress.advance(task_id)
    search_ready = False
    for _ in range(int(3 * 60 / 5)):
        try:
            result = docker.opensearch_api_request("GET", f"_cluster/health")
        except subprocess.CalledProcessError:
            progress.console.print(":hourglass: Search service still starting...")
            time.sleep(5)
            continue

        # Because we run a single node, status will be yellow when ready instead
        # of green because there are no replicas for the indices
        if result['status'] == 'yellow':
            search_ready = True
            break
        else:
            progress.console.print(":hourglass: Search service initializing...")
            time.sleep(5)

    if not search_ready:
        raise Exception("Server (opensearch) failed to be ready within 3 minutes.")

    progress.advance(task_id)
    time.sleep(10)
    progress.advance(task_id)
    return


def create_snapshot_repo() -> None:
    """Create the snapshot repo for snapshots. It is OK to run this if the repo already exists

    Returns:
        None
    """
    data = {"type": "fs",
            "settings":
                {"location": "/mnt/snapshots"}
            }
    result = docker.opensearch_api_request("PUT", "_snapshot/hoss-backup-repository", data)
    if "acknowledged" not in result:
        raise Exception(f"Failed to create snapshot repository: {result}")
    if not result['acknowledged']:
        raise Exception(f"An error occurred while configuring opensearch snapshot repository: {result}")
    time.sleep(5)


def backup(backup_root: str, current_backup_dir: str) -> None:
    """Create a snapshot

    Args:
        backup_root: The root backup directory (e.g. ~/.hoss/backup)
        current_backup_dir: The directory of the current backup (e.g. ~/.hoss/backup/.current)

    Returns:
        None
    """
    snapshot_name = f"hoss-snap-{datetime.datetime.timestamp(datetime.datetime.utcnow())}"
    result = docker.opensearch_api_request("PUT", f"_snapshot/hoss-backup-repository/{snapshot_name}")
    time.sleep(1)

    state = "FAILED"
    for _ in range(120):
        result = docker.opensearch_api_request("GET", f"_snapshot/hoss-backup-repository/{snapshot_name}")
        state = result['snapshots'][0]['state']
        if state == "SUCCESS" or state == "FAILED":
            break
        else:
            time.sleep(5)

    if state != "SUCCESS":
        raise Exception(f"Failed to backup opensearch index: {result}")

    # move snapshot data
    time.sleep(3)
    backup_src = Path(backup_root, '.opensearch')
    backup_dst = Path(current_backup_dir, 'search')
    shutil.copytree(backup_src.as_posix(), backup_dst)


def get_latest_snapshot() -> dict:
    """Get the latest snapshot from opensearch

    Returns:
        a dict containing the snapshot's info
    """
    result = docker.opensearch_api_request("GET", f"_snapshot/hoss-backup-repository/_all")
    sorted_snaps = sorted(result['snapshots'], key=lambda x: x['end_time_in_millis'], reverse=True)
    return sorted_snaps[0]


def delete_snapshot_indicies(snap_to_restore: dict) -> None:
    """Delete all indices in the snapshot to restore.

    We do this because you can't restore an index if it already exists without running
    a rename process.

    Args:
        snap_to_restore: The snapshot to restore, provided by the `_snapshot/hoss-backup-repository/_all` endpoint

    Returns:
        None
    """
    for index in snap_to_restore['indices']:
        result = docker.opensearch_api_request("DELETE", index)
        if not result['acknowledged']:
            raise Exception(f"Failed to delete index '{index}' before restore.")


def restore(snap_to_restore: dict) -> None:
    """Trigger restore in opensearch of the specified snapshot

    Args:
        snap_to_restore: The snapshot to restore, provided by the `_snapshot/hoss-backup-repository/_all` endpoint

    Returns:
        None
    """
    data = {'include_global_state': True}
    result = docker.opensearch_api_request("POST", f"_snapshot/hoss-backup-repository/"
                                                   f"{snap_to_restore['snapshot']}/_restore", data=data)
    if not result['accepted']:
        raise Exception(f"Failed to start index restore: {result}")

    
def update_index(version: str) -> None:
    """Update index mapping to a new version and re-index all existing documents in place"""
    with open(f"hossadm/metadata-index-versions/{version}.json") as f:
        mappings = json.load(f)
    
    # update mapping for existing metadata index
    result = docker.opensearch_api_request("PUT", "metadata-index/_mapping", mappings["mappings"])
    if not result.get('acknowledged'):
        raise Exception(f"Failed to update the metadata index mappings: {result}")

    # re-index documents in place according to the new mappings
    result = docker.opensearch_api_request("POST", "metadata-index/_update_by_query")
    if len(result['failures']) > 0:
        raise Exception(f"Failed to re-index existing metadata according to the new mappings: {result['failures']}")
    print(f"Re-indexed {result['updated']}/{result['total']} documents in the metadata index with no failures")
