import click
from pathlib import Path
import os
import stat
import shutil
import tempfile

from rich.progress import Progress
from dotenv import dotenv_values

from hossadm.console import console
from hossadm import configuration
from hossadm import docker
from hossadm import filesystem
from hossadm import opensearch
from hossadm import version

import time


@click.command()
@click.argument('backup_archive', metavar='BACKUP_ARCHIVE', type=str, nargs=1)
@click.option('--config-dir', '-c', type=str, default="~/.hoss", show_default=True,
              help="Server configuration directory")
@click.option('--endpoint', '-e', type=str, default="http://localhost", show_default=True,
              help="Server root endpoint, including the scheme (e.g. https://hoss.mydomain.com")
@click.pass_context
def restore(ctx, backup_archive, config_dir, endpoint):
    """Restore the server running locally"""
    console.clear()
    config_dir = Path(config_dir).expanduser()

    # Verify system is "reset"
    if Path(config_dir, ".env").exists():
        console.print("\n\n  WARNING: An existing installation is detected.",
                      style="white on red")
        console.print("  Reset the local installation before proceeding with a restore.  \n\n",
                      style="white on red")
        return

    console.print("")
    console.print(f"Starting restore from {backup_archive}...")
    console.print("")
    with Progress(console=console) as progress:
        task_prepare = progress.add_task("Preparing server restore", total=5)
        task_start = progress.add_task("Waiting for server to start", start=False)
        task_finish = progress.add_task("Completing restore", start=False)

        progress.console.print(":clamp:  Unpacking backup file...")
        temp_dir = tempfile.TemporaryDirectory()
        backup_dir = Path(temp_dir.name)
        shutil.unpack_archive(backup_archive, backup_dir)
        config = dotenv_values(Path(backup_dir, 'config', ".env").expanduser().as_posix())

        try:
            filesystem.initialize_config_dir(config_dir)
            filesystem.prepare_for_restore(Path(config['BACKUP_ROOT']))
        except PermissionError as err:
            progress.console.print(f":lock:  Failed to prepare filesystem for restore. "
                                   f"Please delete the following directories "
                                   f"(you may require sudo depending on your UID): "
                                   f"{Path(config['BACKUP_ROOT'], '.db').as_posix()}, "
                                   f"{Path(config['BACKUP_ROOT'], '.opensearch').as_posix()}, "
                                   f"{Path(config_dir, 'data', 'opensearch').as_posix()}, "
                                   f"{Path(config_dir, 'data', 'db').as_posix()}", style="red")
            progress.console.print(err)
            time.sleep(5)
            return

        if int(config['UID']) != 1000:
            opensearch_data_dir = Path(config_dir, 'data', 'opensearch')
            opensearch_backup_dir = Path(config['BACKUP_ROOT'], '.opensearch')
            progress.print("")
            progress.print("")
            progress.print("Please run the following commands to prepare bind mount ownership:")
            progress.print("")
            progress.print(f"sudo chmod g+rwx {opensearch_data_dir.as_posix()}")
            progress.print(f"sudo chgrp 1000 {opensearch_data_dir.as_posix()}")
            progress.print(f"sudo chmod g+rwx {opensearch_backup_dir.as_posix()}")
            progress.print(f"sudo chgrp 1000 {opensearch_backup_dir.as_posix()}")
            progress.print("")
            progress.print("")
            filesystem.wait_for_group(opensearch_data_dir, 1000)
            filesystem.wait_for_group(opensearch_backup_dir, 1000)

        progress.advance(task_prepare)

        progress.console.print(":gear:  Restoring configuration...")

        configuration.restore_files(config_dir, Path(backup_dir, 'config'))

        time.sleep(2)
        progress.advance(task_prepare)

        progress.console.print(":computer_disk:  Restoring volumes...")
        docker.restore_auth(backup_dir, 0, 0)
        docker.restore_ldap(backup_dir, 0, 0)
        progress.advance(task_prepare)

        progress.console.print(":file_cabinet:   Preparing database restore...")
        db_backup_dir = Path(config['BACKUP_ROOT'], '.db')
        backup_src = Path(backup_dir, 'db', 'postgres-hoss_core-backup.dump')
        backup_dst = Path(db_backup_dir, 'postgres-hoss_core-backup.dump')
        shutil.move(backup_src, backup_dst)
        backup_src = Path(backup_dir, 'db', 'postgres-hoss_auth-backup.dump')
        backup_dst = Path(db_backup_dir, 'postgres-hoss_auth-backup.dump')
        shutil.move(backup_src, backup_dst)

        # remove the create database script since we'll be restoring
        create_script = Path(config_dir, 'core', 'db-init-scripts', 'create-databases.sh')
        restore_script = Path(config_dir, 'core', 'db-init-scripts', 'zz-restore-db.sh')
        with open(restore_script, 'wt') as f:
            f.write('#!/bin/bash\n\n')
            f.write('set -e\nset -u\n')
            f.write('echo "Restoring database from backup"\n')
            f.write('/usr/local/bin/pg_restore -U $POSTGRES_USER -d hoss_core /mnt/backup/postgres-hoss_core-backup.dump\n')
            f.write('/usr/local/bin/pg_restore -U $POSTGRES_USER -d hoss_auth /mnt/backup/postgres-hoss_auth-backup.dump\n')
            f.write('rm /mnt/backup/postgres-hoss_core-backup.dump\n')
            f.write('rm /mnt/backup/postgres-hoss_auth-backup.dump\n')
            f.write('echo "Database restore complete!"\n')

        st = os.stat(restore_script.as_posix())
        os.chmod(restore_script.as_posix(), st.st_mode | stat.S_IRUSR | stat.S_IXUSR | stat.S_IRGRP | stat.S_IXGRP | stat.S_IROTH | stat.S_IXOTH)
        st = os.stat(create_script.as_posix())
        os.chmod(create_script.as_posix(), st.st_mode | stat.S_IRUSR | stat.S_IXUSR | stat.S_IRGRP | stat.S_IXGRP | stat.S_IROTH | stat.S_IXOTH)
        progress.advance(task_prepare)

        progress.console.print(":magnifying_glass_tilted_right: Preparing search index restore...")
        # move snapshot data
        backup_src = Path(backup_dir, 'search')
        backup_dst = Path(config['BACKUP_ROOT'], '.opensearch')
        shutil.copytree(backup_src.as_posix(), backup_dst, dirs_exist_ok=True)
        time.sleep(2)
        progress.advance(task_prepare)

        # Done preparing - start waiting for the server to start
        progress.console.print(":person_running: *** Ready for restore. Please start the server! ***")
        progress.start_task(task_start)
        progress.update(task_start, total=4)
        time.sleep(5)
        progress.console.print(":hourglass: Waiting for server to start...")

        # Wait for opensearch to be up and ready
        # Note: opensearch.wait_for_server advances the `task_start` progress bar 4 times internally
        opensearch.wait_for_server(progress, task_start)

        # Done starting - start finishing
        progress.start_task(task_finish)
        progress.update(task_finish, total=5)

        # Verify the server is running a version that is compatible with this backup
        # Currently, we require the server major.minor version to be the same, but the build can be different
        is_supported, backup_version, server_version = version.server_version_is_supported(backup_dir, endpoint)
        if not is_supported:
            console.print(f"\n\n  Error: The server version ({server_version['version']}) is not compatible with the "
                          f"backup version {backup_version['version']}).", style="white on red")
            console.print("  Backups are compatible with a server's `major.minor` version. Install a server at a "
                          "compatible version and then perform and upgrade if needed.  \n\n", style="white on red")
            return

        # Run process to restore opensearch snapshot
        progress.console.print(":magnifying_glass_tilted_right: Configuring search index restore...")
        opensearch.create_snapshot_repo()
        progress.advance(task_finish)

        snap_to_restore = opensearch.get_latest_snapshot()
        progress.advance(task_finish)

        opensearch.delete_snapshot_indicies(snap_to_restore)
        progress.advance(task_finish)

        time.sleep(5)
        progress.console.print(":magnifying_glass_tilted_right: Restoring search index...")
        opensearch.restore(snap_to_restore)
        time.sleep(1)
        progress.advance(task_finish)

        progress.console.print(":broom:  Cleaning up...")
        temp_dir.cleanup()
        docker.clear_snapshot_dir()
        time.sleep(1)
        progress.advance(task_finish)

        progress.console.print(":tada:  Restore Complete  :tada:")
        time.sleep(2)


