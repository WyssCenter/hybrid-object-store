import click
from pathlib import Path
import shutil
import datetime
import time

from rich.prompt import Confirm
from rich.progress import Progress
from dotenv import dotenv_values

from hossadm.console import console
from hossadm import configuration
from hossadm import docker
from hossadm import filesystem
from hossadm import opensearch
from hossadm import version


@click.command()
@click.option('--config-dir', '-c', type=str, default="~/.hoss", show_default=True,
              help="Server configuration directory")
@click.option('--endpoint', '-e', type=str, default="http://localhost", show_default=True,
              help="Server root endpoint, including the scheme (e.g. https://hoss.mydomain.com")
@click.pass_context
def backup(ctx, config_dir, endpoint):
    """Backup the server running locally"""
    console.clear()

    config_dir = Path(config_dir).expanduser()
    config = dotenv_values(Path(config_dir, ".env").expanduser().as_posix())

    # Prep backup location
    current_dir = Path(config['BACKUP_ROOT'], ".current").expanduser()
    uid = int(config['UID'])
    gid = int(config['GID'])
    if Path(current_dir).exists():
        console.print("\n\n  WARNING: An existing backup in-progress is detected.  \n\n",
                      style="white on red")
        should_continue = Confirm.ask("Delete partial backup and create new?")
        if not should_continue:
            console.print("- Backup cancelled.")
            return
        else:
            # Clear .current dir.
            shutil.rmtree(current_dir.as_posix())

    # Prep current dir
    filesystem.initialize_current_dir(current_dir)

    # Make sure db backup dir is clean
    db_core_backup = Path(config['BACKUP_ROOT'], '.db', 'postgres-hoss_core-backup.dump')
    if db_core_backup.exists():
        db_core_backup.unlink()
    db_auth_backup = Path(config['BACKUP_ROOT'], '.db', 'postgres-hoss_auth-backup.dump')
    if db_auth_backup.exists():
        db_auth_backup.unlink()

    console.print("")
    with Progress() as progress:
        task = progress.add_task("Backing Up Local Server", total=7)

        progress.console.print(":gear:  Backing up configuration...")
        current_config_dir = Path(current_dir, 'config')
        configuration.backup_files(config_dir, current_config_dir)
        version.save_server_version(current_dir, endpoint)
        time.sleep(2)
        progress.advance(task)

        progress.console.print(":computer_disk:  Backing up volumes...")
        docker.backup_auth(current_dir, uid, gid)
        docker.backup_ldap(current_dir, uid, gid)
        progress.advance(task)

        progress.console.print(":file_cabinet:   Backing up Core service database...")
        docker.backup_database(config['POSTGRES_USER'], config['POSTGRES_PASSWORD'], "hoss_core")
        backup_dst = Path(current_dir, 'db', 'postgres-hoss_core-backup.dump')
        shutil.move(db_core_backup, backup_dst)
        progress.advance(task)

        progress.console.print(":file_cabinet:   Backing up Auth service database...")
        docker.backup_database(config['POSTGRES_USER'], config['POSTGRES_PASSWORD'], "hoss_auth")
        backup_dst = Path(current_dir, 'db', 'postgres-hoss_auth-backup.dump')
        shutil.move(db_auth_backup, backup_dst)
        progress.advance(task)

        progress.console.print(":magnifying_glass_tilted_right:  Backing up search...")
        opensearch.create_snapshot_repo()
        opensearch.backup(config['BACKUP_ROOT'], current_dir)
        progress.advance(task)

        progress.console.print(":clamp:   Packaging backup...")
        archive_name = Path(config['BACKUP_ROOT'], "backups",
                            f"hoss-backup-{datetime.datetime.utcnow().strftime('%Y-%m-%dT%H%M%SZ')}")
        backup_file = shutil.make_archive(archive_name, 'gztar', current_dir)
        progress.advance(task)

        progress.console.print(":broom:  Cleaning up...")
        shutil.rmtree(current_dir.as_posix())
        time.sleep(2)
        progress.advance(task)

        progress.console.print(":tada:  Backup Complete  :tada:")
        time.sleep(2)

    progress.console.print(f"\n\nBackup archive: {backup_file}")




