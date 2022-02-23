# Developer Documentation
The resources in this directory document development processes, software design, and other useful information for maintaining the Hoss.

## Backend Development

* [API documentation](server/api_docs.md) for the core and auth services are available when running the servers in "dev" mode.
* All changes to the database schema require migration support as outlined [in this document](server/database_migrations.md)
* We use DexIDP for OIDC provider federation. We build a modified version to add recaptcha support. Some consideration is required when updating the Dex container version as [described in this document](server/dex.md).
* The sync service implements this [Sync Policy Spec](server/sync_policy.md)
* Hoss server dev instructions can be found in the [README](../../../README.md)
* Hoss integration tests are critical to ensuring no regressions are introduced. You should **always** add integration tests when developing new features or fixing bugs. The integration test framework is located in the "test" directory and more information can be found in the [README](../../../test/README.md)
* [Docker Compose development](server/docker-compose.md)
* [Makefile development](server/makefile.md)

## Frontend Development
* [UI Sketch design file](ui/hoss.sketch) - This sketch file contains design elements for the Hoss web UI
* [Logo Sketch design file](ui/hoss-logo.sketch) - This sketch file contains the Hoss logo
* [Frontend development](ui/development.md)
* [Frontend architecture](ui/architecture.md)
* [Frontend styling](ui/styling.md)

## Release Processes
* [Hoss Server Release Process](server/server_release.md)
* [Client Library Release Process](client/release.md)
* [Terraform Module Release Process](server/terraform_release.md)

## System Design
* [Goals and Motivation](design/goals_and_motivation.md)
* [High Level System Design](design/system_arch.md)
* [Backup and Restore Design notes](server/backup_restore.md)
* [Metadata Search Design](design/search.md)
