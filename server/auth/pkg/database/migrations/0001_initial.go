package migrations

import (
	"fmt"

	"github.com/go-pg/migrations/v8"
)

func Register0001() {
	migrations.MustRegisterTx(func(db migrations.DB) error {

		fmt.Println("Creating enum role...")
		_, err := db.Exec(`CREATE TYPE role AS ENUM ('admin', 'privileged', 'user')`)
		if err != nil {
			return err
		}

		fmt.Println("Creating table users...")
		_, err = db.Exec(`CREATE TABLE users (
			id bigserial PRIMARY KEY,
			username character varying NOT NULL,
			full_name character varying,
			given_name character varying,
			family_name character varying,
			email character varying,
			email_verified boolean NOT NULL,
			subject character varying,
			role role,
			UNIQUE (username)
		)`)
		if err != nil {
			return err
		}

		fmt.Println("Creating table groups...")
		_, err = db.Exec(`CREATE TABLE groups (
			id bigserial PRIMARY KEY,
			group_name character varying(64) NOT NULL,
			description character varying NOT NULL,
			UNIQUE (group_name)
		)`)
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

		// https://stackoverflow.com/questions/41970461/how-to-generate-a-random-unique-alphanumeric-id-of-length-n-in-postgres-9-6
		// The token prefix 'hp_' is used to identify the token in the Authorization header and provide support for adding
		// future capabilities to tokens (i.e. the system would be able to infer the version of the token format)
		fmt.Println("Creating PAT generator function...")
		_, err = db.Exec(`CREATE EXTENSION pgcrypto`)
		if err != nil {
			return err
		}
		_, err = db.Exec(`CREATE OR REPLACE FUNCTION generate_uid(size INT) RETURNS TEXT AS $$
		DECLARE
		  characters TEXT := 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
		  bytes BYTEA := gen_random_bytes(size);
		  l INT := length(characters);
		  i INT := 0;
		  output TEXT := 'hp_';
		BEGIN
		  WHILE i < size LOOP
			output := output || substr(characters, get_byte(bytes, i) % l + 1, 1);
			i := i + 1;
		  END LOOP;
		  RETURN output;
		END;
		$$ LANGUAGE plpgsql VOLATILE;`)
		if err != nil {
			return err
		}

		fmt.Println("Creating table personal_access_tokens...")
		_, err = db.Exec(`CREATE TABLE personal_access_tokens (
			id bigserial PRIMARY KEY,
			pat character varying NOT NULL DEFAULT generate_uid(30),
			owner_id bigint REFERENCES users ON DELETE SET NULL,
			description character varying,
			UNIQUE (pat)
		)`)
		if err != nil {
			return err
		}

		return nil
	}, func(db migrations.DB) error {
		fmt.Println("Dropping table personal_access_tokens...")
		_, err := db.Exec(`DROP TABLE personal_access_tokens`)
		if err != nil {
			return err
		}

		fmt.Println("Dropping table memberships...")
		_, err = db.Exec(`DROP TABLE memberships`)
		if err != nil {
			return err
		}

		fmt.Println("Dropping table groups...")
		_, err = db.Exec(`DROP TABLE groups`)
		if err != nil {
			return err
		}

		fmt.Println("Dropping table users...")
		_, err = db.Exec(`DROP TABLE users`)
		if err != nil {
			return err
		}

		fmt.Println("Dropping enum role...")
		_, err = db.Exec(`DROP TYPE role`)
		if err != nil {
			return err
		}

		_, err = db.Exec(`DROP EXTENSION pgcrypto CASCADE`)
		if err != nil {
			return err
		}

		return nil
	})
}
