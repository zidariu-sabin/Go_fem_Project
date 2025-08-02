package store

import (
	"database/sql"
	"fmt"
	"io/fs"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/zidariu-sabin/femProject/internal/utils"
)

// var collection = "pg"

func Open() (*sql.DB, error) {
	db, err := sql.Open("pgx", utils.DATABASE_PGS_CONN_STRING)
	if err != nil {
		//format specifier that wraps error
		return nil, fmt.Errorf("db: open %w", err)
	}

	fmt.Println("Connected to Database... ")

	return db, nil
}

// db, migration file structure,
// we have to set the file structure so our migrations module(goose) can know where to pick migrations from
func MigrateFs(db *sql.DB, migrationsFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationsFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)

}

// applying migrations to the database using Up function and checking for errors
func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	err = goose.Up(db, dir)
	//we check if the database is up
	if err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	return nil
}
