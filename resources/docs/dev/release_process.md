# Making a Server Release

First, make sure the version number is properly incremented. In the file `VERSION`, edit the string to the desired number. The version number should be incremented by developers when merging PRs that make significant changes, especially to the API. When cutting a release, double check that the version has been incremented properly.

The major version should only be incremented for large breaking changes or when moving from a beta to stable product.

The minor version should be incremented when significant changes are made, and in particular anything that will break the backup/restore process (e.g. db schema change, file location change). Backup/restore compatibility is enforced by making sure the backup and server have the same `major.minor` version. Most changes can likely just increment the `build` version. 

Additionally, if changes are made to the API that will impact the client library, you must update the version support variables in the library (and make associated changes). To set the server version compatibility for a version of the client library, edit the `MIN_SUPPORTED_SERVER_VERSION` and `MAX_SUPPORTED_SERVER_VERSION` variables in [`hoss/api.py`](https://github.com/gigantum/hoss-client/blob/main/hoss/api.py). For the client library to work with a server, the server's version must be `MIN_SUPPORTED_SERVER_VERSION` <= `CURRENT_SERVER_VERSION` <= `MAX_SUPPORTED_SERVER_VERSION`. When making this change, you should bump the client version and cut a release along-side the server release.

To cut a release, we simply use GitHub releases and the auto-populating changelog functionality. Our releases are "soft" in that we don't actually push any artifacts anywhere yet. The 
user is still expected to pull the repo at the desired version and run `make env`, `make config`, `make build`, `make up`.
