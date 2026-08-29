package migrator

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	appmigrations "github.com/noirbyss/worktrition-app/backend/user-service/migrations"
)

const ensureSchemaMigrationsTableQuery = `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version TEXT PRIMARY KEY,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	)
`

func Run(ctx context.Context, pool *pgxpool.Pool) error {
	if pool == nil {
		return fmt.Errorf("migration pool is required")
	}

	if err := ensureSchemaMigrationsTable(ctx, pool); err != nil {
		return err
	}

	appliedVersions, err := loadAppliedVersions(ctx, pool)
	if err != nil {
		return err
	}

	fileNames, err := migrationFileNames()
	if err != nil {
		return err
	}

	appliedCount := 0
	for _, fileName := range fileNames {
		version := strings.TrimSuffix(fileName, ".up.sql")
		if appliedVersions[version] {
			continue
		}

		if err := applyMigration(ctx, pool, version, fileName); err != nil {
			return err
		}

		appliedCount++
	}

	if appliedCount == 0 {
		slog.Info("database migrations already up to date")
		return nil
	}

	slog.Info("database migrations applied", "count", appliedCount)
	return nil
}

func ensureSchemaMigrationsTable(ctx context.Context, pool *pgxpool.Pool) error {
	if _, err := pool.Exec(ctx, ensureSchemaMigrationsTableQuery); err != nil {
		return fmt.Errorf("ensure schema_migrations table: %w", err)
	}

	return nil
}

func loadAppliedVersions(ctx context.Context, pool *pgxpool.Pool) (map[string]bool, error) {
	rows, err := pool.Query(ctx, `SELECT version FROM schema_migrations`)
	if err != nil {
		return nil, fmt.Errorf("load applied migrations: %w", err)
	}
	defer rows.Close()

	appliedVersions := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("scan applied migration: %w", err)
		}

		appliedVersions[version] = true
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate applied migrations: %w", err)
	}

	return appliedVersions, nil
}

func migrationFileNames() ([]string, error) {
	entries, err := appmigrations.Files.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("read migration directory: %w", err)
	}

	fileNames := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.HasSuffix(name, ".up.sql") {
			fileNames = append(fileNames, name)
		}
	}

	slices.Sort(fileNames)

	return fileNames, nil
}

func applyMigration(ctx context.Context, pool *pgxpool.Pool, version, fileName string) error {
	script, err := appmigrations.Files.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("read migration %s: %w", version, err)
	}

	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin migration %s: %w", version, err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if _, err := tx.Conn().PgConn().Exec(ctx, string(script)).ReadAll(); err != nil {
		return fmt.Errorf("execute migration %s: %w", version, err)
	}

	if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations (version) VALUES ($1)`, version); err != nil {
		return fmt.Errorf("mark migration %s as applied: %w", version, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit migration %s: %w", version, err)
	}

	slog.Info("migration applied", "version", version)

	return nil
}
