# Updating a Server
To update a server, you update the underlying code repository, rebuild, and redeploy. There are several
steps and considerations to make during this process as outlined below.

## Updating a Single Server
To update a server you simply update the configuration and code and redeploy any changes. Follow the
steps below:

```{note}
The update process does not provide a method to lock state. You should update a server when it is not
in use to avoid consistency issues with users manipulating the system while services restart.
```

1) **Backup**
   
   You should backup a server before the update in case something goes wrong. This will let you restore
   the server back to a functioning state and try again if desired. To backup the server, follow the
   instructions in the [Backup](backup-and-restore.md#backup) section.

2) **Review release notes and make required changes**
   
   The [Release Notes](https://github.com/WyssCenter/hybrid-object-store/releases) will indicate what changes are needed to
   local configuration (e.g. a new field was added to a config file, you have to update search indices, etc.).

   It is very important to follow the steps required between the version you have installed and the version you plan to install.
   Currently there is no requirement to make intermediate updates if moving several versions, but that could change in the future.

3) **Update Repository**
   
   Update the server code repository to the desired version. If you are running from `main` you can likely run `git pull`. If
   you are running a specific release, or wish to update to a specific release, checkout the release using `git checkout <version>`
   where version is the tag of the release (e.g. `0.2.5`).

4) **Update configuration**
   
   Next, run `make config` to update an required config files. This can sometimes be skipped depending on the update, but generally is
   safe to run.

5) **Build Images**

   Run `make build` to build any container that has changed during the update.

6) **Restart/Start Updated Services**
   
   Run `make up DETACH=true` to re-create any service that has changed during the update. This may cause additional service to re-create or restart
   if they are dependent on each other. After a minute everything should be up and ready to go.

7) **Backup Again**
   
   Backups are linked to the server version due to possible incompatible changes between versions. This means the backup you ran at the start of this
   process will no longer work if you tried to restore it into the server at the version you just installed. It is always best practice to immediately
   re-run the backup process to make sure your latest backup will restore at the updated version.

   To backup the server, follow the instructions in the [Backup](backup-and-restore.md#backup) section again.


## Updating Linked Servers

When updating servers that are linked via a shared auth service, you must consider the order in which the servers are updated. Because the
auth service needs to be up and running when a server restarts you should be cautious in the order and timing of your updates. Typically you should:

1) Backup all servers
2) Update the server that is running the auth service, following the single server update instructions above
3) Update all other servers, one at a time, following the single server update instructions above