# Backup and Restore

Backup and Restore functionality is provided by the `hossadm` package, which is included in the Hoss repository. You must install the tool to perform
backup and restore actions.

Backup does not include object store data! When a backup is run, the server state is captured, including:

* Server configuration (i.e. contents of the `~/.hoss` directory)
* Database tables
* Internal LDAP server data (if enabled)
* Private keys
* Search index

The backup is packed into a single archive that can then be moved as needed to safe storage. Since this archive will contain sensitive information
like passwords and private keys, it should be treated securely. 

```{warning}
The backup process does not include object store data. You should protect object store data via other means, such as replication or bucket versioning.
```

## Configuring Backup
You must set the backup root directory **at install time**. You cannot change the backup location after a server has started, and must do a backup/restore to modify these settings.

By default, the backup root directory will be located at `~/.hoss/backup`. To change this location, edit the `BACKUP_ROOT` variable in the `~/.hoss/.env` file. Then run `make config` before running `make up.

The structure of the backup directory is:

- $BACKUP_ROOT/backups: Location where successfully created backup archives will be placed.
- $BACKUP_ROOT/.db: A directory that will be bind mounted into the database container. Database dumps are placed here during backup/restore
- $BACKUP_ROOT/.opensearch: A directory that will be bind mounted into the opensearch container. Index snapshots are placed here during backup/restore
- $BACKUP_ROOT/.current: A directory that contains a backup in progress. The contents of this directory are compressed into a `.tar.gz` file after all items are successfully collected.

Currently, it is not recommended that `BACKUP_ROOT` be on an NFS mount due to the need to control file permissions during the backup/restore process. It is safe to move the backup archive after
it has been created.


## Installing the `hossadm` tool

If you have yet to install the `hossadm` tool:

1. Create and activate a new Python3 virtual environment.
2. From the `admin/` directory if the Hoss code repository run `pip3 install -U .`

If the tool has been updated, simply run `pip3 install -U .` again after updating the Hoss code repository to
the desired version.

## Backup

To start a backup, run `hossadm backup` as the user who runs the Hoss (i.e. what user ran `make up`) while the
server is running.

If you are not developing, and running on localhost, you must include the `--endpoint` option to indicate your server's external endpoint, e.g.:

```
hossadm backup --endpoint https://hoss.mycompany.com
```

This will create a backup of the current state of the system. It does NOT backup object store data.

The backup data is compressed into a single archive and placed in the `$BACKUP_ROOT/backups` directory.

## Restore

```{warning}
If you are using minIO locally, be **very** careful with the `make reset` command. This could delete all of your data by removing the data
directory!!!
```

To prepare for a restore, you should have a "clean" Hoss working directory (i.e. `~/.hoss`) before running this command. If you are restoring
on the same system that was previously backed up, you'll need to do `make down` and `make reset` to clear most resources. Also you will likely 
need to manually remove a few directories with `sudo` because of permission changes (e.g. `~/.hoss/data/db`, `~/.hoss/backup/.db`, `~/.hoss/data/opensearch`, `~/.hoss/backup/.opensearch`). If using minIO locally, you should 

If the default `BACKUP_ROOT` location was used, you can remove all content from the working directory except the `backup` directory because you'll need the backup archive in that directory.

To start a restore, run `hossadm restore <PATH TO BACKUP ARCHIVE>` as the user who runs the Hoss (i.e. what user ran `make up`)

If you are not developing, and running on localhost, you must include the `--endpoint` option to indicate your server's external endpoint, e.g.:

```
hossadm restore --endpoint https://hoss.mycompany.com ~/hoss-backup-2022-02-22T014536Z.tar.gz
```

During the restore process, the `hossadm` tool will instruct you to start the server. At this point, in a different terminal run `make up DETACH=true` as the user who runs the Hoss in `server` directory of the Hoss code repository. `hossadm` will detect the server starting up and continue with the restore process.

Finally, if running minIO, it is recommended that you restart the sync and core service to ensure that any timing issues during start up are resolved
immediately.

```
make restart SERVICE_NAME=sync
make restart SERVICE_NAME=core
```
