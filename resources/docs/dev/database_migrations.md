# Database Migrations
The Hoss makes use of the [go-pg](https://github.com/go-pg/pg) Golang ORM for interacting with the Postgres database server. Each service (currently `core` and `auth`) has its own database within the server and maintains its own database schema and relational mappings.

To handle schema changes during updates the [go-pg migrations](https://github.com/go-pg/migrations) library is used.

## Migration Format
Migration files are straight forward but require a specific format.

Requirements:
1. Migration files must be placed under `pkg/database/migrations/` folder
2. Migration files must be called `XXXX_short_name.go`
   - Where `XXXX` is an incrementing number denoting the order in which the migration files will be applied.
   - Where `short_name` is a name comprised of up to several words that can be used to describe what the migration does, so that it is easier for developers to understand what should be going on with the migrations in that file.
3. Migration files must consist of a function called `RegisterXXXX` that calls `migrations.MustRegisterTx()` to register a new set of schema changes.
   - Where `XXXX` is the same number as used in the filename above
   - The call to `migrations.MustRegisterTx()` takes either one or two functions. If two functions are provided the first is called when upgrading to the current version and the second function is called when downgrading from the current version. If only one function is provided there is no downgrade function.
   - All calls to the database within the upgrade or downgrade function are wrapped in a single database transaction, meaning that if there is a problem applying the migration the database schema is left unchanged.
   - Because the `go-pg` ORM is used all migrations are written in Postgresql dialect and may make use of Postgres specific features.
   - As helpful feedback migrations normally print information about what they are about to do, especially if there are multiple steps to the migration. This makes it easier for a developer to debug any Postgresql errors that are encountered.

Once the migration file has been created the final step is to hook the migration into the service. This is accomplished by updating `pkg/database/database.go:Load()` to add a call to `hossMigrations.RegisterXXXX()` after the call to `hossMigrations.Register0001()`.

> Note: Automatic migration file detection is not used and if this final step is not done then the new migrations will not be run.

## Initial Migration
Each service has at least one migration file, called `0001_initial.go`, that contains the initial schema for the service. This typically consists of `CREATE` statements to setup the tables and other resources needed for the service to function. It may also include populating initial data.

## Additional Migrations
Each service may have additional migration files or you may be creating one. These migrations should contain SQL statements that either contain `CREATE` statements to create new resources that are need for new service features, `ALTER` statements to change existing resources to support code changes, or `DROP` statements to remove resources that are no longer used.

> Please keep in mind that during an update there may be data in the database tables and it should not be lost. If needed used a temporary table to store the existing data if you need to drop and recreate a table.

## Unit Testing
Unit tests make use of the migration files and will run the downgrade functions at the end of each test to reset the database to a clean state. This is a good check that the migration downgrade function correctly cleans up what the upgrade function does.
