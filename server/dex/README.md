
## RECAPTCHA Configuration

Dex is an identity service that uses OpenID Connect to drive authentication for other apps.

When auth is enabled via LDAP, including the "internal" `ldap` and `ldap-admin` services, Dex will serve a login page for users. This page
is customized via the `web` directory. This is essentially a copy of the `web` dir from the Dex project, per their
[instructions](https://dexidp.io/docs/templates/#using-your-own-templates).

We have removed the template `password.html` and it is instead represented by `dex/password-captcha.html` and `dex/password.html`.

To enable the Google reCAPTCHA widget, the env var `RECAPTCHA_SITE_KEY` must be set. If this is set, then when `make config` is run,
`dex/password-captcha.html` will be copied into the working directory at `~/.hoss/auth/web/templates` to be served. 
If not, `dex/password.html` will be used.

Similarly, if `RECAPTCHA_SITE_KEY` is set `config-dex-captcha.tmpl` will be rendered and appended to the dex config file at
`~/.hoss/auth/config-dex.yaml`. If not, `config-dex-no-captcha.tmpl` is rendered and appened to the dex config file.

To obtain your site key:

1. Go to [https://www.google.com/recaptcha/admin/create](https://www.google.com/recaptcha/admin/create)
2. Log in
3. Fill out the form, selecting the reCAPTCHA v2 option
4. Enter the domain you expect to use. If developing, be sure to include localhost
5. Set `RECAPTCHA_SITE_KEY` and `RECAPTCHA_SECRET_KEY` to the provided values
