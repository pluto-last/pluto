package main

import (
	"flag"
	"log"
	"os"
	"pluto/cmd/migrate/table"
	"pluto/config"
	"pluto/middleware/db"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

var configPath string

func main() {
	flag.StringVar(&configPath, "config", "./conf/conf.yaml", "config file path")
	flag.Parse()
	db, err := DB(configPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	m := gormigrate.New(db, gormigrate.DefaultOptions,
		[]*gormigrate.Migration{
			initProject,
		},
	)
	if len(os.Args) < 3 {
		if err = m.Migrate(); err != nil {
			log.Fatalln(err)
		}
		log.Println("Migration OK")
		return
	}
	switch os.Args[2] {
	case "down":
		if err = m.RollbackLast(); err != nil {
			log.Fatalln(err)
		}
		log.Println("Rollback OK")
	default:
		if err = m.Migrate(); err != nil {
			log.Fatalln(err)
		}
		log.Println("Migration OK")
	}
}

var initProject = &gormigrate.Migration{
	ID: "202110141753",
	Migrate: func(db *gorm.DB) error {
		err := db.AutoMigrate(
			new(table.User),
			new(table.Role),
			new(table.RolePermissions),
			new(table.Permission),
			new(table.UserRoles),
			new(table.OperationRecord),
			new(table.Task),
			new(table.TaskRecord),
			new(table.TaskResult),
			new(table.Wallet),
			new(table.Robot),
			new(table.Scene),
			new(table.Sip),
			new(table.SceneNode),
			new(table.MonthlyBill),
			new(table.NodeBranch),
			new(table.UserSip),
		).Error
		if err != nil {
			return err
		}
		return nil
	},
}

func DB(confName string) (*gorm.DB, error) {
	specification, err := config.Get(confName)
	if err != nil {
		return nil, err
	}
	gormDB, err := db.Instance(specification)
	if err != nil {
		return nil, err
	}
	return gormDB, nil
}
