package migrations

import (
	"fmt"

	"github.com/go-pg/migrations/v8"
)

func Register0002() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("Creating enum dataset_delete_status...")
		_, err := db.Exec(`CREATE TYPE dataset_delete_status AS ENUM ('NOT_SCHEDULED','SCHEDULED', 'IN_PROGRESS', 'ERROR')`)
		if err != nil {
			return err
		}

		fmt.Println("Adding 'delete_on' and 'delete_status' columns to datasets...")
		_, err = db.Exec(`ALTER TABLE datasets
			ADD COLUMN delete_on timestamp,
			ADD COLUMN delete_status dataset_delete_status;
			UPDATE datasets SET delete_status = 'NOT_SCHEDULED';
			ALTER TABLE datasets ALTER COLUMN delete_status SET NOT NULL;
			ALTER TABLE datasets ALTER COLUMN delete_status SET DEFAULT 'NOT_SCHEDULED';`)
		if err != nil {
			return err
		}

		return nil
	}, func(db migrations.DB) error {
		fmt.Println("Removing 'delete_on' and 'delete_status' columns from datasets...")
		_, err := db.Exec(`ALTER TABLE datasets
			DROP COLUMN delete_on,
			DROP COLUMN delete_status;`)
		if err != nil {
			return err
		}

		fmt.Println("Dropping enum dataset_delete_status...")
		_, err = db.Exec(`DROP TYPE dataset_delete_status`)
		if err != nil {
			return err
		}

		return nil
	})
}
