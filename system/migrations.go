package system

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/GeertJohan/go.rice"
	"github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	"github.com/rubenv/sql-migrate"
)

//RunMigrations applies database migrations, where box - rice box for "migrations" dir, command:
//new - creates new blank migration in "migrations" directory. Edit that file as needed.
//"up", "down"- apply all pending migrations, or undo the last one
//"redo" - rollback last migration, then reapply it
//db - database handler
func RunMigrations(box *rice.Box, db *sqlx.DB, command *string) {
	switch *command {
	case "new":
		migrateNew(box)
	case "up":
		migrateUp(db.DB, box, 0)
	case "down":
		migrateDown(db.DB, box, 1)
	case "redo":
		migrateDown(db.DB, box, 1)
		migrateUp(db.DB, box, 1)
	default:
		logrus.Fatalf("Wrong migration flag %q, acceptable values: up, down", *command)
	}
}

//migrateNew creates new blank migration
func migrateNew(box *rice.Box) {
	if len(flag.Args()) == 0 {
		logrus.Error("Migration's name not specified")
		return
	}
	name := path.Join(box.Name(), fmt.Sprintf("%d_%s.sql", time.Now().Unix(), flag.Arg(0)))
	file, err := os.Create(name)
	if err != nil {
		logrus.Error(err)
		return
	}
	fmt.Fprintf(file, "-- +migrate Up\n")
	fmt.Fprintf(file, "-- SQL in section 'Up' is executed when this migration is applied\n\n\n")
	fmt.Fprintf(file, "-- +migrate Down\n")
	fmt.Fprintf(file, "-- SQL in section 'Down' is executed when this migration is rolled back\n\n\n")
	err = file.Close()
	if err != nil {
		logrus.Error(err)
	} else {
		logrus.Infof("File %s has been successfully created\n", name)
	}
}

//migrateUp applies {{max}} pending db migrations. If max == 0, it applies all
func migrateUp(DB *sql.DB, box *rice.Box, max int) {
	migrations := getRiceMigrations(box)
	n, err := migrate.ExecMax(DB, "postgres", migrations, migrate.Up, max)
	if err != nil {
		logrus.Error(err)
	} else {
		logrus.Infof("%d migration(s) applied", n)
	}
}

//migrateDown rolls back {{max}} db migrations. If max == 0, it rolles back all of them
func migrateDown(DB *sql.DB, box *rice.Box, max int) {
	migrations := getRiceMigrations(box)
	n, err := migrate.ExecMax(DB, "postgres", migrations, migrate.Down, max)
	if err != nil {
		logrus.Error(err)
	} else {
		logrus.Infof("%d migration(s) rolled back", n)
	}
}

//getRiceMigrations builds migration source from go.rice storage
func getRiceMigrations(box *rice.Box) *migrate.MemoryMigrationSource {
	source := &migrate.MemoryMigrationSource{}
	fn := func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql") {
			migFile, err := box.Open(path)
			if err != nil {
				return err
			}
			mig, err := migrate.ParseMigration(path, migFile)
			migFile.Close()
			if err != nil {
				return err
			}
			source.Migrations = append(source.Migrations, mig)
		}
		return nil
	}
	err := box.Walk("", fn)
	if err != nil {
		logrus.Fatal(err)
		return nil
	}
	return source
}
