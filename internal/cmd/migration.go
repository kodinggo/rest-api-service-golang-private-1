package cmd

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"
)

var (
	direction string
	step      int
)

var migrationCMD = &cobra.Command{
	Use:   "migrate",
	Short: "migrate is a command to migrate database schema",
	Run: func(cmd *cobra.Command, args []string) {
		// db environment
		var (
			dbUsername = "root"
			dbPassword = ""
			dbName     = "kodinggo"
			dbHost     = "localhost:3306"
		)

		// prepare connection string
		// charset=utf8mb4 agar dapat menyimpan semua karakter unicode
		// parseTime=true agar dapat diparsing dari timestamp ke tipe time.Time
		connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true",
			dbUsername,
			dbPassword,
			dbHost,
			dbName)
		connDB, err := sql.Open("mysql", connStr)
		if err != nil {
			log.Panicf("failed to connect server db, error: %s", err.Error())
		}

		// Perform migration
		migrations := &migrate.FileMigrationSource{Dir: "./db/migrations"}
		var n int // "n" berapa migration yg di-applied
		if direction == "down" {
			n, err = migrate.ExecMax(connDB, "mysql", migrations, migrate.Down, step)
		} else {
			n, err = migrate.ExecMax(connDB, "mysql", migrations, migrate.Up, step)
		}
		if err != nil {
			log.Panicf("failed to run migration, error: %s", err.Error())
		}

		log.Printf("successfully applied %d migration(s)", n)
	},
}

func init() {
	migrationCMD.Flags().StringVarP(&direction, "direction", "d", "up", "apply/rollback migrations")
	migrationCMD.Flags().IntVarP(&step, "step", "s", 1, "number of step to undo/redo migrations")
	rootCMD.AddCommand(migrationCMD)
}
