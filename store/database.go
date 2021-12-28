package store

import (
	"fmt"
	"os"
	"runtime"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore() (*Store, error) {
	if err := CreateHomeDir(); err != nil {
		return nil, err
	}

	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s/.runit/runit.db", UserHomeDir())), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

// $HOME/.runit is created if it doesn't exist
func CreateHomeDir() error {
	if _, err := os.Stat(fmt.Sprintf("%s/.runit", UserHomeDir())); os.IsNotExist(err) {
		if err := os.Mkdir(fmt.Sprintf("%s/.runit", UserHomeDir()), 0755); err != nil {
			return err
		}
	}
	return nil
}

// UserHomeDir returns the home directory of the current user
func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

// Migrate the schema
func (store *Store) AutoMigrate() {
	// store.db.AutoMigrate(&Service{})
}
