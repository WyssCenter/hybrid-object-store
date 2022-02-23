import click

from hossadm.console import console
from hossadm import user


@click.command()
@click.argument('username', type=str, nargs=1)
@click.option('--endpoint', '-e', type=str, default="http://localhost", show_default=True,
              help="Server root endpoint, including the scheme (e.g. https://hoss.mydomain.com")
@click.pass_context
def remove_user(ctx, username, endpoint):
    """Remove a user's access from a server

    This command is used to remove a user from the system after they have been removed from the external auth provider.
    This step is required because even if a user can no longer log in, the could in theory still exchange a PAT for
    valid credentials. Running this command will delete a user's PATs and group memberships.

    To run this command, you must set the environment variable `HOSS_PAT` to a valid PAT from a user with the `admin` role.
    """
    user.remove_user(username, endpoint)
