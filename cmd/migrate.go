package cmd

import (
	"fmt"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/slbmax/ses-weather-app/assets"
	"github.com/slbmax/ses-weather-app/internal/config"
	"github.com/spf13/cobra"
	"gitlab.com/distributed_lab/kit/kv"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate [up|down]",
	Args:  cobra.ExactArgs(1),
	Short: "Run database migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		getter, err := kv.FromEnv()
		if err != nil {
			return fmt.Errorf("failed to get configer key-value getter: %w", err)
		}

		cfg := config.New(getter)

		var direction migrate.MigrationDirection
		switch args[0] {
		case "up":
			direction = migrate.Up
		case "down":
			direction = migrate.Down
		default:
			return fmt.Errorf("invalid migration type: %s", args[0])
		}

		migrationsFs := &migrate.EmbedFileSystemMigrationSource{
			FileSystem: assets.Migrations,
			Root:       "migrations",
		}

		applied, err := migrate.Exec(cfg.DB().RawDB(), "postgres", migrationsFs, direction)
		if err != nil {
			return fmt.Errorf("failed to run migrations: %w", err)
		}

		cfg.Log().
			WithField("direction", direction).
			WithField("applied", applied).
			Info("migrations applied")

		return nil
	},
}
