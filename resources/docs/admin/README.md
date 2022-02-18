# Admin Documentation for the Hoss

This directory contains documentation for administrators working with the Hoss. It is 
automatically built and published on Read the Docs. Tagging a release on GitHub
will update the "stable" version automatically.

## Editing

To edit the admin docs for the first time:

1. Create a new virtualenv in this directory
   ```
   python3 -m venv ./docs-venv
   source docs-venv/bin/activate 
   ```
2. Install python dependencies
   ```
   pip install -r ./requirements.txt
   ```
3. Make desired edits
4. Run `make html` to generate html documentation locally
5. Open `resources/docs/admin/build/html/index.html` to preview changes