# Hoss Server Admin Tool
The Hoss admin tool provides a CLI to perform administrative operations.

Currently the tool provides backup and restore operations.

## Installing the admin tool
1. From the `admin/` directory run `pip3 install .`


## Configuring a Hoss installation
You must set the backup root directory **at install time**. You cannot change the backup location after a server has started, and must to a backup/restore to modify these settings.

By default, the backup root dirctory will be located at `~/.hoss/backup`. To change this location, edit the `BACKUP_ROOT` variable in the `~/.hoss/.env` file. Then run `make config` before running `make up.

The structure of the backup directory is:

- $BACKUP_ROOT/backups: Location where successfully created backup archives will be placed.
- $BACKUP_ROOT/.db: A directory that will be bind mounted into the database container. Database dumps are placed here during backup/restore
- $BACKUP_ROOT/.opensearch: A directory that will be bind mounted into the opensearch container. Index snapshots are placed here during backup/restore
- $BACKUP_ROOT/.current: A directory that contains a backup in progress. The contents of this directory are compressed into a `.tar.gz` file after all items are successfully collected.

Currently, it is not recommended that the backup directory be on an NFS mount due to the need to control file permissions during the backup/restore process. It is safe to move the backup archive after
it has been created.

## Backup
To start a backup, run `hossadm backup` as the user who runs the Hoss (i.e. what user ran `make up`)

This will create a backup of the current state of the system. It does NOT backup object store data.

The backup data is compressed into a single archive and placed in the `$BACKUP_ROOT/backups` directory.

## Restore

To prepare for a restore, you should have a "clean" Hoss working directory (i.e. `~/.hoss`) before running this command. If you are restoring
on the same system that was previously backed up, you'll need to do `make down` and `make reset` to clear most resources. Also you will likely 
need to manually remove a few directories with `sudo` because of permission changes (e.g. ~/.hoss/data/db, ~/.hoss/backup/.db, ~/.hoss/data/opensearch, ~/.hoss/backup/.opensearch)

If the default `BACKUP_ROOT` location was used, you can remove all content from the working directory except the `backup` directory because you'll need the backup archive in that directory.

To start a restore, run `hossadm restore <PATH TO BACKUP ARCHIVE>` as the user who runs the Hoss (i.e. what user ran `make up`)

During the restore process, the `hossadm` tool will instruct you to start the server. At this point, in a different terminal run `make up DETACH=true` as the user who runs the Hoss.
`hossadm` will detect the server starting up and continue with the restore process.

## Version
To print the version of the current Hoss install, run `hossadm version`

## Update index
To update the metadata index to a new version, run `hossadm update-index --version #.#.#`. This will re-index the existing data in the metadata index to a new version of the index with an updated mapping.
 
It is recommended to run a backup of the server before running the migration because any errors could result in data loss.
 
## Future index development
When `HOSS` development introduces an updated index mapping, make sure to place a json file defining the new version's mapping in the `metadata-index-versions` directory, named according to the server version (`#.#.#.json`). This should be the same index mapping as defined in the `HOSS` sync service's `elastic` package.
 
The migration command is currently intended for migrating between index versions with the same core fields; if a new index with additional fields is required the command may need to be updated with custom scripting to add the new field to existing documents (see [update-by-query](https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-update-by-query.html#docs-update-by-query-api-source) docs).
