package migrations

import (
	"fmt"

	"github.com/go-pg/migrations/v8"
)

func Register0001() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("Creating table users...")
		_, err := db.Exec(`CREATE TABLE users (
			id bigserial PRIMARY KEY,
			username character varying NOT NULL,
			UNIQUE (username)
		)`)
		if err != nil {
			return err
		}

		fmt.Println("Creating table groups...")
		_, err = db.Exec(`CREATE TABLE groups (
			id bigserial PRIMARY KEY,
			group_name character varying NOT NULL,
			UNIQUE (group_name)
		)`)
		if err != nil {
			return err
		}

		fmt.Println("Creating enum objectstoretype...")
		_, err = db.Exec(`CREATE TYPE objectstoretype AS ENUM ('minio', 's3')`)
		if err != nil {
			return err
		}

		fmt.Println("Creating table memberships...")
		_, err = db.Exec(`CREATE TABLE memberships (
			group_id bigint REFERENCES groups ON DELETE CASCADE NOT NULL,
			user_id bigint REFERENCES users ON DELETE CASCADE NOT NULL,
			PRIMARY KEY (group_id, user_id),
			UNIQUE (group_id, user_id)
		)`)
		if err != nil {
			return err
		}

		fmt.Println("Creating table object_stores...")
		_, err = db.Exec(`CREATE TABLE object_stores (
			id bigserial PRIMARY KEY,
			name character varying NOT NULL,
			description character varying,
			endpoint character varying NOT NULL,
			object_store_type objectstoretype NOT NULL,
			region character varying,
			profile character varying,
			role_arn character varying,
			notification_arn character varying,
			UNIQUE (name)
		)`)
		if err != nil {
			return err
		}

		fmt.Println("Creating table namespaces...")
		_, err = db.Exec(`CREATE TABLE namespaces (
			id bigserial PRIMARY KEY,
			name character varying NOT NULL,
			description character varying,
			object_store_id bigint REFERENCES object_stores ON DELETE CASCADE NOT NULL,
			bucket_name character varying NOT NULL,
			UNIQUE (name)
		)`)
		if err != nil {
			return err
		}

		fmt.Println("Creating table datasets...")
		_, err = db.Exec(`CREATE TABLE datasets (
			id bigserial PRIMARY KEY,
			namespace_id bigint REFERENCES namespaces ON DELETE CASCADE NOT NULL,
			name character varying NOT NULL,
			description character varying,
			owner_id bigint REFERENCES users ON DELETE SET NULL,
			created timestamp NOT NULL,
			root_directory character varying NOT NULL,
			sync_enabled boolean NOT NULL,
			sync_type character varying,
			UNIQUE (namespace_id, name)
		)`)
		if err != nil {
			return err
		}

		fmt.Println("Creating enum permission...")
		_, err = db.Exec(`CREATE TYPE permission AS ENUM ('r', 'rw')`)
		if err != nil {
			return err
		}

		fmt.Println("Creating table permissions...")
		_, err = db.Exec(`CREATE TABLE permissions (
			group_id bigint REFERENCES groups ON DELETE CASCADE NOT NULL,
			dataset_id bigint REFERENCES datasets ON DELETE CASCADE NOT NULL,
			permission permission NOT NULL,
			PRIMARY KEY (group_id, dataset_id),
			UNIQUE (group_id, dataset_id)
		)`)
		if err != nil {
			return err
		}

		fmt.Println("Creating table sync_configurations...")
		_, err = db.Exec(`CREATE TABLE sync_configurations (
			id bigserial PRIMARY KEY,
			source_namespace_id bigint REFERENCES namespaces ON DELETE CASCADE NOT NULL,
			target_core_service character varying NOT NULL,
			target_namespace character varying NOT NULL,
			sync_type character varying NOT NULL,
			UNIQUE (source_namespace_id, target_core_service, target_namespace)
		)`)
		if err != nil {
			return err
		}

		// The sync_configuration_meta table is, currently, just used to hold a timestamp
		// of when the last modification to either the sync_configurations table was made
		// or when the datasets.sync_enabled was updated or a dataset was deleted
		//
		// The sync_configuration_meta timestamp is updated using a stored procedure that is the
		// target of triggers defined on the sync_configurations and datasets tables
		//
		// The triggers are only triggered when data in the table is changed, so if the user
		// tries to delete a non-existent record (for example) the trigger will not fire
		//
		// The sync_configuration_meta.last_updated value is used when the Sync Service
		// queries for the latest Sync Configuration information. The Sync Service will
		// include the `If-Modified-Since` HTTP header, so that it will know if there has
		// been any change in the data. That header is compared against the last_updated
		// timestamp to determine if there has been a change in the data and that the Core
		// Service needs to query all of the Sync Configuration data and return it

		fmt.Println("Creating table sync_configuration_meta...")
		_, err = db.Exec(`CREATE TABLE sync_configuration_meta (
			id bigserial PRIMARY KEY,
			last_updated timestamp NOT NULL
		)`)
		if err != nil {
			return err
		}

		// Create the meta data defaults
		// Set the initial timestamp to the current time
		_, err = db.Exec(`INSERT INTO sync_configuration_meta (last_updated) VALUES (now())`)
		if err != nil {
			return err
		}

		// Create the stored procedure that will set the last_updated timestamp to the current
		// (database) time
		fmt.Println("Creating function sync_configuration_updated...")
		_, err = db.Exec(`CREATE FUNCTION sync_configuration_updated() RETURNS TRIGGER AS $sync_updated$
			BEGIN
				IF NEW IS DISTINCT FROM OLD THEN
					UPDATE sync_configuration_meta SET last_updated = now();
				END IF;
				return NULL;
			END
			$sync_updated$ LANGUAGE plpgsql
		`)
		if err != nil {
			return err
		}

		// Create the trigger on the sync_configurations table
		fmt.Println("Creating trigger sync_configuration_updated...")
		_, err = db.Exec(`CREATE TRIGGER sync_configuration_updated
			AFTER INSERT OR UPDATE OR DELETE ON sync_configurations
			FOR EACH ROW EXECUTE PROCEDURE sync_configuration_updated()
		`)
		if err != nil {
			return err
		}

		// Create the trigger on the datasets table
		fmt.Println("Creating trigger dataset_sync_updated...")
		// NOTE: trigger on all dataset deletes as the function doesn't check
		//       to see if the dataset has sync enabled
		_, err = db.Exec(`CREATE TRIGGER dataset_sync_updated
			AFTER UPDATE OF sync_enabled OR DELETE ON datasets
			FOR EACH ROW EXECUTE PROCEDURE sync_configuration_updated()
		`)
		if err != nil {
			return err
		}

		return nil
	}, func(db migrations.DB) error {
		fmt.Println("Dropping trigger dataset_sync_updated...")
		_, err := db.Exec(`DROP TRIGGER IF EXISTS dataset_sync_updated ON datasets`)
		if err != nil {
			return err
		}

		fmt.Println("Dropping trigger sync_configuration_updated...")
		_, err = db.Exec(`DROP TRIGGER IF EXISTS sync_configuration_updated ON sync_configurations`)
		if err != nil {
			return err
		}

		fmt.Println("Dropping function sync_configuration_updated...")
		_, err = db.Exec(`DROP FUNCTION IF EXISTS sync_configuration_updated()`)
		if err != nil {
			return err
		}

		fmt.Println("Dropping table sync_configuration_meta...")
		_, err = db.Exec(`DROP TABLE sync_configuration_meta`)
		if err != nil {
			return err
		}

		fmt.Println("Dropping table sync_configurations...")
		_, err = db.Exec(`DROP TABLE sync_configurations`)
		if err != nil {
			return err
		}

		fmt.Println("Dropping table permission...")
		_, err = db.Exec(`DROP TABLE permissions`)
		if err != nil {
			return err
		}

		fmt.Println("Dropping enum permission...")
		_, err = db.Exec(`DROP TYPE permission`)
		if err != nil {
			return err
		}

		fmt.Println("Dropping table memberships...")
		_, err = db.Exec(`DROP TABLE memberships`)
		if err != nil {
			return err
		}

		fmt.Println("Dropping table datasets...")
		_, err = db.Exec(`DROP TABLE datasets`)
		if err != nil {
			return err
		}

		fmt.Println("Dropping table namespaces...")
		_, err = db.Exec(`DROP TABLE namespaces`)
		if err != nil {
			return err
		}

		fmt.Println("Dropping table object_stores...")
		_, err = db.Exec(`DROP TABLE object_stores`)
		if err != nil {
			return err
		}

		fmt.Println("Dropping enum objectstoretype...")
		_, err = db.Exec(`DROP TYPE objectstoretype`)
		if err != nil {
			return err
		}

		fmt.Println("Dropping table users...")
		_, err = db.Exec(`DROP TABLE users`)
		if err != nil {
			return err
		}

		fmt.Println("Dropping table groups...")
		_, err = db.Exec(`DROP TABLE groups`)
		if err != nil {
			return err
		}

		return nil
	})
}
