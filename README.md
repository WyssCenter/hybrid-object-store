
<p align="center">
    <img alt="Hoss" src="resources/images/Primary-Wordmark.svg" width="400" />
</p>

---

![Build and Test](https://github.com/gigantum/hybrid-object-store/actions/workflows/main.yaml/badge.svg)
[![Documentation Status](https://readthedocs.org/projects/hybrid-object-store/badge/?version=latest)](https://hybrid-object-store.readthedocs.io/en/latest/?badge=latest)


## Running via Docker Compose
You can run the entire system locally via docker compose. By default it will run on localhost.

The following are the safest steps to get running.

1. Install dependencies
   1. The system needs Docker, Docker Compose, make, and git. Instructions will vary depending on your host OS.
2. run `make reset` (if you've previously run the hos on the current machine)
3. run `make setup` (if you've never run the hos on the current machine)
4. run `make env`
   1. (Optional) Edit the `EXTERNAL_HOSTNAME` and `DOMAIN` variables in the `~/.hos/.env` file if you wish to run on something other than localhost. If performing TLS termination outside of the server (e.g. an ALB in AWS), make sure `EXTERNAL_HOSTNAME` starts with `https://`.
   2. (Optional) The system can automatically provision and renew certificates via Let's Encrypt. To enable this feature, set `LETS_ENCRYPT_ENABLED=true`, make sure `EXTERNAL_HOSTNAME` starts with `https://` and `ADMIN_EMAIL` is set. The `ADMIN_EMAIL` will be used for communication (e.g. expiration notices) from Let's Encrypt. Your server must be reachable on port 80 for the ACME challenge to succeed. The ingress proxy automatically will redirect all other traffic on port 80 to 443, so you can use this to lock down server access if desired (e.g allow 80 from anywhere and 443 from a specific CIDR block).
   3. (Optional) If running behind a load balancer or proxy that will be run a health check on the server, you may need to set the `HEALTH_CHECK_HOST` env var. If this is set, an additional routing rule will be added to ensure that calls to the core service at the `HEALTH_CHECK_HOST` will route. For example, when running in AWS behind and ALB, the server will be configured to route requests to the external FQDN, but the ALB will make health check calls to the internal IP. In this case, you should set `HEALTH_CHECK_HOST` to the internal IP of the instance.
   4. (Optional) If you are using the LDAP integration and wish to enable reCAPTCHA on the login page, you must set the `RECAPTCHA_SITE_KEY` and `RECAPTCHA_SECRET_KEY` env vars to valid values before running `make config`. You can learn more about how to obtain these values for your domain [here](server/dex/README.md).
5. run `make config`
   1. (Optional) Customize the look of the Hoss by modifying the primary colors, server name displayed in the top menu bar, and logos. You can do this by editing `~/.hoss/ui/config.json` and replacing `~/.hoss/ui/logo.svg` and `~/.hoss/ui/favicon.png`. If you wish to also change the Dex login page (e.g. you are using LDAP auth), you can also manually edit the icon `~/.hoss/auth/web/static/img/hoss-logo.svg` and login page `~/.hoss/auth/web/templates/password.html`
6. run `make build`
7. run `make up` or `make up DETACH=true`

In this configuration you should be able to reach minio directly on localhost and will be redirected to the Hoss user interface when you load the server in your browser. The core API should be available on `http://localhost/core/v1/`

Note, if running in a scenario where you do not want all of the services to run (e.g. no sync, minio, auth/ldap services), you can edit the `SERVICES` environment variable. The Makefile uses this variable to know which services to include during `make up`. Currently available service names are `opensearch`, `ldap`, `rabbitmq`, `db`, `dex`, `reverse-proxy`, `auth`, `ldap-admin`, `etcd-0`, `minio`, `core`, `ui`, `sync`.

When running with default configuration, test user accounts will automatically be created. **You should always delete these accounts before running the system for real and if not running on localhost**

| username               | password | role       
|------------------------|----------|------------
| admin@example.com      | foo      | admin      
| privileged@example.org | bar      | privileged 
| user@example.org       | password | base user  
| test.user@example.com  | foobar   | base user  


## Monitoring Server Logs

When running with `make up DETACH=true`, logs are not visible. If you wish to check logs, you can either watch or get logs.

To watch logs run `make watch-logs`. This will tail the last 10 lines of all the services and then follow the output in real-time. You can also optionally use the `SERVICE_NAME` variable to view the logs from only 1 service at a time
(e.g. `make watch-logs SERVICE_NAME=core`). Currently available service names are `opensearch`, `ldap`, `rabbitmq`, `db`, `dex`, `reverse-proxy`, `auth`, `ldap-admin`, `etcd-0`, `minio`, `core`, `ui`, `sync`.

To pull all the logs, run `make get-logs`. This will get all available logs and display without color. You can easily write the logs to a file for later analysis doing something like `make get-logs > logs.txt`. 
Additionally, you can use the `SERVICE_NAME` variable to get logs for a single service.

Note, using the default JSON logging driver can fill storage with time. You can configure a different logging driver to send the logs somewhere, but a simple solution is to configure the JSON driver to limit
log storage. To do this, add something like the following to `daemon.json` and restart the docker daemon.

```
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "50m",
    "max-file": "20" 
  }
}
```

## Restarting a Service
If you make a change or need to restart a service for any reason, you can run `make restart SERVICE_NAME=<service name>`. 

## Running Unit Tests
To run unit tests, you need to run the part of the stack. We currently do not mock since it's
easy to just wipe the data and start over when developing. Tests still have proper fixtures to clean up.

The following are the safest steps to get tests to run. You should run them from the `server` directory.

First start with 1 time setup that is required if you have not done these steps yet to start developing

1. Ensure you have the minio client `mc` installed on your host (required for minio related tests)
2. Ensure you have go installed
3. run `make setup`

Then, run the following to set up the system in a configuration where the services can hit other resources (e.g. database, minio, etc.)

1. run `make reset`
2. run `make env`
3. (Optionally) run `make build` (useful to make sure the compiler is happy with everything. If iterating you don't need to keep rebuilding in most cases because the tests run on your host.)
4. run `make up-testing` to run the required services for unit tests to run.
5. When done, run `make down` to shut down the running services.

You are free to then run tests via VSCode or directly via `go test -v ./...` in the desired location. If something get's messed up, simply restart from step 1. and it will reset all the data.

If you have not generated the API docs locally, it's likely you'll get a build failure. When building the containers this happens automatically for you. You must install [swag](https://github.com/swaggo/swag) locally (tool used to auto-generate swagger docs), and then run `make api-docs`. If using go 1.16 or later:

```
go install github.com/swaggo/swag/cmd/swag@latest
make api-docs
```

If `swag` is not found, make sure your `GOPATH` is on your `PATH`.

## Running Integration Tests
To run the integration tests, you need to build and run the full system and configure a virtualenv to run the integration tests. See detailed instructions in the [`test` package README](test/README.md)



