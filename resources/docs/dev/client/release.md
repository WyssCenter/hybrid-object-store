# Client Release Process

The `hoss-client` library is managed in GitHub and available on PYPI. To cut a new release:

1. Make sure the library version has been properly incremented. Developers should be incrementing the version as they merge PRs, but before a release always double check this is correct. To set the version, edit `hoss/version.py`. 
2. Tag main using the same version. Use the GitHub Release tool to create a new release with automatically generated changelog information.
3. GitHub Actions will build and push a release to PYPI for you.
4. ReadTheDocs will automatically update the "stable" docs.