# Using Microsoft Azure Active Directory

You can easily use Microsoft Azure Active Directory (AD) to provide authentication services. This is done by configuring the Hoss server as an OIDC client and
using the Dex [Microsoft connector](https://dexidp.io/docs/connectors/microsoft/)

## Create App Registration
First you must create an App registration to use with the Hoss Server. 

* Navigate in the Azure web portal to your AD.
* In the left navigation panel, select "App registrations".
* Click the "New registration" button at the top of the panel that loads
* Fill out the form to create the app registration
  * Set the `Name` to something reasonable (e.g. "Hoss Server")
  * For supported account types, select "Accounts in this organizational directory only (Default Directory only - Single tenant)". This will limit logins to Microsoft users only in your AD.
  * For the "Redirect URI", select the "Web" platform type and `<EXTERNAL_HOSTNAME>/dex/callback`, where `<EXTERNAL_HOSTNAME>` is the externally accessible hostname for your Hoss server, including the scheme (e.g. https://hoss.myserver.com/dex/callback).

Keep this page open, as you'll need to reference it later when configuring your Hoss Server.

## Configure Groups
Depending on your desired configuration there are 2 or 3 groups that you must have defined in your AD. These can be existing groups, or you can create new groups specifically for the purpose
of Hoss integration. You will configure the Hoss server later, indicating the name of the groups to use. Ideally, users should only be in 1 group.


### Admin Group

The first group to create (or select) is the administrator group. Users that are in this group will be granted the `admin` role in the Hoss. This will let them view, create, and delete all resources in the Hoss. You should typically be selective on users with this role, or even enforce multiple accounts for admins, where they also have an account with reduced privileges for normal
day-to-day use of the Hoss.

After creating or selecting a group, make note of the name for later.

### Privileged Group

The second group to create (or select) is the privileged group, which typically is used by developers who need the ability to create datasets. Users that are in this group will be granted the `privileged` role in the Hoss. This will grant them access to additional capabilities, such as creating & deleting datasets, configuring dataset syncing, and managing groups of which they are a member. 

After creating or selecting a group, make note of the name for later.

### (Optional) User Group

The last group to create (or select) is an optional group for all "other" users. If you choose to do this, you can configure dex to require that users must be in one of the three groups. This can
be useful in situations where you have a large organization and don't want to grant access to the Hoss for *everyone* in your AD. Users in this group will be allowed to log in, but will have
the standard `user` role and no additional capabilities. They must be granted permissions to datasets to interact with any data.

After creating or selecting a group, make note of the name for later.

## Configure Users
Next, make sure to place users into these new groups as needed. If you are using the optional User Group, you'll have to make sure everyone is in one group.


## Configure Hoss Server
Configure the server as you normally would with the following additional steps:

Edit `~/.hoss/auth/config.yaml` and:
* set `issuer` to `<EXTERNAL_HOSTNAME>/dex`, where `<EXTERNAL_HOSTNAME>` is the externally accessible hostname for your Hoss server, including the scheme (e.g. https://hoss.myserver.com/dex).
* set `username_claim` to email. This will tell the Hoss to take the first part of a user's email address and use it as their unique username within the system.
* set `admin_group` to the name of your administrators group from above
* set `privileged_group` to the name of your privileged group from above


Edit `~/.hoss/auth/config-dex.yaml` and:
* set `issuer` to `<EXTERNAL_HOSTNAME>/dex`, where `<EXTERNAL_HOSTNAME>` is the externally accessible hostname for your Hoss server, including the scheme (e.g. https://hoss.myserver.com/dex).
* Under the connectors section, remove the entire entry for the internal LDAP server. You can keep both connectors in place if you wish to have users in both systems access the Hoss server, but this isn't a common configuration.
* Add the Microsoft connector, which would look something like this:
  
  ```yaml
  connectors:
  - type: microsoft
    # Required field for connector id.
    id: microsoft
    # Required field for connector name.
    name: Microsoft
    config:
        # Credentials can be string literals or pulled from the environment.
        clientID: <CLIENT_ID>
        clientSecret: <CLIENT_SECRET>
        redirectURI: <CALLBACK>
        tenant: <TENANT>
        emailToLowercase: true
        groups:
          - <ADMIN_GROUP>
          - <PRIVILEGED_GROUP>
          - <USER_GROUP>
  ```

  where

  * `<CLIENT_ID>` is the "Application (client) ID" of the App registration. You can find this value on the Overview section of the App registration in the Azure web portal.
  * `<CLIENT_SECRET>` is the client secret of the App registration. You must generate a client secret for the App registration if you have not done so. Open the "Certificates & secrets" section of your App registration. Click on the "New client secret" button. The secret will only be displayed at this time, so be sure to copy it!
  * `<CALLBACK>` is `<EXTERNAL_HOSTNAME>/dex/callback`, where `<EXTERNAL_HOSTNAME>` is the externally accessible hostname for your Hoss server, including the scheme (e.g. https://hoss.myserver.com/dex/callback).
  * `<TENANT>` is the "Directory (tenant) ID" of the AD associated with this App registration. You can find this value on the Overview section of the App registration in the Azure web portal.
  * `<ADMIN_GROUP>` is the name of your administrators group from above
  * `<PRIVILEGED_GROUP>` is the name of your privileged group from above
  * `<USER_GROUP>` is the name of your user group from above

Note, if you are not using a user group, and have only defined admin and privileged groups, you should remove the entire `group` section. In this configuration all users in your AD will be able to log into the Hoss.
