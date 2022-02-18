import click

from hossadm.console import console
from hossadm import opensearch


@click.command()
@click.option('--version', '-v', type=str, required=True,
              help="Server version to which metadata index should be updated, in the form #.#.#")
@click.pass_context
def update_index(ctx, version):
    """Update the metadata index mapping to a new version"""
    console.print(f"\nUpdating metadata index to version v{version}")

    # migrate index
    opensearch.update_index(version)

