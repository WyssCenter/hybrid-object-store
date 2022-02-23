# Revoking Access

Because the Hoss system is built around auth integration, it is possible for a user to be removed from the authentication provider, yet still have access to the Hoss via PATs. Because of this, there is a tool and recommended process for revoking access to a Hoss server.

1) Remove user's access from the Authentication provider. If this is the internal LDAP server, this is referring to deleting the user from the server. If it is some other provider (e.g. Azure AD) this could mean deleting the user, deactivating the user, or removing a group that granted them access to the server. Regardless of the auth provider you have configured, the user should not be able to successfully log into the Hoss.
2) Use the `hossadm` library to remove the user's PATs and group memberships. In the example below, the server is running at `https://hoss.mycompany.com` and the user we wish to deactivate is `user1`

   ```
   hossadm remove-user --endpoint https://hoss.mycompany.com user1
   ```

3) Wait X hours, where X is the JWT expiration time set in your auth service. After this time, you can be guaranteed that the user will no longer be able to access the system in any way.


## Installing the `hossadm` tool

The `hossadm` tool should typically be used **on the server**. If you have yet to install the `hossadm` tool:

1. Create and activate a new Python3 virtual environment.
   * For example, run `python3 -m venv ./hossadm-venv` in your home directory
   * Then run `source ~/hossadm-venv/bin/activate` 
2. From the `admin/` directory of the Hoss source code repository run `pip3 install -U .`

If the tool has been updated, simply run `pip3 install -U .` again after updating the Hoss code repository to
the desired version.