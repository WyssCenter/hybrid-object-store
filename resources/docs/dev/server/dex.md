## dexidp integration

Dex is an identity service that uses OpenID Connect to drive authentication for other apps. The auth service interfaces with
dex to integrate with identity providers and then generate JWTs to use within the Hoss system.

## Images

We currently support two ways to run dex. The first is to simply use the official docker image
provided by the dex team. The second is to use a fork that we maintain that adds the ability
to use Google's reCAPTCHA on the login page hosted by Dex. This page is used with the LDAP
integration, including when the internal LDAP provider is enabled.

### Using the official image
To use the official image, simply leave the `RECAPTCHA_SITE_KEY` env var empty in the `~/.hoss/.env` file.

With this var unset, when `make config` is run, the `dex/dex-config-no-captcha.tmpl` file will be appended to 
the end of the dex config and the `dex/password.html` file will be copied into the `~/.hoss/auth/web/templates` directory. This
directory is mounted into the dex container and used to customize the dex web interface.

When other `make` commands are run, the `dex/docker-compose.yaml` file will be added to the compose files used.

### Using the fork

To use the fork and enable reCAPTCHA, set the `RECAPTCHA_SITE_KEY` and `RECAPTCHA_SECRET_KEY` env vars in the `~/.hoss/.env` file to [valid values](../../../server/dex/README.md).

With these vars set, when `make config` is run, the `dex/dex-config-captcha.tmpl` file will be appended to 
the end of the dex config and the `dex/password-captcha.html` file will be copied into the `~/.hoss/auth/web/templates` directory. This
directory is mounted into the dex container and used to customize the dex web interface.

When other `make` commands are run, the `dex/docker-compose-recaptcha.yaml` file will be added to the compose files used.

## Updating Dex
The version of dex used is in each docker compose file. When updating dex, you should update both images.

To update the official image, edit the `image:` line in `server/dex/docker-compose.yaml`.

To update the forked image:

1) Update the [dex fork](https://github.com/WyssCenter/dex) from the upstream
2) Merge main or the tag you updated into the branch `recaptcha`
3) Update the `VERSION` env var in the make file to the desired version number
4) Run `make docker-image`
5) Run `make docker-push`
6) edit the `image:` line in `server/dex/docker-compose-recaptcha.yaml`.
