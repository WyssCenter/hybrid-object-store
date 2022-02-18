# User Interface
There are a few items that can be customized in the user interface by editing `~/.hoss/ui/config.json`. The fields are described below:

* `server_name`: The "name" of the server. This will be displayed to users in the top menu bar, and can help users locate themselves when switching between servers. This value defaults to the hostname of the server.
* `colors`:
  * `primary`: The hex value for the primary color of the app
  * `secondary`: The hex value for the secondary color of the app

Once changes are made to this file, simply refresh the page to see the effect. You do not need to restart the service.

You may also wish to change the logo and favicon. This can be done by replacing the `~/.hoss/ui/logo.svg` and `~/.hoss/ui/favicon.png` files. Currently you must keep the file types the same.

Note, this won't change the Dex login page if you are using the internal auth provider. If you wish to also modify this, you must manually do it in `~/.hoss/auth/web/themes/hoss`.