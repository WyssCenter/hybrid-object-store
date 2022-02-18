import click

from hossadm.commands.backup import backup
from hossadm.commands.restore import restore
from hossadm.commands.version import version
from hossadm.commands.update_index import update_index

CONTEXT_SETTINGS = dict(help_option_names=['-h', '--help'])


@click.group(help="A Command Line Interface to administer a Hoss server.",
             context_settings=CONTEXT_SETTINGS)
def cli():
    pass


# Add commands from package
cli.add_command(backup)
cli.add_command(restore)
cli.add_command(version)
cli.add_command(update_index)


if __name__ == '__main__':
    cli()
