from typing import List, Optional
import shlex
import subprocess
from pathlib import Path
import json


def _run_command(volumes: List[str], uid: int, gid: int, command: str) -> None:
    vol_str = " ".join([f"-v {v}" for v in volumes])
    cmd_str = f'docker run --user="{uid}:{gid}" {vol_str} busybox {command}'

    subprocess.run(shlex.split(cmd_str), check=True)


def _check_or_create_volume(volume_name: str) -> None:
    cmd = f'docker volume ls -q -f "name={volume_name}"'
    result = subprocess.run(shlex.split(cmd), check=True, capture_output=True)
    exists = result.stdout.decode().strip()

    if not exists:
        cmd = f'docker volume create {volume_name}'
        subprocess.run(shlex.split(cmd), check=True, capture_output=True)


def backup_auth(current_backup_dir: Path, uid: int, gid: int) -> None:
    vols = ["server_auth-secrets:/mnt/auth",
            f"{current_backup_dir.expanduser().as_posix()}/volumes:/mnt/backup"]

    # Copy to tmp to modify perms so copying to local fs does not then need sudo
    # to remove file or modify perms on copy
    cmd = 'cp -a /mnt/auth /mnt/backup'
    _run_command(vols, 0, 0, cmd)

    cmd = f'chown -R {uid}:{gid} /mnt/backup/auth'
    _run_command(vols, 0, 0, cmd)


def restore_auth(current_backup_dir: Path, uid: int, gid: int) -> None:
    _check_or_create_volume("server_auth-secrets")
    vols = ["server_auth-secrets:/mnt/auth",
            f"{current_backup_dir.expanduser().as_posix()}/volumes/auth:/mnt/backup"]
    cmd = "cp -a /mnt/backup/. /mnt/auth/"
    _run_command(vols, uid, gid, cmd)

    # Set to correct ownership & permissions
    cmd = "chown 1001:1001 -R /mnt/auth/"
    _run_command(vols, uid, gid, cmd)

    cmd = "chmod 600 /mnt/auth/private.pem"
    _run_command(vols, uid, gid, cmd)


def backup_ldap(current_backup_dir: Path, uid: int, gid: int) -> None:
    vols = ["server_ldap-vol0:/mnt/ldap-vol0",
            f"{current_backup_dir.expanduser().as_posix()}/volumes:/mnt/backup"]
    cmd = "cp -a /mnt/ldap-vol0 /mnt/backup"
    _run_command(vols, 0, 0, cmd)
    cmd = f'chown -R {uid}:{gid} /mnt/backup/ldap-vol0'
    _run_command(vols, 0, 0, cmd)

    vols = ["server_ldap-vol1:/mnt/ldap-vol1",
            f"{current_backup_dir.expanduser().as_posix()}/volumes:/mnt/backup"]
    cmd = "cp -a /mnt/ldap-vol1 /mnt/backup"
    _run_command(vols, 0, 0, cmd)
    cmd = f'chown -R {uid}:{gid} /mnt/backup/ldap-vol1'
    _run_command(vols, 0, 0, cmd)


def restore_ldap(current_backup_dir: Path, uid: int, gid: int) -> None:
    _check_or_create_volume("server_ldap-vol0")
    _check_or_create_volume("server_ldap-vol1")

    vols = ["server_ldap-vol0:/mnt/ldap-vol0",
            f"{current_backup_dir.expanduser().as_posix()}/volumes/ldap-vol0:/mnt/backup"]
    cmd = "cp -a /mnt/backup/. /mnt/ldap-vol0/"
    _run_command(vols, uid, gid, cmd)

    cmd = "chown 911:911 -R /mnt/ldap-vol0/"
    _run_command(vols, uid, gid, cmd)

    vols = ["server_ldap-vol1:/mnt/ldap-vol1",
            f"{current_backup_dir.expanduser().as_posix()}/volumes/ldap-vol1:/mnt/backup"]
    cmd = "cp -a /mnt/backup/. /mnt/ldap-vol1/"
    _run_command(vols, uid, gid, cmd)

    cmd = "chown 911:911 -R /mnt/ldap-vol1/"
    _run_command(vols, uid, gid, cmd)


def backup_database(user: str, password: str, database_name: str) -> None:
    result = subprocess.run(shlex.split('docker ps -aqf "name=server_db_1"'), check=True, capture_output=True)
    container_id = result.stdout.decode().strip()

    cmd = f'docker exec -e PGPASSWORD="{password}" {container_id} /bin/bash -c "/usr/local/bin/pg_dump -Fc -U {user} ' \
          f'{database_name} > /mnt/backup/postgres-{database_name}-backup.dump"'

    subprocess.run(shlex.split(cmd), check=True)


def opensearch_api_request(method: str, path: str, data: Optional[dict] = None) -> dict:
    result = subprocess.run(shlex.split('docker ps -aqf "name=server_opensearch_1"'), check=True, capture_output=True)
    container_id = result.stdout.decode().strip()

    if data:
        data_str = f"-d '{json.dumps(data)}' "
    else:
        data_str = ""

    if method == "GET":
        header = ""
    else:
        header = '-H "Content-Type: application/json" '

    cmd = shlex.split(f'docker exec -i {container_id} /bin/bash ')
    cmd.append("-c")
    cmd.append(f'curl -s {header}-X {method} {data_str}http://localhost:9200/{path}')

    result = subprocess.run(cmd, check=True, capture_output=True)
    result_str = result.stdout.decode().strip()
    if not result_str:
        # Open search should always return something OR check=True would catch the error
        raise Exception("Got an unexpected empty response from the opensearch API.")

    return json.loads(result_str)


def clear_snapshot_dir() -> None:
    result = subprocess.run(shlex.split('docker ps -aqf "name=server_opensearch_1"'),
                            check=True, capture_output=True)
    container_id = result.stdout.decode().strip()

    cmd = shlex.split(f'docker exec -i --user="0:0" {container_id} /bin/bash -c "rm -rf /mnt/snapshots/*"')
    subprocess.run(cmd, check=True)


def container_is_running(name: str) -> bool:
    result = subprocess.run(shlex.split(f'docker ps -qf "name={name}"'), capture_output=True)
    if result.stdout.decode().strip():
        return True
    else:
        return False
