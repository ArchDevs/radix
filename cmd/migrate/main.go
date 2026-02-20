package main

import (
	"fmt"
	"log"
	"os"

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
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all up migrations",
	Run: func(cmd *cobra.Command, args []string) {
		m := createMigrate()
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to run up migrations: %v", err)
		}
		fmt.Println("✓ Migrations applied successfully")
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		m := createMigrate()
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to run down migrations: %v", err)
		}
		fmt.Println("✓ Migrations rolled back successfully")
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
