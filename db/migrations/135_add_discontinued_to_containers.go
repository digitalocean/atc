package migrations

import "github.com/concourse/atc/dbng/migration"

func AddDiscontinuedToContainers(tx migration.LimitedTx) error {
	_, err := tx.Exec(`
  ALTER TABLE containers
  ADD COLUMN discontinued bool NOT NULL DEFAULT false
`)
	if err != nil {
		return err
	}
	return nil
}
