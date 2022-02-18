package migrations

import (
	"fmt"

	"github.com/go-pg/migrations/v8"
)

func Register0003() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("Altering table datasets (adding column sync_policy) ...")
		_, err := db.Exec(`ALTER TABLE datasets
			ADD COLUMN sync_policy character varying
		`)
		if err != nil {
			return err
		}

		fmt.Println("Adding default sync_policy to table datasets...")
		_, err = db.Exec(`UPDATE datasets
			SET sync_policy = '{"Version":"1","Statements":[]}'
			WHERE sync_enabled
		`)
		if err != nil {
			return err
		}

		// Create the trigger on the datasets table - specific to the sync_policy field
		fmt.Println("Creating trigger dataset_sync_policy_updated...")
		_, err = db.Exec(`CREATE TRIGGER dataset_sync_policy_updated
			AFTER UPDATE OF sync_policy ON datasets
			FOR EACH ROW EXECUTE PROCEDURE sync_configuration_updated()
		`)
		if err != nil {
			return err
		}

		return nil
	}, func(db migrations.DB) error {
		fmt.Println("Dropping trigger dataset_sync_policy_updated...")
		_, err := db.Exec(`DROP TRIGGER IF EXISTS dataset_sync_policy_updated ON datasets`)
		if err != nil {
			return err
		}

		fmt.Println("Altering table datasets (dropping column sync_policy) ...")
		_, err = db.Exec(`ALTER TABLE datasets
			DROP COLUMN IF EXISTS sync_policy
		`)
		if err != nil {
			return err
		}

		return nil
	})
}
