# Hoss Backup & Restore Design Notes

Initial HOS back up/restore capability will:
* Back up the state of the system at a specific version of the HOS
* Restore the state of the system to the same version it was backed up from
* NOT backup object storage


## Resources

The following items must be backed up.

* Config
    * ~/.hoss/.env
    * ~/.hoss/traefik.yaml
    * ~/.hoss/sync
    * ~/.hoss/elasticsearch
    * ~/.hoss/core
    * ~/.hoss/auth
* Auth
    * The contents of the `server_auth-secrets` volume
        * This is the private pem
* ldap
    * The contents of volumes `server_ldap-vol0` and `server_ldap-vol1`
        * This is the ldap config and database
* Postgres
    * Backup all of the databases
        * Some examples on how to do this [https://simplebackups.com/blog/docker-postgres-backup-restore-guide-with-examples/](https://simplebackups.com/blog/docker-postgres-backup-restore-guide-with-examples/)
* Minio
    * None
* Ectd
    * None
    * Policies will re-render as users login
* Elasticsearch
    * [https://opensearch.org/docs/latest/opensearch/snapshot-restore/](https://opensearch.org/docs/latest/opensearch/snapshot-restore/)
* Sync
    * None

## Backup

Backup will be done via the `hossadm` administrator python package that is located in the `admin` directory of this repository.

User will not need sudo, but will need to be able to run docker commands.

A new environment variable `BACKUP_ROOT` will be added to the .env file. The default location will be `$HOSS_WORKING_DIR/backup`, This directory will have the following structure:

```
$BACKUP_ROOT/
    |_ .current/
        |_ config/
        |_ volumes/
        |_ db/
        |_ search/
    |_ .opensearch
    |_ .db
    |_ backups

```

1. The .current backup folder is the location where backup resources are written as the backup is in process. It is cleared if there are any contents (maybe with a warning if not empty?) at the start of the back up process.
2. .db directory is bind mounted into the database container and will contain snapshots during backup and restore.
3. the .opensearch directory is bind mounted into the opensearch container and will contain snapshots during backup and operation. These snapshots are cumulative.
4. All config data is copied to the `$BACKUP_ROOT/config` dir
5. Auth volume is backed up
    1. Run something like `docker run -it -v server_auth-secrets:/mnt/auth -v /backup-root/.current/volumes:/mnt/backup busybox cp -R /mnt/auth /mnt/backup`
6. Ldap-admin volumes are backed up
    1. Run something like `docker run -it -v server_ldap-vol0:/mnt/ldap-vol0 -v /backup-root/.current/volumes:/mnt/backup busybox cp -R /mnt/ldap-vol0 /mnt/backup`
    2. Run something like `docker run -it -v server_ldap-vol1:/mnt/ldap-vol1 -v /backup-root/.current/volumes:/mnt/backup busybox cp -R /mnt/ldap-vol1 /mnt/backup`
7. Run a backup of the postgresdb
    1. Postgres container modified to mount /backup-root/.db:/mnt/backup
    2. Run something like this: 

        docker exec &lt;postgresql_container> /bin/bash \


         -c "export PGPASSWORD=&lt;postgresql_password> \


             && /usr/bin/pg_dump -U &lt;postgresql_user> hos_auth" \


          > /mnt/backup/postgres-hos_auth-backup.sql


        docker exec &lt;postgresql_container> /bin/bash \


         -c "export PGPASSWORD=&lt;postgresql_password> \


             && /usr/bin/pg_dump -U &lt;postgresql_user> hos_core" \


          > /mnt/backup/postgres-hos_core-backup.sql

    3. Move postgres-backup.sql to /backup-root/.current/database
8. Run elasticsearch snapshot and tar data
    1. [https://opensearch.org/docs/latest/opensearch/snapshot-restore/#about-snapshots](https://opensearch.org/docs/latest/opensearch/snapshot-restore/#about-snapshots)
    2. Snapshots will go into /backup-root/.elasticsearch
    3. Once ready, contents of /backup-root/.elasticsearch will be tar’d and placed in /backup-root/.current/search
9. RabbitMQ volume is backed up
    1.  Run something like `docker run -it -v server_rabbitmq-data:/mnt/rabbitmq -v /backup-root/.current/volumes:/mnt/backup busybox cp -R /mnt/rabbitmq /mnt/backup`
10. The /backup-root/.current/ dir is tar.gz’d and placed in the /backup-root/backups directory


## Restore

Restore will be done via a bash script where the user provides the path to the backup tar.gz file



1. First the backup file is uncompressed to temp space
2. The `~/.hos` directory is created and populated
3. All volumes are created if missing and populated using the reverse process as backup
4. Prepare postgres restore
    1. Place the sql dump in the `/backup-root/.db` directory
    2. Wipe the current postgres data directory if it exists
    3. Write an additional startup script to do a restore
        1. Run a restore using something like:
            1. pg_restore -U postgres -d hos_auth /tmp/db/postgres-hos_auth-backup.sql
            2. pg_restore -U postgres -d hos_core /tmp/db/postgres-hos_core-backup.sql
5. Prep elasticsearch restore
    4. Wipe the /backup-root/.elasticsearch location
    5. Place the archived snapshots in /backup-root/.elasticsearch
6. Run `make up` and wait for system to be ready
7. Run elasticsearch restore
    6. [https://opensearch.org/docs/latest/opensearch/snapshot-restore/#restore-snapshots](https://opensearch.org/docs/latest/opensearch/snapshot-restore/#restore-snapshots)