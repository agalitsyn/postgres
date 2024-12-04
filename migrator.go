package postgres

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"

	"github.com/jackc/tern/v2/migrate"
)

const migrationsTableName = "schema_version"

func logMigrationProcess(version int32, name, direction, _ string) {
	slog.Info("migrating", "version", version, "name", name, "direction", direction)
}

func MigrateUp(ctx context.Context, db *DB, migrationsFiles fs.FS) error {
	conn, err := db.Pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("could not acquire connection: %w", err)
	}
	defer conn.Release()

	m, err := migrate.NewMigrator(ctx, conn.Conn(), migrationsTableName)
	if err != nil {
		return fmt.Errorf("could not create migrator: %w", err)
	}

	if err = m.LoadMigrations(migrationsFiles); err != nil {
		return fmt.Errorf("could not load migrations: %w", err)
	}

	m.OnStart = logMigrationProcess
	if err = m.Migrate(ctx); err != nil {
		return fmt.Errorf("could not migrate: %w", err)
	}

	return nil
}
