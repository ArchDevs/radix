package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/ArchDevs/radix/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration tool",
	Long:  `A CLI tool for managing database migrations using golang-migrate`,
}

func init() {
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(forceCmd)
}

var steps int

var upCmd = &cobra.Command{
	Use:   "up [steps]",
	Short: "Apply up migrations (all or specified steps)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		m := createMigrate()

		if len(args) > 0 {
			n, err := strconv.Atoi(args[0])
			if err != nil {
				log.Fatalf("Invalid step count: %v", err)
			}
			if err := m.Steps(n); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("Failed to run up migrations: %v", err)
			}
			fmt.Printf("✓ Applied %d migration(s) successfully\n", n)
		} else {
			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("Failed to run up migrations: %v", err)
			}
			fmt.Println("✓ All migrations applied successfully")
		}
	},
}

var downCmd = &cobra.Command{
	Use:   "down [steps]",
	Short: "Rollback migrations (all or specified steps)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		m := createMigrate()

		if len(args) > 0 {
			n, err := strconv.Atoi(args[0])
			if err != nil {
				log.Fatalf("Invalid step count: %v", err)
			}
			if err := m.Steps(-n); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("Failed to run down migrations: %v", err)
			}
			fmt.Printf("✓ Rolled back %d migration(s) successfully\n", n)
		} else {
			if err := m.Down(); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("Failed to run down migrations: %v", err)
			}
			fmt.Println("✓ All migrations rolled back successfully")
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Check current migration version",
	Run: func(cmd *cobra.Command, args []string) {
		m := createMigrate()
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("Failed to get version: %v", err)
		}
		fmt.Printf("Version: %d, Dirty: %v\n", version, dirty)
	},
}

var forceCmd = &cobra.Command{
	Use:   "force [version]",
	Short: "Force set migration version (for dirty state recovery)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		m := createMigrate()
		version, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalf("Invalid version: %v", err)
		}
		if err := m.Force(version); err != nil {
			log.Fatalf("Failed to force version: %v", err)
		}
		fmt.Printf("✓ Forced migration version to: %d\n", version)
	},
}

func createMigrate() *migrate.Migrate {
	_ = godotenv.Load()
	cfg := config.Load()

	m, err := migrate.New(
		"file://migrations",
		"sqlite3://"+cfg.DB.DataSource,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	return m
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
