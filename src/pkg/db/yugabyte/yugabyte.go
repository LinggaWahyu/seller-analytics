package yugabyte

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// YugabyteDBConfig, database config for Yugabyte
type YugabyteDBConfig struct {
	Username string
	Password string
	Host     string
	Name     string
	Port     string
}

// NewDatabase, creates a new gorm db connection to a YugabyteDB instance
// note that this works because YSQL is postgres equivalent
func NewDatabase(cfg YugabyteDBConfig) (*gorm.DB, error) {
	conn := fmt.Sprintf("host= %s port = %s user = %s password = %s dbname = %s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Password,
		cfg.Name)
	var err error
	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
