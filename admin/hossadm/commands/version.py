import click

import hossadm
from hossadm.console import console
from hossadm.version import get_server_version


@click.command()
@click.option('--endpoint', '-e', type=str, default="http://localhost", show_default=True,
              help="Server root endpoint, including the scheme (e.g. https://hoss.mydomain.com")
def version(endpoint):
    """Print version info

    /f
    Returns:
        None
    """
    console.print(f"\nhossadm: v{hossadm.__version__}")
    server_version = get_server_version(endpoint)
    console.print(f"Hoss server: v{server_version['version']} (build {server_version['build']})\n")
