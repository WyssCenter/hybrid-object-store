# Monitoring Logs

The services running on a Hoss server generate log messages that can be useful when debugging issues or trying to
understand usage. As a system administrator with access to the host that is running the Hoss software, you
can easily view and extract log data.

When running the server in the foreground (via `make up`), all log data is printed to stdout and should appear in your console.

When running the server in the background with `make up DETACH=true`, logs are not visible. If you wish to check logs, you can either watch or get logs.

To watch logs run `make watch-logs`. This will tail the last 10 lines of all the services and then follow the output in real-time. You can also optionally use the `SERVICE_NAME` variable to view the logs from only 1 service at a time
(e.g. `make watch-logs SERVICE_NAME=core`). Currently available service names are `opensearch`, `ldap`, `rabbitmq`, `db`, `dex`, `reverse-proxy`, `auth`, `ldap-admin`, `etcd-0`, `minio`, `core`, `ui`, `sync`.

To get all available log data, run `make get-logs`. This will get all available logs and display without color codes. You can then write the logs to a file for later analysis doing something like `make get-logs > logs.txt`. 
Additionally, you can use the `SERVICE_NAME` variable to get logs for a single service (e.g. `make get-logs SERVICE_NAME=core`).


## Configure Docker Logging Driver

Note, using the default JSON logging driver can fill storage with time. You can configure a different logging driver to send the logs somewhere, but a simple solution is to configure the JSON driver to limit
log storage. To do this, add something like the following to `daemon.json` and restart the docker daemon.

```json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "50m",
    "max-file": "10" 
  }
}
```

This will limit log files to 50 MB in size before rolling the file, with a maximum of 10 files. In this configuration, each service would generate at a maximum 500MB of log data, for a total 
maximum of 6.5GB. You should modify these values to meet your needs and hardware configuration.